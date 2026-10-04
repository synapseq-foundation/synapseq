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
	TrackBinauralBeat: BinauralModifier{},
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
