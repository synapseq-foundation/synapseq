// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
			tr0.SourceName = tr2.SourceName
		}

		if tr2.Type == t.TrackSilence {
			tr2.Carrier = tr1.Carrier
			tr2.Resonance = tr1.Resonance
			tr2.NoiseSmooth = tr1.NoiseSmooth
			tr2.Effect.Intensity = tr1.Effect.Intensity
			tr2.Effect.Value = tr1.Effect.Value
			tr2.Effect.Type = tr1.Effect.Type
			tr2.SourceName = tr1.SourceName
			tr2.Waveform = tr1.Waveform
		}

		last.CrossfadeOut[ch] = t.TrackCrossfade{}
		next.CrossfadeIn[ch] = t.TrackCrossfade{}
		prepareEffectOnOffTransition(tr0, tr1, tr2)
		if requiresBoundaryCrossfade(*tr1, *tr2) {
			if isActiveTrack(*tr1) {
				last.CrossfadeOut[ch] = t.TrackCrossfade{Active: true, Track: *tr1}
				*tr1 = *lastCrossfadeSteadyTrack(tr0, tr1)
			}
			if isActiveTrack(*tr2) {
				next.CrossfadeIn[ch] = t.TrackCrossfade{Active: true, Track: *tr2}
			}
			continue
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
		tr1.SourceName = tr2.SourceName
	}
	return nil
}

func requiresBoundaryCrossfade(last, next t.Track) bool {
	if isSilenceTrack(last) || isSilenceTrack(next) {
		return false
	}
	if isActiveTrack(last) != isActiveTrack(next) {
		return true
	}
	if !isActiveTrack(last) || !isActiveTrack(next) {
		return false
	}
	if isEffectOnOffTransition(last, next) {
		return false
	}
	return last.Type != next.Type ||
		last.Effect.Type != next.Effect.Type ||
		last.SourceName != next.SourceName
}

func prepareEffectOnOffTransition(start, end, next *t.Track) {
	if !isEffectOnOffTransition(*end, *next) {
		return
	}
	if end.Effect.Type == t.EffectOff {
		start.Effect.Type = next.Effect.Type
		start.Effect.Value = next.Effect.Value
		start.Effect.Intensity = 0
		return
	}
	next.Effect.Type = end.Effect.Type
	next.Effect.Value = end.Effect.Value
	next.Effect.Intensity = 0
}

func isEffectOnOffTransition(last, next t.Track) bool {
	if !isActiveTrack(last) || !isActiveTrack(next) {
		return false
	}
	if last.Type != next.Type || last.SourceName != next.SourceName {
		return false
	}
	return (last.Effect.Type == t.EffectOff && next.Effect.Type != t.EffectOff) ||
		(last.Effect.Type != t.EffectOff && next.Effect.Type == t.EffectOff)
}

func isActiveTrack(track t.Track) bool {
	return track.Type != t.TrackOff && track.Type != t.TrackSilence
}

func isSilenceTrack(track t.Track) bool {
	return track.Type == t.TrackSilence
}

func lastCrossfadeSteadyTrack(start, end *t.Track) *t.Track {
	if isActiveTrack(*start) {
		return start
	}
	return end
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
