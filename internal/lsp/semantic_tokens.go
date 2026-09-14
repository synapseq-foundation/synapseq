// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package lsp

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
	"go.lsp.dev/protocol"
)

const (
	semanticTokenComment uint32 = iota
	semanticTokenKeyword
	semanticTokenNumber
	semanticTokenString
	semanticTokenParameter
	semanticTokenVariable
)

var semanticTokenLegend = []string{
	string(protocol.SemanticTokenTypesComment),
	string(protocol.SemanticTokenTypesKeyword),
	string(protocol.SemanticTokenTypesNumber),
	string(protocol.SemanticTokenTypesString),
	string(protocol.SemanticTokenTypesParameter),
	string(protocol.SemanticTokenTypesVariable),
}

var syntaxKeywords = map[string]struct{}{
	t.KeywordOff:     {},
	t.KeywordSilence: {},
	t.KeywordOption + t.KeywordOptionSampleRate: {},
	t.KeywordOption + t.KeywordOptionVolume:     {},
	t.KeywordOption + t.KeywordOptionAmbiance:   {},
	t.KeywordOption + t.KeywordOptionMusic:      {},
	t.KeywordOption + t.KeywordOptionExtends:    {},
	t.KeywordOption + t.KeywordOptionWaveform:   {},
	t.KeywordOption + t.KeywordOptionTransition: {},
	t.KeywordWaveform:                           {},
	t.KeywordSine:                               {},
	t.KeywordSquare:                             {},
	t.KeywordTriangle:                           {},
	t.KeywordSawtooth:                           {},
	t.KeywordTone:                               {},
	t.KeywordBinaural:                           {},
	t.KeywordMonaural:                           {},
	t.KeywordIsochronic:                         {},
	t.KeywordNoise:                              {},
	t.KeywordWhite:                              {},
	t.KeywordPink:                               {},
	t.KeywordBrown:                              {},
	t.KeywordAmbiance:                           {},
	t.KeywordMusic:                              {},
	t.KeywordPure:                               {},
	t.KeywordTransitionSteady:                   {},
	t.KeywordTransitionEaseOut:                  {},
	t.KeywordTransitionEaseIn:                   {},
	t.KeywordTransitionSmooth:                   {},
	t.KeywordFrom:                               {},
	t.KeywordAs:                                 {},
	t.KeywordTemplate:                           {},
}

var syntaxParameters = map[string]struct{}{
	t.KeywordAmplitude:  {},
	t.KeywordLeft:       {},
	t.KeywordRight:      {},
	t.KeywordEffect:     {},
	t.KeywordPan:        {},
	t.KeywordModulation: {},
	t.KeywordDoppler:    {},
	t.KeywordShift:      {},
	t.KeywordIntensity:  {},
	t.KeywordSmooth:     {},
	t.KeywordTrack:      {},
}

type semanticToken struct {
	line      uint32
	start     uint32
	length    uint32
	tokenType uint32
}

type lexicalToken struct {
	value string
	start int
	end   int
}

func semanticTokenData(text string) []uint32 {
	tokens := semanticTokens(text)
	data := make([]uint32, 0, len(tokens)*5)
	var previousLine uint32
	var previousStart uint32
	for index, token := range tokens {
		deltaLine := token.line - previousLine
		deltaStart := token.start
		if index > 0 && deltaLine == 0 {
			deltaStart -= previousStart
		}
		data = append(data, deltaLine, deltaStart, token.length, token.tokenType, 0)
		previousLine = token.line
		previousStart = token.start
	}
	return data
}

func semanticTokens(text string) []semanticToken {
	result := []semanticToken{}
	for lineIndex, line := range strings.Split(text, "\n") {
		line = strings.TrimSuffix(line, "\r")
		lexical := scanSemanticTokens(line)
		if len(lexical) == 0 {
			continue
		}

		if strings.HasPrefix(lexical[0].value, t.KeywordComment) {
			result = append(result, makeSemanticToken(line, uint32(lineIndex), lexical[0].start, len(line), semanticTokenComment))
			continue
		}

		for tokenIndex, token := range lexical {
			tokenType := semanticTokenTypeForLine(lexical, tokenIndex)
			result = append(result, makeSemanticToken(line, uint32(lineIndex), token.start, token.end, tokenType))
		}
	}
	return result
}

func semanticTokenTypeForLine(tokens []lexicalToken, index int) uint32 {
	value := tokens[index].value
	if index > 0 && isOptionPath(tokens[0].value, index) {
		return semanticTokenString
	}
	return semanticTokenType(value)
}

func isOptionPath(option string, index int) bool {
	switch option {
	case t.KeywordOption + t.KeywordOptionExtends:
		return index == 1
	case t.KeywordOption + t.KeywordOptionAmbiance, t.KeywordOption + t.KeywordOptionMusic:
		return index == 2
	default:
		return false
	}
}

func scanSemanticTokens(line string) []lexicalToken {
	tokens := []lexicalToken{}
	for index := 0; index < len(line); {
		runeValue, size := utf8.DecodeRuneInString(line[index:])
		if unicode.IsSpace(runeValue) {
			index += size
			continue
		}

		start := index
		for index < len(line) {
			runeValue, size = utf8.DecodeRuneInString(line[index:])
			if unicode.IsSpace(runeValue) {
				break
			}
			index += size
		}
		tokens = append(tokens, lexicalToken{value: line[start:index], start: start, end: index})
	}
	return tokens
}

func makeSemanticToken(line string, lineIndex uint32, start, end int, tokenType uint32) semanticToken {
	utf16Start := utf16Offset(line, start)
	utf16End := utf16Offset(line, end)
	return semanticToken{
		line:      lineIndex,
		start:     uint32(utf16Start),
		length:    uint32(utf16End - utf16Start),
		tokenType: tokenType,
	}
}

func semanticTokenType(value string) uint32 {
	if _, ok := syntaxKeywords[value]; ok {
		return semanticTokenKeyword
	}
	if _, ok := syntaxParameters[value]; ok {
		return semanticTokenParameter
	}
	if isSemanticNumber(value) {
		return semanticTokenNumber
	}
	if isSemanticString(value) {
		return semanticTokenString
	}
	return semanticTokenVariable
}

func isSemanticNumber(value string) bool {
	if strings.Count(value, ":") == 2 {
		parts := strings.Split(value, ":")
		for _, part := range parts {
			if _, err := strconv.Atoi(part); err != nil {
				return false
			}
		}
		return true
	}
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

func isSemanticString(value string) bool {
	return strings.Contains(value, "/") || strings.Contains(value, "://") || strings.HasPrefix(value, ".")
}
