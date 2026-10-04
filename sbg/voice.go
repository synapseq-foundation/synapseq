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

type parsedTone struct {
	carrier   float64
	hasBeat   bool
	beat      float64
	reverse   bool
	amplitude float64
}

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
	case token == "bell" ||
		strings.HasPrefix(token, "bell+") ||
		strings.HasPrefix(token, "bell-") ||
		strings.HasPrefix(token, "bell/"):
		return voice{}, errors.New("bell voices are not supported by SynapSeq")
	case strings.HasPrefix(token, "spin:"):
		tone, err := parseTone(strings.TrimPrefix(token, "spin:"), true)
		return voice{
			kind:      voiceSpin,
			carrier:   tone.carrier,
			beat:      tone.beat,
			amplitude: tone.amplitude,
		}, err
	default:
		tone, err := parseTone(token, false)
		kind := voiceTone
		if tone.hasBeat {
			kind = voiceBinaural
		}
		return voice{
			kind:      kind,
			reverse:   kind == voiceBinaural && tone.reverse,
			carrier:   tone.carrier,
			beat:      tone.beat,
			amplitude: tone.amplitude,
		}, err
	}
}

func parseTone(token string, requireBeat bool) (parsedTone, error) {
	toneText, amplitudeText, found := strings.Cut(token, "/")
	missingComponent := !found || toneText == "" || amplitudeText == ""
	if missingComponent || strings.Contains(amplitudeText, "/") {
		return parsedTone{}, errors.New("expected carrier[+|-beat]/amplitude")
	}
	amplitude, err := parseDecimal(amplitudeText)
	if err != nil {
		return parsedTone{}, fmt.Errorf("invalid amplitude: %w", err)
	}

	signIndex := -1
	for index := 1; index < len(toneText); index++ {
		if toneText[index] != '+' && toneText[index] != '-' {
			continue
		}
		if signIndex >= 0 {
			return parsedTone{}, errors.New("multiple beat signs")
		}
		signIndex = index
	}
	if signIndex < 0 {
		if requireBeat {
			return parsedTone{}, errors.New("missing beat frequency")
		}
		carrier, err := parseDecimal(toneText)
		return parsedTone{carrier: carrier, amplitude: amplitude}, err
	}
	if signIndex == len(toneText)-1 {
		return parsedTone{}, errors.New("missing beat frequency")
	}
	carrier, err := parseDecimal(toneText[:signIndex])
	if err != nil {
		return parsedTone{}, fmt.Errorf("invalid carrier: %w", err)
	}
	beat, err := parseDecimal(toneText[signIndex+1:])
	if err != nil {
		return parsedTone{}, fmt.Errorf("invalid beat: %w", err)
	}
	return parsedTone{
		carrier:   carrier,
		hasBeat:   true,
		reverse:   toneText[signIndex] == '-',
		beat:      beat,
		amplitude: amplitude,
	}, nil
}

func parseDecimal(value string) (float64, error) {
	if value == "" {
		return 0, errors.New("empty number")
	}
	seenDecimal := false
	for _, char := range value {
		if (char < '0' || char > '9') && char != '.' {
			return 0, fmt.Errorf("unexpected character %q", char)
		}
		if char != '.' {
			continue
		}
		if seenDecimal {
			return 0, errors.New("more than one decimal point")
		}
		seenDecimal = true
	}
	number, err := strconv.ParseFloat(value, 64)
	if err != nil || math.IsInf(number, 0) || math.IsNaN(number) {
		return 0, fmt.Errorf("invalid number %q", value)
	}
	return number, nil
}
