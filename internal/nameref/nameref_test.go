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

package nameref

import (
	"strings"
	"testing"
)

func TestIsValid_Valid(ts *testing.T) {
	tests := []string{
		"rain",
		"Rain",
		"rain_01",
		"rain-01",
		"A1_b-2",
		"abcdefghijklmnopqrst",
	}

	for _, name := range tests {
		if err := IsValid(name); err != nil {
			ts.Errorf("expected valid name %q, got error: %v", name, err)
		}
	}
}

func TestIsValid_Invalid(ts *testing.T) {
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
		{"abcdefghijklmnopqrstu", "cannot be longer"},
	}

	for _, test := range tests {
		err := IsValid(test.name)
		if err == nil {
			ts.Errorf("expected error for invalid name %q, got nil", test.name)
			continue
		}

		if !strings.Contains(err.Error(), test.errContains) {
			ts.Errorf("invalid name %q: expected error containing %q, got %q", test.name, test.errContains, err.Error())
		}
	}
}