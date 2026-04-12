//go:build wasm

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

package sequence

import (
	"fmt"

	r "github.com/synapseq-foundation/synapseq/v4/internal/resource"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func resolveParsedOptions(_ string, parsedOptions *t.ParseOptions) error {
	for _, path := range parsedOptions.Ambiance {
		if !r.IsRemoteFile(path) {
			return fmt.Errorf("WASM only supports remote URLs for ambiance audio")
		}
	}

	for _, path := range parsedOptions.Extends {
		if !r.IsRemoteFile(path) {
			return fmt.Errorf("WASM only supports remote URLs for extends")
		}
	}

	return nil
}