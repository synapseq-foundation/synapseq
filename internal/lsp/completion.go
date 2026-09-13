// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package lsp

import (
	"sort"
	"strings"
	"unicode/utf8"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
	"go.lsp.dev/protocol"
)

type symbols struct {
	presets     []string
	ambiance    []string
	music       []string
	waveforms   []string
	transitions []string
}

func complete(text string, position protocol.Position) protocol.CompletionItemSlice {
	line := lineAt(text, position.Line)
	prefix := linePrefix(line, position.Character)
	values := completionValues(prefix, collectSymbols(text))
	items := make(protocol.CompletionItemSlice, 0, len(values))
	for _, value := range values {
		items = append(items, protocol.CompletionItem{
			Label: value,
			Kind:  protocol.CompletionItemKindKeyword,
		})
	}
	return items
}

func completionValues(prefix string, index symbols) []string {
	fields := strings.Fields(prefix)
	trimmed := strings.TrimSpace(prefix)
	if strings.HasPrefix(trimmed, t.KeywordOption) && len(fields) <= 1 {
		return []string{
			t.KeywordOption + t.KeywordOptionSampleRate,
			t.KeywordOption + t.KeywordOptionVolume,
			t.KeywordOption + t.KeywordOptionAmbiance,
			t.KeywordOption + t.KeywordOptionMusic,
			t.KeywordOption + t.KeywordOptionWaveform,
			t.KeywordOption + t.KeywordOptionTransition,
			t.KeywordOption + t.KeywordOptionExtends,
		}
	}
	if strings.HasPrefix(prefix, "  ") {
		return trackCompletions(fields, index)
	}
	if len(fields) >= 1 && strings.Contains(fields[0], ":") {
		if len(fields) == 1 || len(fields) == 2 {
			return append([]string{t.KeywordSilence}, index.presets...)
		}
		return append(t.BuiltinTransitionNames(), index.transitions...)
	}
	return []string{}
}

func trackCompletions(fields []string, index symbols) []string {
	if len(fields) <= 1 {
		return []string{t.KeywordTone, t.KeywordNoise, t.KeywordAmbiance, t.KeywordMusic, t.KeywordWaveform, t.KeywordTrack}
	}
	if fields[0] == t.KeywordWaveform && len(fields) == 1 {
		return append([]string{t.KeywordSine, t.KeywordSquare, t.KeywordTriangle, t.KeywordSawtooth}, index.waveforms...)
	}
	if fields[0] == t.KeywordAmbiance && len(fields) <= 2 {
		return index.ambiance
	}
	if fields[0] == t.KeywordMusic && len(fields) <= 2 {
		return index.music
	}
	if fields[0] == t.KeywordNoise && len(fields) <= 2 {
		return []string{t.KeywordWhite, t.KeywordPink, t.KeywordBrown}
	}
	if fields[0] == t.KeywordTone && len(fields) >= 2 {
		return []string{t.KeywordBinaural, t.KeywordMonaural, t.KeywordIsochronic, t.KeywordEffect, t.KeywordAmplitude}
	}
	return []string{t.KeywordEffect, t.KeywordAmplitude, t.KeywordSmooth, t.KeywordIntensity}
}

func collectSymbols(text string) symbols {
	index := symbols{
		presets:     []string{},
		ambiance:    []string{},
		music:       []string{},
		waveforms:   []string{},
		transitions: []string{},
	}
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "  ") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 || strings.HasPrefix(fields[0], t.KeywordComment) {
			continue
		}
		if strings.HasPrefix(fields[0], t.KeywordOption) {
			if len(fields) < 3 {
				continue
			}
			switch fields[0] {
			case t.KeywordOption + t.KeywordOptionAmbiance:
				index.ambiance = append(index.ambiance, fields[1])
			case t.KeywordOption + t.KeywordOptionMusic:
				index.music = append(index.music, fields[1])
			case t.KeywordOption + t.KeywordOptionWaveform:
				index.waveforms = append(index.waveforms, fields[1])
			case t.KeywordOption + t.KeywordOptionTransition:
				index.transitions = append(index.transitions, fields[1])
			}
			continue
		}
		if !strings.Contains(fields[0], ":") {
			index.presets = append(index.presets, fields[0])
		}
	}
	sort.Strings(index.presets)
	sort.Strings(index.ambiance)
	sort.Strings(index.music)
	sort.Strings(index.waveforms)
	sort.Strings(index.transitions)
	return index
}

func lineAt(text string, target uint32) string {
	lines := strings.Split(text, "\n")
	if int(target) >= len(lines) {
		return ""
	}
	return strings.TrimSuffix(lines[target], "\r")
}

func linePrefix(line string, character uint32) string {
	byteOffset := 0
	units := uint32(0)
	for byteOffset < len(line) && units < character {
		runeValue, size := utf8.DecodeRuneInString(line[byteOffset:])
		if runeValue > 0xFFFF {
			units += 2
		} else {
			units++
		}
		byteOffset += size
	}
	return line[:byteOffset]
}
