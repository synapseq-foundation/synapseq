// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package types

// AIProvider identifies an internal AI provider behavior.
type AIProvider string

const (
	// AIProviderDefault selects the standard provider behavior.
	AIProviderDefault         AIProvider = "default"
	// AIProviderAppleFoundation selects Apple's fm serve behavior.
	AIProviderAppleFoundation AIProvider = "apple-foundation"
)
