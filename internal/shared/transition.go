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

package shared

import (
	"math"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const (
	periodStepMinLegMs = 5000
	periodStepHardCap  = 12
)

// MaxPeriodSteps returns the maximum number of steps supported for a period duration.
func MaxPeriodSteps(durationMs int) int {
	stepSlots := durationMs / periodStepMinLegMs
	if stepSlots <= 1 {
		return 0
	}

	maxSteps := (stepSlots - 1) / 2
	if maxSteps > periodStepHardCap {
		return periodStepHardCap
	}
	if maxSteps < 0 {
		return 0
	}

	return maxSteps
}

// ApplyTransitionAlpha maps linear progress through the configured transition curve.
func ApplyTransitionAlpha(progress float64, transition t.TransitionType) float64 {
	alpha := clampUnit(progress)

	switch transition {
	case t.TransitionEaseOut:
		alpha = math.Log1p(math.Expm1(t.TransitionCurveK)*alpha) / t.TransitionCurveK
	case t.TransitionEaseIn:
		alpha = math.Expm1(t.TransitionCurveK*alpha) / math.Expm1(t.TransitionCurveK)
	case t.TransitionSmooth:
		raw := 1.0 / (1.0 + math.Exp(-t.TransitionCurveK*(alpha-0.5)))
		min := 1.0 / (1.0 + math.Exp(t.TransitionCurveK*0.5))
		max := 1.0 / (1.0 + math.Exp(-t.TransitionCurveK*0.5))
		alpha = (raw - min) / (max - min)
	}

	return alpha
}

// StepAlpha computes the effective interpolation alpha for a period.
//
// Steps=0 preserves the current monotonic transition. Steps>0 creates an
// alternating trajectory with 2*steps+1 legs so the period always starts at
// TrackStart and ends exactly at TrackEnd.
func StepAlpha(progress float64, transition t.TransitionType, steps int) float64 {
	progress = clampUnit(progress)
	if steps <= 0 {
		return ApplyTransitionAlpha(progress, transition)
	}
	if progress >= 1 {
		return 1
	}

	totalLegs := 2*steps + 1
	legSpan := 1.0 / float64(totalLegs)
	legIndex := int(math.Floor(progress / legSpan))
	if legIndex >= totalLegs {
		legIndex = totalLegs - 1
	}

	legStart := float64(legIndex) * legSpan
	legProgress := clampUnit((progress - legStart) / legSpan)
	curved := ApplyTransitionAlpha(legProgress, transition)
	if legIndex%2 == 0 {
		return curved
	}

	return 1 - curved
}

func clampUnit(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}

	return value
}
