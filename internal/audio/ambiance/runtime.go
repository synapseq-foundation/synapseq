// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package ambiance

import (
	"github.com/synapseq-foundation/synapseq/v4/internal/audio/audiosource"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type SampleAudio = audiosource.SampleAudio
type Runtime = audiosource.Runtime

func NewRuntime(periods []t.Period, ambiance map[string]string, sampleRate int, newAudio audiosource.NewAudioFunc) (*Runtime, error) {
	return audiosource.NewRuntime(periods, ambiance, sampleRate, audiosource.RuntimeOptions{
		TrackType:  t.TrackAmbiance,
		SourceKind: "ambiance",
		Scope:      audiosource.BufferScopeSource,
		NewAudio:   newAudio,
	})
}

func NewTestRuntime(sampleCount int) *Runtime {
	return audiosource.NewTestRuntime(sampleCount)
}
