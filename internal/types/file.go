// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package types

const (
	// MaxTextFileSize is the maximum allowed size for text files (32KB)
	MaxTextFileSize = 32 * 1024
	// MaxAmbianceFileSize is the maximum allowed size for ambiance audio files (20MB)
	MaxAmbianceFileSize = 20 * 1024 * 1024
	// MaxMusicFileSize is the maximum allowed size for music audio files (80MB)
	MaxMusicFileSize = 80 * 1024 * 1024
)

// FileFormat represents the format of the input/output file
type FileFormat int

const (
	FormatText FileFormat = iota
	FormatAmbiance
)

// String returns the string representation of the FileFormat
func (ff FileFormat) String() string {
	switch ff {
	case FormatText:
		return "text"
	case FormatAmbiance:
		return "ambiance"
	default:
		return "unknown"
	}
}

// AmbianceAudioFormat represents a supported ambiance source audio format.
type AmbianceAudioFormat int

const (
	AmbianceAudioUnknown AmbianceAudioFormat = iota
	AmbianceAudioWAV
	AmbianceAudioMP3
)

func (aaf AmbianceAudioFormat) String() string {
	switch aaf {
	case AmbianceAudioWAV:
		return "wav"
	case AmbianceAudioMP3:
		return "mp3"
	default:
		return "unknown"
	}
}
