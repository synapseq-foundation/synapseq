// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package ambiance

import (
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type finiteSampleAudio struct {
	data [][]int
	pos  []int
}

func newFiniteSampleAudio(data [][]int) *finiteSampleAudio {
	copied := make([][]int, len(data))
	for i := range data {
		copied[i] = append([]int(nil), data[i]...)
	}
	return &finiteSampleAudio{
		data: copied,
		pos:  make([]int, len(data)),
	}
}

func (fa *finiteSampleAudio) ReadSamplesAt(index int, samples []int, numSamples int) (int, error) {
	if numSamples > len(samples) {
		numSamples = len(samples)
	}
	if index < 0 || index >= len(fa.data) {
		return 0, nil
	}

	for i := 0; i < numSamples; i++ {
		if fa.pos[index] >= len(fa.data[index]) {
			samples[i] = 0
			continue
		}
		samples[i] = fa.data[index][fa.pos[index]]
		fa.pos[index]++
	}
	return numSamples, nil
}

func (fa *finiteSampleAudio) Close() error {
	return nil
}

func TestMusicRuntimeKeepsEOFStatePerChannel(ts *testing.T) {
	var p0, p1, p2 t.Period
	p0.Time = 0
	p1.Time = 1000
	p2.Time = 2000

	p0.TrackStart[0] = t.Track{Type: t.TrackMusic, SourceName: "meditation"}
	p1.TrackStart[1] = t.Track{Type: t.TrackMusic, SourceName: "meditation"}

	sourceData := [][]int{{10, 11, 12, 13}}
	runtime, err := NewMusicRuntime(
		[]t.Period{p0, p1, p2},
		map[string]string{"meditation": "meditation.mp3"},
		44100,
		func(paths []string, sampleRate int) (SampleAudio, error) {
			return newFiniteSampleAudio(sourceData), nil
		},
	)
	if err != nil {
		ts.Fatalf("NewMusicRuntime: %v", err)
	}
	defer runtime.Close()

	channels := make([]t.Channel, t.NumberOfChannels)
	channels[0].Track.Type = t.TrackMusic
	runtime.UpdateChannelIndex(0, 0, t.TrackMusic)
	runtime.CollectActiveIndices(channels)
	runtime.PrepareBuffers(3)

	gotCh0 := append([]int(nil), runtime.ChannelBuffer(0)...)
	wantEOF := []int{10, 11, 12, 13, 0, 0}
	if !equalSamples(gotCh0, wantEOF) {
		ts.Fatalf("channel 0 expected EOF-padded samples %v, got %v", wantEOF, gotCh0)
	}

	channels[0].Track.Type = t.TrackOff
	channels[1].Track.Type = t.TrackMusic
	runtime.UpdateChannelIndex(0, 1, t.TrackOff)
	runtime.UpdateChannelIndex(1, 1, t.TrackMusic)
	runtime.CollectActiveIndices(channels)
	runtime.PrepareBuffers(3)

	gotCh1 := append([]int(nil), runtime.ChannelBuffer(1)...)
	if !equalSamples(gotCh1, wantEOF) {
		ts.Fatalf("channel 1 should start the same music from the beginning, want %v got %v", wantEOF, gotCh1)
	}
}

func equalSamples(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
