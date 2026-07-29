// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sbg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func parseMusicOption(fields []string) (string, bool, error) {
	for index, field := range fields {
		if field != "-m" {
			continue
		}
		if index+1 >= len(fields) {
			return "", true, errors.New("-m requires an audio file")
		}
		return fields[index+1], true, nil
	}
	return "", false, nil
}

func parseSampleRateOption(fields []string) (int, bool, error) {
	for index, field := range fields {
		if field != "-r" {
			continue
		}
		if index+1 >= len(fields) {
			return 0, true, errors.New("-r requires a sample rate")
		}
		sampleRate, err := strconv.Atoi(fields[index+1])
		if err != nil {
			return 0, true, errors.New("-r requires an integer sample rate")
		}
		return sampleRate, true, nil
	}
	return 0, false, nil
}

func validName(value string) bool {
	if value == "" || !asciiLetter(value[0]) {
		return false
	}
	for index := 1; index < len(value); index++ {
		char := value[index]
		if !asciiLetter(char) && (char < '0' || char > '9') && !strings.ContainsRune("_.+-", rune(char)) {
			return false
		}
	}
	return true
}

func asciiLetter(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func lineError(source string, line int, message string) error {
	return fmt.Errorf("parse %q line %d: %s", source, line, message)
}
