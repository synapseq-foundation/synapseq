/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

package shared

import "fmt"

// IsValidNamedRef checks if a name is valid for a named reference
func IsValidNamedRef(name string) error {
	isLetter := func(b byte) bool {
		return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
	}
	isDigit := func(b byte) bool {
		return b >= '0' && b <= '9'
	}

	if len(name) == 0 {
		return fmt.Errorf("reference name cannot be empty")
	}

	first := name[0]
	if !isLetter(first) {
		return fmt.Errorf("reference name must start with a letter: %q", name)
	}

	for i := 1; i < len(name); i++ {
		ch := name[i]
		if !(isLetter(ch) || isDigit(ch) || ch == '_' || ch == '-') {
			return fmt.Errorf("invalid character in reference name %q: %q", name, string(ch))
		}
	}

	return nil
}
