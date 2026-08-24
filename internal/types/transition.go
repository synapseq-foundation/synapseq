// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package types

import "strings"

const (
	MinTransitionPoints = 2
	MaxTransitionPoints = 256
)

// TransitionDefinition contains one normalized monotonic interpolation curve.
type TransitionDefinition struct {
	Name   string
	Points []float64
}

// IsBuiltinTransitionName reports whether name conflicts with a built-in transition.
func IsBuiltinTransitionName(name string) bool {
	switch strings.ToLower(name) {
	case KeywordTransitionSteady, KeywordTransitionEaseOut, KeywordTransitionEaseIn, KeywordTransitionSmooth:
		return true
	default:
		return false
	}
}

// BuiltinTransitionNames returns the transition names accepted without a definition.
func BuiltinTransitionNames() []string {
	return []string{KeywordTransitionSteady, KeywordTransitionEaseOut, KeywordTransitionEaseIn, KeywordTransitionSmooth}
}
