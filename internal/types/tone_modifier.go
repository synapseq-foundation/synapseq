// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package types

// ToneModifier represents configuration specific to a tone type.
// Implementations are restricted to this package.
type ToneModifier interface {
	isToneModifier()
	String() string
}

// BinauralMode specifies the frequency ordering of a binaural pair.
type BinauralMode int

const (
	// BinauralLeftRight places the higher frequency on the left.
	BinauralLeftRight BinauralMode = iota
	// BinauralRightLeft places the higher frequency on the right.
	BinauralRightLeft
)

// BinauralModifier configures a binaural tone. Store it as a value, not a pointer.
type BinauralModifier struct {
	Mode BinauralMode
}

func (BinauralModifier) isToneModifier() {}

// BinauralMode returns the effective mode, including the implicit left-right default.
func (tr Track) BinauralMode() BinauralMode {
	if modifier, ok := tr.ToneModifier.(BinauralModifier); ok {
		return modifier.Mode
	}
	return BinauralLeftRight
}

// String returns the binaural mode keyword.
func (modifier BinauralModifier) String() string {
	switch modifier.Mode {
	case BinauralLeftRight:
		return KeywordBinauralLeftRight
	case BinauralRightLeft:
		return KeywordBinauralRightLeft
	default:
		return ""
	}
}

var defaultToneModifiers = map[TrackType]ToneModifier{
	TrackBinauralBeat:   BinauralModifier{},
	TrackIsochronicBeat: IsochronicModifier{},
}

// ToneModifierString returns the configured or default modifier's keyword.
// Tracks without a modifier return an empty string.
func (tr Track) ToneModifierString() string {
	modifier := tr.ToneModifier
	if modifier == nil {
		modifier = defaultToneModifiers[tr.Type]
	}
	if modifier == nil {
		return ""
	}
	return modifier.String()
}

// IsochronicMode specifies the pulse envelope of an isochronic tone.
type IsochronicMode int

const (
	// IsochronicStandard preserves the original pulse envelope.
	IsochronicStandard IsochronicMode = iota
	// IsochronicShort compresses the envelope to half a cycle at 75 percent gain.
	IsochronicShort
	// IsochronicWide uses a fixed envelope over 75 percent of each cycle.
	IsochronicWide
)

// IsochronicModifier configures an isochronic tone; store it as a value.
type IsochronicModifier struct{ Mode IsochronicMode }

func (IsochronicModifier) isToneModifier() {}
func (modifier IsochronicModifier) String() string {
	switch modifier.Mode {
	case IsochronicStandard:
		return KeywordIsochronicStandard
	case IsochronicWide:
		return KeywordIsochronicWide
	case IsochronicShort:
		return KeywordIsochronicShort
	default:
		return ""
	}
}

// IsochronicMode returns the effective mode, including the implicit standard default.
func (tr Track) IsochronicMode() IsochronicMode {
	if modifier, ok := tr.ToneModifier.(IsochronicModifier); ok {
		return modifier.Mode
	}
	return IsochronicStandard
}
