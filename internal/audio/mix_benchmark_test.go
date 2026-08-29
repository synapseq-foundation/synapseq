// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package audio

import (
	"strconv"
	"testing"

	amb "github.com/synapseq-foundation/synapseq/v4/internal/audio/ambiance"
	efx "github.com/synapseq-foundation/synapseq/v4/internal/audio/effects"
	wt "github.com/synapseq-foundation/synapseq/v4/internal/audio/wavetable"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

var benchmarkSample int

func BenchmarkMix(b *testing.B) {
	for _, tracks := range []int{1, 4, 16} {
		b.Run("tracks="+strconv.Itoa(tracks), func(b *testing.B) {
			renderer := newBenchmarkMixRenderer(tracks)
			samples := make([]int, t.BufferSize*audioChannels)

			b.ReportAllocs()
			b.SetBytes(int64(len(samples) * 2))
			b.ResetTimer()
			for range b.N {
				renderer.mix(samples)
			}
			benchmarkSample = samples[0]
		})
	}
}

func BenchmarkMixShift(b *testing.B) {
	sources := []struct {
		name string
		kind t.TrackType
	}{
		{name: "pure", kind: t.TrackPureTone},
		{name: "binaural", kind: t.TrackBinauralBeat},
		{name: "noise", kind: t.TrackPinkNoise},
		{name: "ambiance", kind: t.TrackAmbiance},
	}
	effects := []struct {
		name   string
		effect t.Effect
	}{
		{name: "off"},
		{name: "shift", effect: t.Effect{Type: t.EffectShift, Value: 10, Intensity: 0.5}},
	}

	for _, source := range sources {
		for _, effect := range effects {
			for _, sampleRate := range []int{44100, 96000} {
				for _, tracks := range []int{1, 4, 16} {
					name := "source=" + source.name + "/effect=" + effect.name +
						"/rate=" + strconv.Itoa(sampleRate) + "/tracks=" + strconv.Itoa(tracks)
					b.Run(name, func(b *testing.B) {
						renderer := newBenchmarkShiftRenderer(tracks, source.kind, effect.effect, sampleRate)
						samples := make([]int, t.BufferSize*audioChannels)

						b.ReportAllocs()
						b.SetBytes(int64(len(samples) * 2))
						b.ResetTimer()
						for range b.N {
							renderer.mix(samples)
						}
						benchmarkSample = samples[0]
					})
				}
			}
		}
	}
}

func newBenchmarkMixRenderer(tracks int) *AudioRenderer {
	renderer := newMixTestRenderer()
	for channel := range tracks {
		trackType := t.TrackPureTone
		carrier := 180.0 + float64(channel*10)
		resonance := 0.0
		if channel%2 == 1 {
			trackType = t.TrackBinauralBeat
			resonance = 8
		}

		renderer.channels[channel] = t.Channel{
			Track: t.Track{
				Type:      trackType,
				Carrier:   carrier,
				Resonance: resonance,
				Waveform:  t.WaveformSine,
			},
			WaveformStart: int(wt.SineID),
			WaveformEnd:   int(wt.SineID),
			Type:          trackType,
			Amplitude:     [2]int{256, 256},
			Increment: [2]int{
				frequencyToIncrement(44100, carrier),
				frequencyToIncrement(44100, carrier+resonance),
			},
		}
		renderer.activeChannels[channel] = channel
		renderer.signals[channel] = channelSignalState{
			resolved:  true,
			kind:      trackType,
			waveform:  efx.WaveformMorph{Start: wt.SineID, End: wt.SineID},
			amplitude: [2]int{256, 256},
			increment: [2]int{
				frequencyToIncrement(44100, carrier),
				frequencyToIncrement(44100, carrier+resonance),
			},
		}
	}
	renderer.activeCount = tracks
	renderer.hasChannelPlan = true

	return renderer
}

func newBenchmarkShiftRenderer(tracks int, trackType t.TrackType, effect t.Effect, sampleRate int) *AudioRenderer {
	renderer := newMixTestRenderer()
	renderer.SampleRate = sampleRate
	renderer.effectProcessor = efx.NewProcessor(sampleRate, renderer.waveTables)
	if trackType == t.TrackAmbiance {
		renderer.ambianceState = amb.NewTestRuntime(1)
		buffer := make([]int, t.BufferSize*audioChannels)
		for frame := range t.BufferSize {
			buffer[frame*2] = 12000
			buffer[frame*2+1] = -8000
		}
		renderer.ambianceState.SetChannelBuffer(0, buffer)
	}

	for channel := range tracks {
		carrier := 180.0 + float64(channel*10)
		increment := [2]int{frequencyToIncrement(sampleRate, carrier)}
		resonance := 0.0
		if trackType == t.TrackBinauralBeat {
			resonance = 8
			increment = [2]int{
				frequencyToIncrement(sampleRate, carrier+resonance/2),
				frequencyToIncrement(sampleRate, carrier-resonance/2),
			}
		}

		track := t.Track{
			Type:      trackType,
			Carrier:   carrier,
			Resonance: resonance,
			Waveform:  t.WaveformSine,
			Effect:    effect,
		}
		renderer.channels[channel] = t.Channel{
			Track:         track,
			WaveformStart: int(wt.SineID),
			WaveformEnd:   int(wt.SineID),
			Type:          trackType,
			Amplitude:     [2]int{256, 256},
			Increment:     increment,
		}
		renderer.activeChannels[channel] = channel
		renderer.signals[channel] = channelSignalState{
			resolved:  true,
			kind:      trackType,
			effect:    effect,
			waveform:  efx.WaveformMorph{Start: wt.SineID, End: wt.SineID},
			amplitude: [2]int{256, 256},
			increment: increment,
		}
		if trackType == t.TrackAmbiance {
			renderer.ambianceState.SetChannelIndex(channel, 0)
		}
	}
	renderer.activeCount = tracks
	renderer.hasChannelPlan = true
	return renderer
}
