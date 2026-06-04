// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package music

import (
	"github.com/synapseq-foundation/synapseq/v4/internal/audio/audiosource"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type SampleAudio = audiosource.SampleAudio
type Runtime = audiosource.Runtime

func NewRuntime(periods []t.Period, music map[string]string, sampleRate int, newAudio audiosource.NewAudioFunc) (*Runtime, error) {
	return audiosource.NewRuntime(periods, music, sampleRate, audiosource.RuntimeOptions{
		TrackType:  t.TrackMusic,
		SourceKind: "music",
		Scope:      audiosource.BufferScopeChannel,
		NewAudio:   newAudio,
	})
}
