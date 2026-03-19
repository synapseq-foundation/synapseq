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

package templates

import (
	"embed"
	"fmt"
)

//go:embed *.spsq
var files embed.FS

// GetTemplateContent returns the content of a template file by name
func GetTemplateContent(name string) ([]byte, error) {
	template, err := files.ReadFile(name + ".spsq")
	if err != nil {
		return nil, fmt.Errorf("template %q not found", name)
	}
	return template, nil
}
