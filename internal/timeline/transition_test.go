// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package timeline

import (
	"math"
	"testing"
)

func TestApplyTransitionPoints(ts *testing.T) {
	points := []float64{0, 0.2, 1}
	for _, test := range []struct {
		progress float64
		want     float64
	}{
		{0, 0},
		{0.25, 0.1},
		{0.5, 0.2},
		{0.75, 0.6},
		{1, 1},
	} {
		if got := ApplyTransitionPoints(test.progress, points); math.Abs(got-test.want) > 1e-9 {
			ts.Errorf("ApplyTransitionPoints(%v) = %v, want %v", test.progress, got, test.want)
		}
	}
}

func TestStepAlphaWithPointsReappliesCurvePerStep(ts *testing.T) {
	points := []float64{0, 0.2, 1}
	if got := StepAlphaWithPoints(0.75, 0, points, 2); math.Abs(got-0.4) > 1e-9 {
		ts.Fatalf("StepAlphaWithPoints() = %v, want 0.4", got)
	}
}
