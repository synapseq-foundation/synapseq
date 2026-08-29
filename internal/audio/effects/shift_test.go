// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package effects

import (
	"math"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const shiftTestSampleRate = 44100

func TestApplyShiftZeroIntensityPreservesStereo(ts *testing.T) {
	processor := NewProcessor(shiftTestSampleRate, nil)
	effect := t.Effect{Type: t.EffectShift, Value: 10, Intensity: 0}

	for sample := range shiftTaps * 2 {
		left, right := processor.ApplyShift(0, effect, 1000+sample, -500-sample)
		if left != 1000+sample || right != -500-sample {
			ts.Fatalf("sample %d: got [%d %d], want [%d %d]", sample, left, right, 1000+sample, -500-sample)
		}
	}
}

func TestApplyShiftZeroSeparationPreservesStereo(ts *testing.T) {
	processor := NewProcessor(shiftTestSampleRate, nil)
	effect := t.Effect{Type: t.EffectShift, Value: 0, Intensity: 1}

	for sample := range shiftTaps * 2 {
		left, right := processor.ApplyShift(0, effect, 1000+sample, -500-sample)
		if left != 1000+sample || right != -500-sample {
			ts.Fatalf("sample %d: got [%d %d], want [%d %d]", sample, left, right, 1000+sample, -500-sample)
		}
	}
}

func TestApplyShiftCreatesSymmetricFrequencySeparation(ts *testing.T) {
	processor := NewProcessor(shiftTestSampleRate, nil)
	effect := t.Effect{Type: t.EffectShift, Value: 10, Intensity: 1}
	const carrier = 1000.0
	const amplitude = 12000.0

	// Warm the FIR and then analyze exactly one second for integer-bin isolation.
	for sample := range shiftTestSampleRate {
		input := int(math.Round(amplitude * math.Cos(2*math.Pi*carrier*float64(sample)/shiftTestSampleRate)))
		processor.ApplyShift(0, effect, input, input)
	}

	left := make([]int, shiftTestSampleRate)
	right := make([]int, shiftTestSampleRate)
	for sample := range shiftTestSampleRate {
		position := sample + shiftTestSampleRate
		input := int(math.Round(amplitude * math.Cos(2*math.Pi*carrier*float64(position)/shiftTestSampleRate)))
		left[sample], right[sample] = processor.ApplyShift(0, effect, input, input)
	}

	leftHigh := spectralMagnitude(left, carrier+effect.Value/2, shiftTestSampleRate)
	leftLow := spectralMagnitude(left, carrier-effect.Value/2, shiftTestSampleRate)
	rightHigh := spectralMagnitude(right, carrier+effect.Value/2, shiftTestSampleRate)
	rightLow := spectralMagnitude(right, carrier-effect.Value/2, shiftTestSampleRate)
	if leftHigh < leftLow*10 {
		ts.Fatalf("left channel did not shift up: high=%f low=%f", leftHigh, leftLow)
	}
	if rightLow < rightHigh*10 {
		ts.Fatalf("right channel did not shift down: low=%f high=%f", rightLow, rightHigh)
	}
}

func TestApplyShiftDerivesWetSignalFromBinauralStereoPair(ts *testing.T) {
	processor := NewProcessor(shiftTestSampleRate, nil)
	effect := t.Effect{Type: t.EffectShift, Value: 10, Intensity: 1}
	const carrier = 1000.0
	const beat = 10.0
	const amplitude = 12000.0

	for sample := range shiftTestSampleRate {
		left, right := binauralShiftTestInput(sample, carrier, beat, amplitude)
		processor.ApplyShift(0, effect, left, right)
	}

	left := make([]int, shiftTestSampleRate)
	right := make([]int, shiftTestSampleRate)
	for sample := range shiftTestSampleRate {
		position := sample + shiftTestSampleRate
		inLeft, inRight := binauralShiftTestInput(position, carrier, beat, amplitude)
		left[sample], right[sample] = processor.ApplyShift(0, effect, inLeft, inRight)
	}

	leftUpper := spectralMagnitude(left, carrier+beat, shiftTestSampleRate)
	leftLower := spectralMagnitude(left, carrier-beat, shiftTestSampleRate)
	rightUpper := spectralMagnitude(right, carrier+beat, shiftTestSampleRate)
	rightLower := spectralMagnitude(right, carrier-beat, shiftTestSampleRate)
	if leftUpper < leftLower*10 {
		ts.Fatalf("binaural-derived left wet signal did not shift upward: upper=%f lower=%f", leftUpper, leftLower)
	}
	if rightLower < rightUpper*10 {
		ts.Fatalf("binaural-derived right wet signal did not shift downward: lower=%f upper=%f", rightLower, rightUpper)
	}
}

func TestResetShiftClearsChannelState(ts *testing.T) {
	processor := NewProcessor(shiftTestSampleRate, nil)
	effect := t.Effect{Type: t.EffectShift, Value: 10, Intensity: 1}
	processor.ApplyShift(3, effect, 1000, 1000)
	processor.ResetShift(3)

	if processor.shiftStates[3] != (shiftState{}) {
		ts.Fatal("shift state was not cleared")
	}
}

func TestApplyShiftParameterChangesPreserveRuntimeState(ts *testing.T) {
	processor := NewProcessor(shiftTestSampleRate, nil)
	processor.ApplyShift(0, t.Effect{Type: t.EffectShift, Value: 10, Intensity: 0.25}, 1000, 500)
	before := processor.shiftStates[0]

	processor.ApplyShift(0, t.Effect{Type: t.EffectShift, Value: 12, Intensity: 0.75}, 800, 400)
	after := processor.shiftStates[0]

	if after.writeIndex != ringIndex(before.writeIndex+1) {
		ts.Fatalf("write index reset during parameter change: before=%d after=%d", before.writeIndex, after.writeIndex)
	}
	if after.samples != before.samples+1 {
		ts.Fatalf("oscillator sample count reset during parameter change: before=%d after=%d", before.samples, after.samples)
	}
	if after.separation != 12 {
		ts.Fatalf("separation = %f, want 12", after.separation)
	}
}

func TestApplyShiftHasNoSteadyStateAllocations(ts *testing.T) {
	processor := NewProcessor(shiftTestSampleRate, nil)
	effect := t.Effect{Type: t.EffectShift, Value: 10, Intensity: 0.5}
	allocations := testing.AllocsPerRun(100, func() {
		for range 1024 {
			processor.ApplyShift(0, effect, 1000, -500)
		}
	})
	if allocations != 0 {
		ts.Fatalf("ApplyShift allocated %f times per run", allocations)
	}
}

func BenchmarkApplyShift(b *testing.B) {
	processor := NewProcessor(shiftTestSampleRate, nil)
	effect := t.Effect{Type: t.EffectShift, Value: 10, Intensity: 0.5}
	b.ReportAllocs()
	b.SetBytes(1024 * 2 * 2)
	b.ResetTimer()
	for range b.N {
		for range 1024 {
			processor.ApplyShift(0, effect, 1000, -500)
		}
	}
}

func spectralMagnitude(samples []int, frequency float64, sampleRate int) float64 {
	var realPart, imaginaryPart float64
	for sample, value := range samples {
		phase := 2 * math.Pi * frequency * float64(sample) / float64(sampleRate)
		realPart += float64(value) * math.Cos(phase)
		imaginaryPart -= float64(value) * math.Sin(phase)
	}
	return math.Hypot(realPart, imaginaryPart)
}

func binauralShiftTestInput(sample int, carrier, beat, amplitude float64) (int, int) {
	position := float64(sample) / shiftTestSampleRate
	left := int(math.Round(amplitude * math.Cos(2*math.Pi*(carrier+beat/2)*position)))
	right := int(math.Round(amplitude * math.Cos(2*math.Pi*(carrier-beat/2)*position)))
	return left, right
}
