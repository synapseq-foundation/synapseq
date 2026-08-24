// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sequence

import (
	"github.com/synapseq-foundation/synapseq/v4/internal/parser"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func buildTrackFromDeclaration(sourceFile string, lineNumber int, lineText string, decl *parser.ParsedTrackDeclaration) (*t.Track, error) {
	rightAmplitude := decl.AmplitudePercent
	if decl.HasRightAmplitude {
		rightAmplitude = decl.RightAmplitudePercent
	}

	track := &t.Track{
		Type:      decl.Type,
		Carrier:   decl.Carrier,
		Resonance: decl.Resonance,
		Amplitude: [2]t.AmplitudeType{
			t.AmplitudePercentToRawChannel(decl.AmplitudePercent),
			t.AmplitudePercentToRawChannel(rightAmplitude),
		},
		SourceName:  decl.SourceName,
		NoiseSmooth: decl.NoiseSmooth,
		Waveform:    decl.Waveform,
		Effect: t.Effect{
			Type:      decl.EffectType,
			Value:     decl.EffectValue,
			Intensity: t.IntensityPercentToRaw(decl.EffectIntensityPercent),
		},
	}

	if err := track.Validate(); err != nil {
		return nil, lineDiagnostic(sourceFile, lineNumber, lineText, err.Error())
	}

	return track, nil
}
