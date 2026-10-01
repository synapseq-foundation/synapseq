// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sbg

import (
	"strconv"
	"strings"
)

const maxPresetNameLength = 20

func validName(value string) bool {
	if value == "" || !asciiLetter(value[0]) {
		return false
	}
	for index := 1; index < len(value); index++ {
		char := value[index]
		isDigit := char >= '0' && char <= '9'
		isSeparator := strings.ContainsRune("_.+-", rune(char))
		if !asciiLetter(char) && !isDigit && !isSeparator {
			return false
		}
	}
	return true
}

func asciiLetter(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func convertedNames(definitions []nameDef) map[string]string {
	result := make(map[string]string, len(definitions))
	used := map[string]struct{}{"silence": {}}
	for _, definition := range definitions {
		var name strings.Builder
		for _, char := range strings.ToLower(definition.name) {
			isLetter := char >= 'a' && char <= 'z'
			isDigit := char >= '0' && char <= '9'
			isSeparator := char == '_' || char == '-'
			if isLetter || isDigit || isSeparator {
				name.WriteRune(char)
				continue
			}
			name.WriteByte('-')
		}
		base := name.String()
		if len(base) > maxPresetNameLength {
			base = base[:maxPresetNameLength]
		}
		candidate := base
		for suffix := 2; ; suffix++ {
			if _, exists := used[candidate]; !exists {
				break
			}
			ending := "-" + strconv.Itoa(suffix)
			candidate = base
			if len(candidate)+len(ending) > maxPresetNameLength {
				candidate = candidate[:maxPresetNameLength-len(ending)]
			}
			candidate += ending
		}
		used[candidate] = struct{}{}
		result[definition.name] = candidate
	}
	return result
}
