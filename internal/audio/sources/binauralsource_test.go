/*
 * SynapSeq - Text-Driven Audio Sequencer for Brainwave Entrainment
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

package sources

import (
	"math"
	"testing"

	efx "github.com/synapseq-foundation/synapseq/v4/internal/audio/effects"
	wt "github.com/synapseq-foundation/synapseq/v4/internal/audio/wavetable"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestBinauralSampleUsesCompiledWaveformMorph(ts *testing.T) {
	processor := efx.NewProcessor(44100, wt.Init())
	source := NewBinaural(Signal{
		Waveform: efx.WaveformMorph{Start: t.WaveformSawtooth, End: t.WaveformSawtooth, Alpha: 0},
		Amplitude: [2]int{4096, 2048},
	})

	left, right := source.Sample(processor, t.PhasePrecision, t.PhasePrecision)
	wave := wt.Init()[int(t.WaveformSawtooth)][1]
	if left != 4096*wave {
		ts.Fatalf("unexpected binaural left sample: got %d, want %d", left, 4096*wave)
	}
	if right != 2048*wave {
		ts.Fatalf("unexpected binaural right sample: got %d, want %d", right, 2048*wave)
	}
}

func TestMonauralSampleUsesCompiledWaveformMorph(ts *testing.T) {
	processor := efx.NewProcessor(44100, wt.Init())
	source := NewMonaural(Signal{
		Waveform: efx.WaveformMorph{Start: t.WaveformSine, End: t.WaveformSquare, Alpha: 0.25},
		Amplitude: [2]int{4096, 0},
	})

	got := source.Sample(processor, t.PhasePrecision, t.PhasePrecision)
	sine := float64(wt.Init()[int(t.WaveformSine)][1])
	square := float64(wt.Init()[int(t.WaveformSquare)][1])
	blended := sine + (square-sine)*0.25
	wave := int(math.Round(blended))
	want := (4096 * (wave + wave)) >> 1
	if got != want {
		ts.Fatalf("unexpected monaural sample: got %d, want %d", got, want)
	}
}

func TestIsochronicSampleUsesCompiledWaveformMorph(ts *testing.T) {
	processor := efx.NewProcessor(44100, wt.Init())
	source := NewIsochronic(Signal{
		Waveform: efx.WaveformMorph{Start: t.WaveformSine, End: t.WaveformSquare, Alpha: 0.25},
		Amplitude: [2]int{4096, 0},
	})

	got := source.Sample(processor, t.PhasePrecision, 0.5)
	want := int(float64(4096*processor.WaveformSampleForMorph(efx.WaveformMorph{Start: t.WaveformSine, End: t.WaveformSquare, Alpha: 0.25}, t.PhasePrecision)) * 0.5)
	if got != want {
		ts.Fatalf("unexpected isochronic sample: got %d, want %d", got, want)
	}
}