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

package sequence

import (
	"github.com/synapseq-foundation/synapseq/v4/internal/parser"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func buildTrackFromDeclaration(sourceFile string, lineNumber int, lineText string, decl *parser.ParsedTrackDeclaration) (*t.Track, error) {
	track := &t.Track{
		Type:         decl.Type,
		Carrier:      decl.Carrier,
		Resonance:    decl.Resonance,
		Amplitude:    t.AmplitudePercentToRaw(decl.AmplitudePercent),
		AmbianceName: decl.AmbianceName,
		NoiseSmooth:  decl.NoiseSmooth,
		Waveform:     decl.Waveform,
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