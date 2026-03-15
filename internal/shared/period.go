/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
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
	"fmt"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// AdjustPeriods adjusts the tracks in the overlapping periods
func AdjustPeriods(last, next *t.Period) error {
	for ch := range t.NumberOfChannels {
		tr0 := &last.TrackStart[ch]
		tr1 := &last.TrackEnd[ch]
		tr2 := &next.TrackStart[ch]

		// Apply Fade-In
		if tr0.Type == t.TrackSilence {
			tr0.Type = tr2.Type
			tr0.Carrier = tr2.Carrier
			tr0.Resonance = tr2.Resonance
			tr0.Amplitude = 0
			tr0.NoiseSmooth = tr2.NoiseSmooth
			tr0.Waveform = tr2.Waveform
			tr0.Effect.Type = tr2.Effect.Type
			tr0.Effect.Value = tr2.Effect.Value
			tr0.Effect.Intensity = tr2.Effect.Intensity
			tr0.AmbianceName = tr2.AmbianceName
		}

		// Apply Fade-Out
		if tr2.Type == t.TrackSilence {
			tr2.Carrier = tr1.Carrier
			tr2.Resonance = tr1.Resonance
			tr2.NoiseSmooth = tr1.NoiseSmooth
			tr2.Effect.Intensity = tr1.Effect.Intensity
			tr2.Effect.Value = tr1.Effect.Value
			tr2.Effect.Type = tr1.Effect.Type
			tr2.AmbianceName = tr1.AmbianceName
			tr2.Waveform = tr1.Waveform
		}

		// Validate if previus period has a track on and next period turn it off or vice-versa
		if (tr1.Type != t.TrackOff && tr1.Type != t.TrackSilence && tr2.Type == t.TrackOff) ||
			(tr1.Type == t.TrackOff && tr2.Type != t.TrackOff && tr2.Type != t.TrackSilence) {
			return diag.Validation(
				fmt.Sprintf("channel %d cannot switch directly between an active track and off on consecutive timeline entries: %s -> %s", ch+1, tr1.Type.String(), tr2.Type.String()),
			).WithHint("the current timeline entry conflicts with the previous one on this channel; insert a silence entry between them")
		}

		// Determine if both periods have a track on
		if tr1.Type != t.TrackOff &&
			tr1.Type != t.TrackSilence &&
			tr2.Type != t.TrackOff &&
			tr2.Type != t.TrackSilence {
			// No slide allowed between different track types, effect types, or ambiance sources.
			if tr1.Type != tr2.Type {
				return diag.Validation(
					fmt.Sprintf("channel %d cannot move directly from %s to %s on consecutive timeline entries", ch+1, tr1.Type.String(), tr2.Type.String()),
				).WithHint("the current timeline entry reuses this channel with an incompatible track type; insert a silence entry between the two presets")
			}
			if tr1.Effect.Type != tr2.Effect.Type {
				return diag.Validation(
					fmt.Sprintf("channel %d cannot move directly from %s effect to %s effect on consecutive timeline entries", ch+1, tr1.Effect.Type.String(), tr2.Effect.Type.String()),
				).WithHint("the current timeline entry reuses this channel with an incompatible effect; insert a silence entry between the two presets")
			}
			if tr1.AmbianceName != tr2.AmbianceName {
				return diag.Validation(
					fmt.Sprintf("channel %d cannot move directly from ambiance %q to %q on consecutive timeline entries", ch+1, tr1.AmbianceName, tr2.AmbianceName),
				).WithHint("the current timeline entry reuses this channel with a different ambiance source; insert a silence entry between the two presets")
			}
		}

		// Carry forward the track settings from the end of the last period to the start of the next period
		tr1.Type = tr2.Type
		tr1.Effect.Type = tr2.Effect.Type
		tr1.Effect.Value = tr2.Effect.Value
		tr1.Carrier = tr2.Carrier
		tr1.Resonance = tr2.Resonance
		tr1.Amplitude = tr2.Amplitude
		tr1.NoiseSmooth = tr2.NoiseSmooth
		tr1.Effect.Intensity = tr2.Effect.Intensity
		tr1.Waveform = tr2.Waveform
		tr1.AmbianceName = tr2.AmbianceName
	}
	return nil
}
