// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
