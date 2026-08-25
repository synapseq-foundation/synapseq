// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package audio

import (
	"strconv"
	"testing"

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
