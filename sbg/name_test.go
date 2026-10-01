// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sbg

import "testing"

func TestConvertedNamesPreserveCollisionOrder(t *testing.T) {
	tests := []struct{ source, want string }{
		{
			source: "Alpha.Beta",
			want:   "alpha-beta",
		},
		{
			source: "Alpha+Beta",
			want:   "alpha-beta-2",
		},
		{
			source: "alpha-beta",
			want:   "alpha-beta-3",
		},
		{
			source: "silence",
			want:   "silence-2",
		},
		{
			source: "Silence",
			want:   "silence-3",
		},
		{
			source: "ABCDEFGHIJKLMNOPQRSTUV",
			want:   "abcdefghijklmnopqrst",
		},
		{
			source: "abcdefghijklmnopqrstzz",
			want:   "abcdefghijklmnopqr-2",
		},
		{
			source: "Keep_Underscore",
			want:   "keep_underscore",
		},
	}
	definitions := make([]nameDef, 0, len(tests))
	for _, test := range tests {
		definitions = append(definitions, nameDef{name: test.source})
	}
	names := convertedNames(definitions)
	for _, test := range tests {
		if names[test.source] != test.want {
			t.Errorf("name %q = %q, want %q", test.source, names[test.source], test.want)
		}
	}
}
