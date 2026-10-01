// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sbg

import (
	"fmt"
	"strings"
)

func parseDefinition(name, voicesText string) (nameDef, error) {
	definition := nameDef{name: name}
	for _, token := range strings.Fields(voicesText) {
		parsedVoice, err := parseVoice(token)
		if err != nil {
			return nameDef{}, fmt.Errorf("voice %q: %v", token, err)
		}
		definition.voices = append(definition.voices, parsedVoice)
	}
	if len(definition.voices) == 0 {
		return nameDef{}, fmt.Errorf("NameDef %q has no voices", name)
	}
	return definition, nil
}
