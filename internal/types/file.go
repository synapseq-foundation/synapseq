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

package types

const (
	// MaxTextFileSize is the maximum allowed size for text files (32KB)
	MaxTextFileSize = 32 * 1024
	// MaxBackgroundFileSize is the maximum allowed size for background files (20MB)
	MaxBackgroundFileSize = 20 * 1024 * 1024
	// MaxStructuredFileSize is the maximum allowed size for structured files (128KB)
	MaxStructuredFileSize = 128 * 1024
)

// FileFormat represents the format of the input/output file
type FileFormat int

const (
	FormatText FileFormat = iota
	FormatJSON
	FormatXML
	FormatYAML
	FormatWAV
	FormatUnknown
)

// String returns the string representation of the FileFormat
func (ff FileFormat) String() string {
	switch ff {
	case FormatText:
		return "text"
	case FormatJSON:
		return "json"
	case FormatXML:
		return "xml"
	case FormatYAML:
		return "yaml"
	case FormatWAV:
		return "wav"
	default:
		return "unknown"
	}
}
