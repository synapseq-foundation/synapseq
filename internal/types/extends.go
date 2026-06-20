// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package types

// Extends represents the configuration for options and presets to extend from.
type Extends struct {
	Presets []Preset
	Options *ParseOptions
}
