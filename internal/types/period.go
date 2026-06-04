// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package types

import "fmt"

// TransitionCurveK is the curve constant for logarithmic, exponential, and sigmoid transitions
const TransitionCurveK = 6.0

// TransitionType defines the type of slide for track transitions
type TransitionType int

const (
	TransitionSteady TransitionType = iota
	TransitionEaseOut
	TransitionEaseIn
	TransitionSmooth
)

// String returns the string representation of the TransitionType
func (s TransitionType) String() string {
	switch s {
	case TransitionSteady:
		return "steady"
	case TransitionEaseOut:
		return "ease-out"
	case TransitionEaseIn:
		return "ease-in"
	case TransitionSmooth:
		return "smooth"
	default:
		return "unknown"
	}
}

// TransitionTypeString returns the TransitionType for a given transition string
func TransitionTypeString(t string) TransitionType {
	switch t {
	case KeywordTransitionSteady:
		return TransitionSteady
	case KeywordTransitionEaseOut:
		return TransitionEaseOut
	case KeywordTransitionEaseIn:
		return TransitionEaseIn
	case KeywordTransitionSmooth:
		return TransitionSmooth
	default:
		return TransitionSteady
	}
}

// Period represents a time period with track configurations
type Period struct {
	Time         int
	TrackStart   [NumberOfChannels]Track
	TrackEnd     [NumberOfChannels]Track
	CrossfadeIn  [NumberOfChannels]TrackCrossfade
	CrossfadeOut [NumberOfChannels]TrackCrossfade
	Transition   TransitionType
	Steps        int
}

// TrackCrossfade describes an automatic boundary fade for a channel.
type TrackCrossfade struct {
	Active bool
	Track  Track
}

// TimeString returns the time of this period as a formatted string
func (p *Period) TimeString() string {
	hh := p.Time / 3600000
	mm := (p.Time % 3600000) / 60000
	ss := (p.Time % 60000) / 1000
	return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
}
