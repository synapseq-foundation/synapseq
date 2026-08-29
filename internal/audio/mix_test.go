// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package audio

import (
	"math"
	"testing"

	amb "github.com/synapseq-foundation/synapseq/v4/internal/audio/ambiance"
	efx "github.com/synapseq-foundation/synapseq/v4/internal/audio/effects"
	audiosync "github.com/synapseq-foundation/synapseq/v4/internal/audio/sync"
	wt "github.com/synapseq-foundation/synapseq/v4/internal/audio/wavetable"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestAudioRendererMix_PureToneFirstSampleAndPhaseAdvance(ts *testing.T) {
	renderer := newMixTestRenderer()
	renderer.channels[0] = t.Channel{
		Track: t.Track{
			Type:     t.TrackPureTone,
			Waveform: t.WaveformSquare,
		},
		WaveformStart: int(wt.SquareID),
		WaveformEnd:   int(wt.SquareID),
		Type:          t.TrackPureTone,
		Amplitude:     [2]int{4096, 4096},
		Increment:     [2]int{t.PhasePrecision, 0},
	}

	samples := renderer.mix(make([]int, t.BufferSize*audioChannels))

	expectedRaw := renderer.channels[0].Amplitude[0] * renderer.waveTables[int(wt.SquareID)][1]
	expected := clampPCM16(expectedRaw >> audioBitShift)

	if samples[0] != expected || samples[1] != expected {
		ts.Fatalf("unexpected first stereo sample: got [%d %d], want [%d %d]", samples[0], samples[1], expected, expected)
	}

	expectedOffset := (t.BufferSize * t.PhasePrecision) & phaseMask
	if renderer.channels[0].Offset[0] != expectedOffset {
		ts.Fatalf("unexpected phase offset: got %d, want %d", renderer.channels[0].Offset[0], expectedOffset)
	}
}

func TestAudioRendererMix_PanEffectRoutesMonoSignal(ts *testing.T) {
	renderer := newMixTestRenderer()
	renderer.channels[0] = t.Channel{
		Track: t.Track{
			Type:     t.TrackPureTone,
			Waveform: t.WaveformSquare,
			Effect: t.Effect{
				Type:      t.EffectPan,
				Intensity: t.IntensityPercentToRaw(100),
			},
		},
		WaveformStart: int(wt.SquareID),
		WaveformEnd:   int(wt.SquareID),
		Type:          t.TrackPureTone,
		Amplitude:     [2]int{4096, 4096},
		Increment:     [2]int{t.PhasePrecision, 0},
		Effect: t.EffectState{
			Increment: int(t.SineTableSize/4) * t.PhasePrecision,
		},
	}

	samples := renderer.mix(make([]int, t.BufferSize*audioChannels))

	expectedRaw := renderer.channels[0].Amplitude[0] * renderer.waveTables[int(wt.SquareID)][1]
	expectedRight := clampPCM16(expectedRaw >> audioBitShift)

	if samples[0] != 0 || samples[1] != expectedRight {
		ts.Fatalf("unexpected pan output: got [%d %d], want [0 %d]", samples[0], samples[1], expectedRight)
	}

	expectedEffectOffset := (t.BufferSize * int(t.SineTableSize/4) * t.PhasePrecision) & phaseMask
	if renderer.channels[0].Effect.Offset != expectedEffectOffset {
		ts.Fatalf("unexpected effect phase offset: got %d, want %d", renderer.channels[0].Effect.Offset, expectedEffectOffset)
	}
}

func TestAudioRendererMix_ModulationAffectsStereoWithSharedPhase(ts *testing.T) {
	renderer := newMixTestRenderer()
	renderer.channels[0] = t.Channel{
		Track: t.Track{
			Type:     t.TrackBinauralBeat,
			Waveform: t.WaveformSawtooth,
			Effect: t.Effect{
				Type:      t.EffectModulation,
				Intensity: t.IntensityPercentToRaw(100),
			},
		},
		WaveformStart: int(wt.SawtoothID),
		WaveformEnd:   int(wt.SawtoothID),
		Type:          t.TrackBinauralBeat,
		Amplitude:     [2]int{4096, 4096},
		Increment:     [2]int{t.PhasePrecision, t.PhasePrecision},
		Effect: t.EffectState{
			Increment: int(t.SineTableSize/2) * t.PhasePrecision,
		},
	}

	samples := renderer.mix(make([]int, t.BufferSize*audioChannels))

	baseSample := renderer.channels[0].Amplitude[0] * renderer.waveTables[int(wt.SawtoothID)][1]
	modOffset := renderer.channels[0].Effect.Increment
	modFactor := renderer.effectProcessor.CalcModulationFactor(&renderer.channels[0], modOffset)
	gain := 0.3 + 0.7*modFactor
	expected := clampPCM16(int(math.Round(float64(baseSample)*gain)) >> audioBitShift)

	if samples[0] != expected || samples[1] != expected {
		ts.Fatalf("unexpected modulation output: got [%d %d], want [%d %d]", samples[0], samples[1], expected, expected)
	}

	expectedEffectOffset := (t.BufferSize * int(t.SineTableSize/2) * t.PhasePrecision) & phaseMask
	if renderer.channels[0].Effect.Offset != expectedEffectOffset {
		ts.Fatalf("unexpected modulation phase offset: got %d, want %d", renderer.channels[0].Effect.Offset, expectedEffectOffset)
	}
}

func TestAudioRendererMix_SawtoothModulationUsesThresholdedCurve(ts *testing.T) {
	renderer := newMixTestRenderer()
	channel := &renderer.channels[0]
	channel.Track = t.Track{
		Type:     t.TrackIsochronicBeat,
		Waveform: t.WaveformSawtooth,
		Effect: t.Effect{
			Type:      t.EffectModulation,
			Intensity: t.IntensityPercentToRaw(100),
		},
	}
	channel.WaveformStart = int(wt.SawtoothID)
	channel.WaveformEnd = int(wt.SawtoothID)
	channel.Effect.Offset = int(t.SineTableSize/2) * t.PhasePrecision

	got := renderer.effectProcessor.ApplyModulationToCurrentPhase(channel, channel.Track.Effect, efx.WaveformMorphFromChannel(channel), 1000)
	modFactor := renderer.effectProcessor.CalcModulationFactor(channel, channel.Effect.Offset)
	expected := int(math.Round(1000 * (0.3 + 0.7*modFactor)))
	if got != expected {
		ts.Fatalf("unexpected sawtooth modulation output for current curve: got %d, want %d", got, expected)
	}
}

func TestAudioRendererMix_PureToneMorphsBetweenWaveforms(ts *testing.T) {
	renderer := newMixTestRenderer()
	renderer.channels[0] = t.Channel{
		Track: t.Track{
			Type:     t.TrackPureTone,
			Waveform: t.WaveformSine,
		},
		WaveformStart: int(wt.SineID),
		WaveformEnd:   int(wt.SquareID),
		WaveformAlpha: 0.25,
		Type:          t.TrackPureTone,
		Amplitude:     [2]int{4096, 4096},
		Increment:     [2]int{t.PhasePrecision, 0},
	}

	samples := renderer.mix(make([]int, t.BufferSize*audioChannels))

	sine := float64(renderer.waveTables[int(wt.SineID)][1])
	square := float64(renderer.waveTables[int(wt.SquareID)][1])
	blended := sine + (square-sine)*0.25
	expectedRaw := int(math.Round(float64(renderer.channels[0].Amplitude[0]) * blended))
	expected := clampPCM16(expectedRaw >> audioBitShift)

	if samples[0] != expected || samples[1] != expected {
		ts.Fatalf("unexpected morphed sample: got [%d %d], want [%d %d]", samples[0], samples[1], expected, expected)
	}
}

func TestAudioRendererMix_AppliesStereoAmplitudeAfterPan(ts *testing.T) {
	renderer := newMixTestRenderer()
	renderer.channels[0] = t.Channel{
		Track: t.Track{
			Type:     t.TrackPureTone,
			Waveform: t.WaveformSquare,
		},
		WaveformStart: int(wt.SquareID),
		WaveformEnd:   int(wt.SquareID),
		Type:          t.TrackPureTone,
		Amplitude:     [2]int{4096, 2048},
		Increment:     [2]int{t.PhasePrecision, 0},
	}

	samples := renderer.mix(make([]int, t.BufferSize*audioChannels))
	if samples[0] != 32767 || samples[1] != 16383 {
		ts.Fatalf("unexpected stereo amplitude output: got [%d %d]", samples[0], samples[1])
	}
}

func TestAudioRendererMix_UsesCrossfadeCueAmplitude(ts *testing.T) {
	var p0, p1 t.Period
	p0.Time = 0
	p1.Time = 60_000

	track := t.Track{Type: t.TrackPureTone, Carrier: 440, Amplitude: t.AmplitudePercentToRaw(40), Waveform: t.WaveformSquare}
	p0.TrackStart[0] = track
	p0.TrackEnd[0] = track
	p0.CrossfadeOut[0] = t.TrackCrossfade{Active: true, Track: track}

	plan := compileRenderPlan([]t.Period{p0, p1}, 44100)
	renderer := newMixTestRenderer()

	fullCue := plan.cue(0, 30_000)
	renderer.applyCueSignalState(fullCue)
	renderer.channels[0].Offset = [2]int{}
	fullSamples := renderer.mix(make([]int, t.BufferSize*audioChannels))

	halfCue := plan.cue(0, 45_000)
	renderer.applyCueSignalState(halfCue)
	renderer.channels[0].Offset = [2]int{}
	halfSamples := renderer.mix(make([]int, t.BufferSize*audioChannels))

	if absInt(halfSamples[0]) >= absInt(fullSamples[0]) {
		ts.Fatalf("expected crossfade cue to reduce mixed output: full=%d half=%d", fullSamples[0], halfSamples[0])
	}
}

func TestAudioRendererMix_ModulationSlewsAbruptSquareGainChanges(ts *testing.T) {
	renderer := newMixTestRenderer()
	channel := &renderer.channels[0]
	channel.Track = t.Track{
		Type:     t.TrackAmbiance,
		Waveform: t.WaveformSquare,
		Effect: t.Effect{
			Type:      t.EffectModulation,
			Intensity: t.IntensityPercentToRaw(100),
		},
	}
	channel.WaveformStart = int(wt.SquareID)
	channel.WaveformEnd = int(wt.SquareID)
	channel.Effect.ModulationGain = 1
	channel.Effect.ModulationInitialized = true

	got := renderer.effectProcessor.ApplyModulationToCurrentPhase(channel, channel.Track.Effect, efx.WaveformMorphFromChannel(channel), 1000)
	expectedFloor := int(1000 * 0.3)
	if got <= expectedFloor || got >= 1000 {
		ts.Fatalf("unexpected slewed modulation output: got %d, want between %d and 1000", got, expectedFloor)
	}
	if channel.Effect.ModulationGain >= 1 || channel.Effect.ModulationGain <= 0.3 {
		ts.Fatalf("unexpected slewed modulation gain: got %f", channel.Effect.ModulationGain)
	}
}

func TestAudioRendererMix_PanUsesWaveformAndSlewsSquareSwitches(ts *testing.T) {
	renderer := newMixTestRenderer()
	channel := &renderer.channels[0]
	channel.Track = t.Track{
		Type:     t.TrackAmbiance,
		Waveform: t.WaveformSquare,
		Effect: t.Effect{
			Type:      t.EffectPan,
			Intensity: t.IntensityPercentToRaw(100),
		},
	}
	channel.WaveformStart = int(wt.SquareID)
	channel.WaveformEnd = int(wt.SquareID)
	channel.Effect.PanPosition = -1
	channel.Effect.PanInitialized = true
	channel.Effect.Offset = int(t.SineTableSize/4) * t.PhasePrecision

	left, right := renderer.effectProcessor.ApplyPan(channel, 1000, 1000)
	if left <= 0 || right <= 0 {
		ts.Fatalf("unexpected hard-switched pan output: got [%d %d]", left, right)
	}
	if channel.Effect.PanPosition <= -1 || channel.Effect.PanPosition >= 1 {
		ts.Fatalf("unexpected slewed pan position: got %f", channel.Effect.PanPosition)
	}
	if channel.Effect.PanPosition >= -0.9 {
		ts.Fatalf("expected a small initial slew step, got pan position %f", channel.Effect.PanPosition)
	}
	if right <= 0 || left >= 1000 {
		ts.Fatalf("expected pan to start moving toward right channel without hard switch: got [%d %d]", left, right)
	}
}

func TestAudioRendererMix_AmbianceUsesPreparedStereoBuffer(ts *testing.T) {
	renderer := newMixTestRenderer()
	renderer.ambianceState = amb.NewTestRuntime(1)
	renderer.ambianceState.SetChannelBuffer(0, []int{20000, -10000})
	renderer.ambianceState.SetChannelIndex(0, 0)
	renderer.channels[0] = t.Channel{
		Track: t.Track{
			Type: t.TrackAmbiance,
		},
		Type:      t.TrackAmbiance,
		Amplitude: [2]int{3, 3},
	}

	samples := renderer.mix(make([]int, t.BufferSize*audioChannels))

	expectedLeft := clampPCM16((20000 * 16 * 3) >> audioBitShift)
	expectedRight := clampPCM16((-10000 * 16 * 3) >> audioBitShift)

	if samples[0] != expectedLeft || samples[1] != expectedRight {
		ts.Fatalf("unexpected ambiance output: got [%d %d], want [%d %d]", samples[0], samples[1], expectedLeft, expectedRight)
	}
}

func TestAudioRendererMix_AmbianceAppliesShiftDryWet(ts *testing.T) {
	renderer := newMixTestRenderer()
	renderer.ambianceState = amb.NewTestRuntime(1)
	buffer := make([]int, t.BufferSize*audioChannels)
	for frame := range t.BufferSize {
		buffer[frame*2] = 20000
		buffer[frame*2+1] = -10000
	}
	renderer.ambianceState.SetChannelBuffer(0, buffer)
	renderer.ambianceState.SetChannelIndex(0, 0)
	renderer.channels[0] = t.Channel{
		Track: t.Track{
			Type:   t.TrackAmbiance,
			Effect: t.Effect{Type: t.EffectShift, Value: 10, Intensity: 0.5},
		},
		Type:      t.TrackAmbiance,
		Amplitude: [2]int{3, 3},
	}

	samples := renderer.mix(make([]int, t.BufferSize*audioChannels))
	expectedLeft := clampPCM16((10000 * 16 * 3) >> audioBitShift)
	expectedRight := clampPCM16((-5000 * 16 * 3) >> audioBitShift)
	if samples[0] != expectedLeft || samples[1] != expectedRight {
		ts.Fatalf("unexpected shifted ambiance onset: got [%d %d], want [%d %d]", samples[0], samples[1], expectedLeft, expectedRight)
	}
	if samples[200] == expectedLeft && samples[201] == expectedRight {
		ts.Fatal("shift wet path did not reach the mixer after FIR warmup")
	}
}

func TestAudioRendererMix_ShiftProcessesGeneratedTrackFamilies(ts *testing.T) {
	const warmupFrames = 64
	tests := []struct {
		name      string
		trackType t.TrackType
		increment [2]int
	}{
		{name: "pure", trackType: t.TrackPureTone, increment: [2]int{frequencyToIncrement(44100, 1000)}},
		{name: "binaural", trackType: t.TrackBinauralBeat, increment: [2]int{frequencyToIncrement(44100, 1005), frequencyToIncrement(44100, 995)}},
		{name: "monaural", trackType: t.TrackMonauralBeat, increment: [2]int{frequencyToIncrement(44100, 1005), frequencyToIncrement(44100, 995)}},
		{name: "isochronic", trackType: t.TrackIsochronicBeat, increment: [2]int{frequencyToIncrement(44100, 1000), frequencyToIncrement(44100, 10)}},
		{name: "white noise", trackType: t.TrackWhiteNoise},
		{name: "pink noise", trackType: t.TrackPinkNoise},
		{name: "brown noise", trackType: t.TrackBrownNoise},
	}

	for _, test := range tests {
		ts.Run(test.name, func(ts *testing.T) {
			renderer := newMixTestRenderer()
			renderer.channels[0] = t.Channel{
				Track: t.Track{
					Type:      test.trackType,
					Waveform:  t.WaveformSine,
					Effect:    t.Effect{Type: t.EffectShift, Value: 100, Intensity: 1},
					Amplitude: t.AmplitudePercentToRaw(25),
				},
				WaveformStart: int(wt.SineID),
				WaveformEnd:   int(wt.SineID),
				Type:          test.trackType,
				Amplitude:     [2]int{1024, 1024},
				Increment:     test.increment,
			}

			samples := renderer.mix(make([]int, t.BufferSize*audioChannels))
			var diverged bool
			for frame := warmupFrames; frame < t.BufferSize; frame++ {
				if samples[frame*2] != samples[frame*2+1] {
					diverged = true
					break
				}
			}
			if !diverged {
				ts.Fatalf("shift did not create stereo divergence for %s", test.name)
			}
		})
	}
}

func TestAudioRendererMix_UsesOnlyActiveCueChannels(ts *testing.T) {
	renderer := newMixTestRenderer()
	renderer.channels[0] = t.Channel{
		Track:         t.Track{Type: t.TrackPureTone, Waveform: t.WaveformSquare},
		WaveformStart: int(wt.SquareID),
		WaveformEnd:   int(wt.SquareID),
		Type:          t.TrackPureTone,
		Amplitude:     [2]int{4096, 4096},
		Increment:     [2]int{t.PhasePrecision, 0},
	}
	for channel := 1; channel < t.NumberOfChannels; channel++ {
		renderer.channels[channel] = t.Channel{
			Track:         t.Track{Type: t.TrackPureTone, Waveform: t.WaveformSquare},
			WaveformStart: int(wt.SquareID),
			WaveformEnd:   int(wt.SquareID),
			Type:          t.TrackPureTone,
			Amplitude:     [2]int{4096, 4096},
			Increment:     [2]int{t.PhasePrecision, 0},
		}
	}

	var cue audiosync.Cue
	cue.Channels[0] = audiosync.ChannelCue{
		Track:         renderer.channels[0].Track,
		WaveformStart: wt.SquareID,
		WaveformEnd:   wt.SquareID,
		Amplitude:     [2]int{4096, 4096},
		Increment:     [2]int{t.PhasePrecision, 0},
	}
	renderer.applyCueSignalState(cue)

	samples := renderer.mix(make([]int, t.BufferSize*audioChannels))
	if samples[0] != 32767 || samples[1] != 32767 {
		ts.Fatalf("unexpected active-channel output: [%d %d]", samples[0], samples[1])
	}
	if renderer.channels[1].Offset != [2]int{} {
		ts.Fatalf("inactive channel phase advanced: %v", renderer.channels[1].Offset)
	}
}

func newMixTestRenderer() *AudioRenderer {
	renderer := &AudioRenderer{
		waveTables:      wt.Init(),
		noiseGenerator:  NewNoiseGenerator(),
		effectProcessor: efx.NewProcessor(44100, wt.Init()),
		ambianceState:   amb.NewTestRuntime(0),
		AudioRendererOptions: &AudioRendererOptions{
			SampleRate: 44100,
			Volume:     100,
		},
	}

	return renderer
}
