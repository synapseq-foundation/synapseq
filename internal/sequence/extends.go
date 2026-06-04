// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sequence

import (
	"fmt"
	"path/filepath"

	r "github.com/synapseq-foundation/synapseq/v4/internal/resource"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// extends loads preset and option definitions from a .spsc file.
func extends(fileName string) (*t.Extends, error) {
	rawContent, err := r.GetFile(fileName, t.FormatText)
	if err != nil {
		return nil, err
	}

	absInputFile, err := filepath.Abs(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve absolute path: %w", err)
	}

	baseDir := filepath.Dir(absInputFile)

	return parseExtendsContent(rawContent, absInputFile, baseDir)
}
