// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package audiosource

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
	runtime, err := NewRuntime(
		[]t.Period{p0, p1, p2},
		map[string]string{"meditation": "meditation.mp3"},
		44100,
		RuntimeOptions{
			TrackType:  t.TrackMusic,
			SourceKind: "music",
			Scope:      BufferScopeChannel,
			NewAudio: func(paths []string, sampleRate int) (SampleAudio, error) {
				return newFiniteSampleAudio(sourceData), nil
			},
		},
	)
	if err != nil {
		ts.Fatalf("NewRuntime: %v", err)
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

func TestReachableSourcesIncludesCrossfades(ts *testing.T) {
	var period t.Period
	period.TrackStart[0] = t.Track{Type: t.TrackMusic, SourceName: "start"}
	period.TrackEnd[1] = t.Track{Type: t.TrackMusic, SourceName: "end"}
	period.CrossfadeOut[2] = t.TrackCrossfade{Track: t.Track{Type: t.TrackMusic, SourceName: "out"}}
	period.CrossfadeIn[3] = t.TrackCrossfade{Track: t.Track{Type: t.TrackMusic, SourceName: "in"}}

	sources := map[string]string{
		"start":  "start.mp3",
		"end":    "end.mp3",
		"out":    "out.mp3",
		"in":     "in.mp3",
		"unused": "unused.mp3",
	}
	reachable := ReachableSources([]t.Period{period}, sources, t.TrackMusic)
	if len(reachable) != 4 {
		ts.Fatalf("expected four reachable sources, got %v", reachable)
	}
	if _, ok := reachable["unused"]; ok {
		ts.Fatalf("unexpected unused source: %v", reachable)
	}
}

func TestRuntimeLoadsOnlyReachableSources(ts *testing.T) {
	periods := []t.Period{{}}
	periods[0].TrackStart[0] = t.Track{Type: t.TrackAmbiance, SourceName: "rain"}
	var loadedPaths []string

	runtime, err := NewRuntime(
		periods,
		map[string]string{"rain": "rain.wav", "unused": "unused.wav"},
		44100,
		RuntimeOptions{
			TrackType:  t.TrackAmbiance,
			SourceKind: "ambiance",
			Scope:      BufferScopeSource,
			NewAudio: func(paths []string, sampleRate int) (SampleAudio, error) {
				loadedPaths = append([]string{}, paths...)
				return newFiniteSampleAudio([][]int{{}}), nil
			},
		},
	)
	if err != nil {
		ts.Fatalf("NewRuntime: %v", err)
	}
	defer runtime.Close()

	if len(loadedPaths) != 1 || loadedPaths[0] != "rain.wav" {
		ts.Fatalf("expected only rain.wav to load, got %v", loadedPaths)
	}
}

func TestChannelRuntimeDeduplicatesIdenticalSourcePaths(ts *testing.T) {
	periods := []t.Period{{}}
	periods[0].TrackStart[0] = t.Track{Type: t.TrackMusic, SourceName: "intro"}
	periods[0].TrackStart[1] = t.Track{Type: t.TrackMusic, SourceName: "repeat"}
	var loadedPaths []string

	runtime, err := NewRuntime(
		periods,
		map[string]string{"intro": "track.mp3", "repeat": "track.mp3"},
		44100,
		RuntimeOptions{
			TrackType:  t.TrackMusic,
			SourceKind: "music",
			Scope:      BufferScopeChannel,
			NewAudio: func(paths []string, sampleRate int) (SampleAudio, error) {
				loadedPaths = append([]string(nil), paths...)
				return newFiniteSampleAudio([][]int{{}}), nil
			},
		},
	)
	if err != nil {
		ts.Fatalf("NewRuntime: %v", err)
	}
	defer runtime.Close()

	if len(loadedPaths) != 1 || loadedPaths[0] != "track.mp3" {
		ts.Fatalf("expected one shared source path, got %v", loadedPaths)
	}
	runtime.UpdateChannelIndex(0, 0, t.TrackMusic)
	runtime.UpdateChannelIndex(1, 0, t.TrackMusic)
	if runtime.channelIdx[0] != 0 || runtime.channelIdx[1] != 0 {
		ts.Fatalf("expected aliases to share source index, got %v", runtime.channelIdx[:2])
	}
}

func TestSourceRuntimeReleasesInactiveBuffers(ts *testing.T) {
	runtime := NewTestRuntime(2)
	runtime.samplesByIdx[0] = make([]int, 32)
	runtime.samplesByIdx[1] = make([]int, 32)

	channels := make([]t.Channel, t.NumberOfChannels)
	runtime.CollectActiveIndices(channels)

	if runtime.samplesByIdx[0] != nil || runtime.samplesByIdx[1] != nil {
		ts.Fatalf("expected inactive source buffers to be released, got %v", runtime.samplesByIdx)
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
