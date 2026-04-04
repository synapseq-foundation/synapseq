//go:build !wasm

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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func writeSeqFile(tst *testing.T, content string) string {
	tst.Helper()
	dir := tst.TempDir()
	p := filepath.Join(dir, "seq.spsq")
	if err := os.WriteFile(p, []byte(strings.TrimSpace(content)+"\n"), 0o600); err != nil {
		tst.Fatalf("write temp sequence: %v", err)
	}
	return p
}

func writeRelFile(tst *testing.T, dir, relPath, content string) string {
	tst.Helper()

	path := filepath.Join(dir, relPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		tst.Fatalf("mkdir temp rel file dir: %v", err)
	}

	if err := os.WriteFile(path, []byte(strings.TrimSpace(content)+"\n"), 0o600); err != nil {
		tst.Fatalf("write temp rel file: %v", err)
	}

	return path
}

func eqTrackGotWant(got, want t.Track) bool {
	if got.Type != want.Type {
		return false
	}
	if got.Carrier != want.Carrier {
		return false
	}
	if got.Resonance != want.Resonance {
		return false
	}
	if got.Amplitude != want.Amplitude {
		return false
	}
	if got.Waveform != want.Waveform {
		return false
	}
	if got.Effect.Type != want.Effect.Type {
		return false
	}
	if got.Effect.Intensity != want.Effect.Intensity {
		return false
	}
	return true
}

func hasTrack(tracks [t.NumberOfChannels]t.Track, want t.Track) bool {
	for ch := range t.NumberOfChannels {
		if eqTrackGotWant(tracks[ch], want) {
			return true
		}
	}
	return false
}

func TestLoadTextSequence_Success(ts *testing.T) {
	seq := `
# Options
@samplerate 48000
@volume 80
@ambiance testnoise testdata/noise

# Presets
alpha
  noise brown amplitude 40
  tone 300 binaural 10 amplitude 20

beta
  noise brown amplitude 40
  tone 300 binaural 14 amplitude 15

# Timeline
00:00:00 alpha
00:01:00 beta
`
	path := writeSeqFile(ts, seq)
	abPath := filepath.Join(filepath.Dir(path), "testdata", "noise.wav")

	result, err := LoadTextSequence(path)
	if err != nil {
		ts.Fatalf("LoadTextSequence error: %v", err)
	}

	opts := result.Options
	if opts.SampleRate != 48000 || opts.Volume != 80 {
		ts.Fatalf("unexpected options: %+v", *opts)
	}
	if opts.Ambiance["testnoise"] != abPath {
		ts.Fatalf("unexpected ambiance path: got %q want %q", opts.Ambiance["testnoise"], abPath)
	}

	periods := result.Periods
	if len(periods) != 2 {
		ts.Fatalf("expected 2 periods, got %d", len(periods))
	}
	if periods[0].Time != 0 || periods[1].Time != 60_000 {
		ts.Fatalf("unexpected period times: %+v", []int{periods[0].Time, periods[1].Time})
	}

	wantNoise := t.Track{
		Type:      t.TrackBrownNoise,
		Amplitude: t.AmplitudePercentToRaw(40),
	}
	wantToneAlpha := t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   300,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(20),
	}
	wantToneBeta := t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   300,
		Resonance: 14,
		Amplitude: t.AmplitudePercentToRaw(15),
	}

	if !hasTrack(periods[0].TrackStart, wantNoise) || !hasTrack(periods[0].TrackStart, wantToneAlpha) {
		ts.Fatalf("missing alpha tracks in period[0]: %+v", periods[0].TrackStart)
	}
	if !hasTrack(periods[1].TrackStart, wantToneBeta) {
		ts.Fatalf("missing beta tracks in period[1]: %+v", periods[1].TrackStart)
	}
}

func TestLoadTextSequence_Error_TimelineBeforePreset(ts *testing.T) {
	seq := `
# Timeline first
00:00:00 alpha
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for timeline before preset")
	}
}

func TestLoadTextSequence_Error_OptionsAfterPreset(ts *testing.T) {
	seq := `
alpha
  tone 100 binaural 1 amplitude 1
@volume 80
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for options after preset")
	}
}

func TestLoadTextSequence_Error_FirstTimelineNotZero(ts *testing.T) {
	seq := `
alpha
  tone 100 binaural 1 amplitude 1
00:00:15 alpha
00:01:00 alpha
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for first timeline not 00:00:00")
	}
}

func TestLoadTextSequence_Error_OverlappingTimeline(ts *testing.T) {
	seq := `
alpha
  tone 100 binaural 1 amplitude 1
00:00:00 alpha
00:01:00 alpha
00:00:30 alpha
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for overlapping timeline")
	}
}

func TestLoadTextSequence_Error_StepsExceedDurationLimit(ts *testing.T) {
	seq := `
alpha
  tone 100 binaural 1 amplitude 1
beta
  tone 120 binaural 4 amplitude 2
00:00:00 alpha steady 3
00:00:30 beta
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for steps above duration limit")
	}
	if !strings.Contains(err.Error(), "uses 3 steps") {
		ts.Fatalf("expected steps validation error, got %v", err)
	}
}

func TestLoadTextSequence_Error_BackgroundWithoutFile(ts *testing.T) {
	seq := `
alpha
  background amplitude 10
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for background track without background option")
	}
}

func TestLoadTextSequence_TrackDiagnosticIncludesSource(ts *testing.T) {
	seq := `
alpha
  tone 300 binaual 10 amplitude 20
00:00:00 alpha
00:01:00 alpha
`
	path := writeSeqFile(ts, seq)

	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatal("expected parse error")
	}

	diagnostic, ok := diag.As(err)
	if !ok {
		ts.Fatalf("expected diag.Diagnostic, got %T", err)
	}
	if diagnostic.Span.File != path {
		ts.Fatalf("expected source file %q, got %q", path, diagnostic.Span.File)
	}
	if diagnostic.Span.Line != 2 {
		ts.Fatalf("expected source line 2, got %d", diagnostic.Span.Line)
	}
	if diagnostic.Span.Column != 12 || diagnostic.Span.EndColumn != 19 {
		ts.Fatalf("expected typo at 12..19, got %d..%d", diagnostic.Span.Column, diagnostic.Span.EndColumn)
	}
	if diagnostic.Suggestion != "did you mean \"binaural\"?" {
		ts.Fatalf("expected typo suggestion, got %q", diagnostic.Suggestion)
	}
}

func TestLoadTextSequence_Error_DuplicatePreset(ts *testing.T) {
	seq := `
alpha
alpha
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for duplicate preset")
	}
}

func TestLoadTextSequence_Error_MaxPresets(ts *testing.T) {
	var b strings.Builder
	for i := 1; i <= t.MaxPresets-1; i++ {
		b.WriteString("p")
		b.WriteString(strings.Repeat("x", i%5+1))
		b.WriteString("\n")
	}
	b.WriteString("overflow\n")
	path := writeSeqFile(ts, b.String())
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for maximum number of presets reached")
	}
}

func TestLoadTextSequence_Comments(ts *testing.T) {
	seq := `
## This is a comment on at the top

# Presets
alpha
  tone 100 binaural 1 amplitude 1

## Another comment line

# Timeline
00:00:00 alpha
##Another comment line in between
00:01:00 alpha
`
	path := writeSeqFile(ts, seq)
	result, err := LoadTextSequence(path)
	if err != nil {
		ts.Fatalf("LoadTextSequence with comments error: %v", err)
	}

	cmms := result.Comments
	if cmms == nil || len(cmms) != 3 {
		ts.Fatalf("expected 2 comment blocks, got %d", len(cmms))
	}
	if cmms[0] != "This is a comment on at the top" {
		ts.Fatalf("unexpected first comment block: %q", cmms[0])
	}
	if cmms[1] != "Another comment line" {
		ts.Fatalf("unexpected second comment block: %q", cmms[1])
	}
	if cmms[2] != "Another comment line in between" {
		ts.Fatalf("unexpected third comment block: %q", cmms[2])
	}
}

func TestLoadTextSequence_Error_PresetEmpty(ts *testing.T) {
	seq := `
alpha
00:00:00 alpha
00:01:00 alpha
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for empty preset")
	}
}

func TestLoadTextSequence_Error_MultipleBackgrounds(ts *testing.T) {
	wd, _ := os.Getwd()
	bg := filepath.Join(wd, "testdata", "noise.wav")
	if _, err := os.Stat(bg); err != nil {
		ts.Fatalf("missing test bg file: %v", err)
	}
	seq := `
@background testdata/noise.wav
alpha
  background amplitude 20
  background amplitude 30
00:00:00 alpha
00:01:00 alpha
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for multiple background tracks in a preset")
	}
}

func TestLoadTextSequence_Error_AtLeastTwoPeriods(ts *testing.T) {
	seq := `
alpha
  tone 100 binaural 1 amplitude 1
00:00:00 alpha
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for less than two periods")
	}
}

func TestLoadTextSequence_LoadsExternalPresetViaExtends(ts *testing.T) {
	dir := ts.TempDir()

	presetRel := "external-presets.spsc"
	presetPath := filepath.Join(dir, presetRel)
	presetContent := `
# Presets
preparation
  noise pink amplitude 50
  tone 300 binaural 10 amplitude 10
`
	if err := os.WriteFile(presetPath, []byte(presetContent+"\n"), 0o600); err != nil {
		ts.Fatalf("write external presets: %v", err)
	}

	seqContent := "# Options\n@extends external-presets\n\n# Timeline\n00:00:00 preparation\n00:01:00 preparation\n"
	seqPath := filepath.Join(dir, "seq.spsq")
	if err := os.WriteFile(seqPath, []byte(strings.TrimSpace(seqContent)+"\n"), 0o600); err != nil {
		ts.Fatalf("write temp sequence: %v", err)
	}

	result, err := LoadTextSequence(seqPath)
	if err != nil {
		ts.Fatalf("LoadTextSequence with external extends error: %v", err)
	}

	if len(result.Periods) != 2 {
		ts.Fatalf("expected 2 periods, got %d", len(result.Periods))
	}

	wantNoise := t.Track{
		Type:      t.TrackPinkNoise,
		Amplitude: t.AmplitudePercentToRaw(50),
	}
	wantTone := t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   300,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(10),
	}

	if !hasTrack(result.Periods[0].TrackStart, wantNoise) || !hasTrack(result.Periods[0].TrackStart, wantTone) {
		ts.Fatalf("missing preparation tracks in period[0]: %+v", result.Periods[0].TrackStart)
	}
	if !hasTrack(result.Periods[1].TrackStart, wantNoise) || !hasTrack(result.Periods[1].TrackStart, wantTone) {
		ts.Fatalf("missing preparation tracks in period[1]: %+v", result.Periods[1].TrackStart)
	}
}

func TestLoadTextSequence_ExtendsOverrideOptionsByOrder(ts *testing.T) {
	dir := ts.TempDir()

	writeRelFile(ts, dir, "config/base.spsc", `
@samplerate 44100

alpha
  tone 100 binaural 1 amplitude 1
`)

	writeRelFile(ts, dir, "config/override.spsc", `
@samplerate 48000
`)

	seqPath := writeRelFile(ts, dir, "seq.spsq", `
@extends config/base
@extends config/override

00:00:00 alpha
00:01:00 alpha
`)

	result, err := LoadTextSequence(seqPath)
	if err != nil {
		ts.Fatalf("LoadTextSequence with ordered extends error: %v", err)
	}

	if result.Options.SampleRate != 48000 {
		ts.Fatalf("expected SampleRate=48000 after ordered extends override, got %d", result.Options.SampleRate)
	}
}

func TestLoadTextSequence_MergesAmbianceFromExtendsAndMainFile(ts *testing.T) {
	dir := ts.TempDir()

	writeRelFile(ts, dir, "packs/forest.spsc", `
@ambiance forest audio/forest

alpha
  ambiance forest amplitude 20
`)

	writeRelFile(ts, dir, "packs/river.spsc", `
@ambiance river audio/river
`)

	seqPath := writeRelFile(ts, dir, "seq.spsq", `
@extends packs/forest
@extends packs/river
@ambiance wind audio/wind

00:00:00 alpha
00:01:00 alpha
`)

	result, err := LoadTextSequence(seqPath)
	if err != nil {
		ts.Fatalf("LoadTextSequence with merged ambiance error: %v", err)
	}

	want := map[string]string{
		"forest": filepath.Join(dir, "packs", "audio", "forest.wav"),
		"river":  filepath.Join(dir, "packs", "audio", "river.wav"),
		"wind":   filepath.Join(dir, "audio", "wind.wav"),
	}

	if len(result.Options.Ambiance) != len(want) {
		ts.Fatalf("expected %d ambiance entries, got %d: %+v", len(want), len(result.Options.Ambiance), result.Options.Ambiance)
	}

	for name, path := range want {
		if result.Options.Ambiance[name] != path {
			ts.Fatalf("expected ambiance %q to resolve to %q, got %q", name, path, result.Options.Ambiance[name])
		}
	}
}
