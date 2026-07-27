// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sbg

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func parseVoice(token string) (voice, error) {
	switch {
	case token == "-":
		return voice{kind: voiceOff}, nil
	case strings.HasPrefix(token, "pink/"):
		amplitude, err := parseDecimal(strings.TrimPrefix(token, "pink/"))
		return voice{kind: voicePink, amplitude: amplitude}, err
	case token == "mix" || strings.HasPrefix(token, "mix/"):
		amplitude, err := parseDecimal(strings.TrimPrefix(token, "mix/"))
		return voice{kind: voiceMix, amplitude: amplitude}, err
	case strings.HasPrefix(token, "spin:"):
		width, _, beat, amplitude, err := parseTone(strings.TrimPrefix(token, "spin:"), true)
		return voice{kind: voiceSpin, width: width, beat: beat, amplitude: amplitude}, err
	default:
		carrier, hasBeat, beat, amplitude, err := parseTone(token, false)
		kind := voiceTone
		if hasBeat {
			kind = voiceBinaural
		}
		return voice{kind: kind, carrier: carrier, beat: beat, amplitude: amplitude}, err
	}
}

func parseTone(token string, requireBeat bool) (float64, bool, float64, float64, error) {
	toneText, amplitudeText, found := strings.Cut(token, "/")
	if !found || toneText == "" || amplitudeText == "" || strings.Contains(amplitudeText, "/") {
		return 0, false, 0, 0, errors.New("expected carrier[+|-beat]/amplitude")
	}
	amplitude, err := parseDecimal(amplitudeText)
	if err != nil {
		return 0, false, 0, 0, fmt.Errorf("invalid amplitude: %w", err)
	}

	signIndex := -1
	for index := 1; index < len(toneText); index++ {
		if toneText[index] == '+' || toneText[index] == '-' {
			if signIndex >= 0 {
				return 0, false, 0, 0, errors.New("multiple beat signs")
			}
			signIndex = index
		}
	}
	if signIndex < 0 {
		if requireBeat {
			return 0, false, 0, 0, errors.New("missing beat frequency")
		}
		carrier, err := parseDecimal(toneText)
		return carrier, false, 0, amplitude, err
	}
	if signIndex == len(toneText)-1 {
		return 0, false, 0, 0, errors.New("missing beat frequency")
	}
	carrier, err := parseDecimal(toneText[:signIndex])
	if err != nil {
		return 0, false, 0, 0, fmt.Errorf("invalid carrier: %w", err)
	}
	beat, err := parseDecimal(toneText[signIndex+1:])
	if err != nil {
		return 0, false, 0, 0, fmt.Errorf("invalid beat: %w", err)
	}
	return carrier, true, beat, amplitude, nil
}

func parseDecimal(value string) (float64, error) {
	if value == "" {
		return 0, errors.New("empty number")
	}
	for index, char := range value {
		if (char < '0' || char > '9') && char != '.' {
			return 0, fmt.Errorf("unexpected character %q", char)
		}
		if char == '.' && strings.Contains(value[index+1:], ".") {
			return 0, errors.New("more than one decimal point")
		}
	}
	number, err := strconv.ParseFloat(value, 64)
	if err != nil || math.IsInf(number, 0) || math.IsNaN(number) {
		return 0, fmt.Errorf("invalid number %q", value)
	}
	return number, nil
}
