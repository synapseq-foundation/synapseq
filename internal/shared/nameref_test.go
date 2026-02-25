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

import (
	"strings"
	"testing"
)

func TestIsValidNamedRef_Valid(t *testing.T) {
	tests := []string{
		"rain",
		"Rain",
		"rain_01",
		"rain-01",
		"A1_b-2",
		"abcdefghijklmnopqrst", // MaxNamedRefLength (20)
	}

	for _, name := range tests {
		if err := IsValidNamedRef(name); err != nil {
			t.Errorf("expected valid name %q, got error: %v", name, err)
		}
	}
}

func TestIsValidNamedRef_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		errContains string
	}{
		{"", "cannot be empty"},
		{"1rain", "must start with a letter"},
		{"_rain", "must start with a letter"},
		{"-rain", "must start with a letter"},
		{"rain name", "invalid character"},
		{"rain.name", "invalid character"},
		{"rain@name", "invalid character"},
		{"abcdefghijklmnopqrstu", "cannot be longer"}, // 21 chars
	}

	for _, tt := range tests {
		err := IsValidNamedRef(tt.name)
		if err == nil {
			t.Errorf("expected error for invalid name %q, got nil", tt.name)
			continue
		}

		if !strings.Contains(err.Error(), tt.errContains) {
			t.Errorf("invalid name %q: expected error containing %q, got %q", tt.name, tt.errContains, err.Error())
		}
	}
}
