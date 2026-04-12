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

package timeline

import (
	"fmt"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// AdjustPeriods adjusts the tracks in the overlapping periods.
func AdjustPeriods(last, next *t.Period) error {
	if err := validatePeriodSteps(last, next); err != nil {
		return err
	}

	for ch := range t.NumberOfChannels {
		tr0 := &last.TrackStart[ch]
		tr1 := &last.TrackEnd[ch]
		tr2 := &next.TrackStart[ch]

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

		if (tr1.Type != t.TrackOff && tr1.Type != t.TrackSilence && tr2.Type == t.TrackOff) ||
			(tr1.Type == t.TrackOff && tr2.Type != t.TrackOff && tr2.Type != t.TrackSilence) {
			return diag.Validation(
				fmt.Sprintf("channel %d cannot switch directly between an active track and off on consecutive timeline entries: %s -> %s", ch+1, tr1.Type.String(), tr2.Type.String()),
			).WithHint("the current timeline entry conflicts with the previous one on this channel; insert a silence entry between them")
		}

		if tr1.Type != t.TrackOff &&
			tr1.Type != t.TrackSilence &&
			tr2.Type != t.TrackOff &&
			tr2.Type != t.TrackSilence {
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

func validatePeriodSteps(last, next *t.Period) error {
	if last == nil || next == nil {
		return nil
	}
	if last.Steps < 0 {
		return diag.Validation(fmt.Sprintf("timeline %s cannot use negative steps: %d", last.TimeString(), last.Steps))
	}

	durationMs := next.Time - last.Time
	maxSteps := MaxPeriodSteps(durationMs)
	if last.Steps > maxSteps {
		return diag.Validation(
			fmt.Sprintf("timeline %s uses %d steps but the maximum for this %d-second interval is %d", last.TimeString(), last.Steps, durationMs/1000, maxSteps),
		).WithHint("reduce steps or increase the time before the next timeline entry")
	}

	return nil
}
