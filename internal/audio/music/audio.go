// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package music

import (
	"github.com/synapseq-foundation/synapseq/v4/internal/audio/audiosource"
	r "github.com/synapseq-foundation/synapseq/v4/internal/resource"
)

type Audio = audiosource.Audio

func NewAudio(filePaths []string, expectedSampleRate int) (*Audio, error) {
	return audiosource.New(filePaths, expectedSampleRate, audiosource.Options{
		PlaybackMode: audiosource.PlaybackFinite,
		LoadFile:     r.GetMusicFile,
		SourceKind:   "music",
	})
}
