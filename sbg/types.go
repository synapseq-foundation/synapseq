// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sbg

import "time"

type voiceKind int

const (
	voiceOff voiceKind = iota
	voicePink
	voiceTone
	voiceBinaural
	voiceSpin
	voiceMix
)

type voice struct {
	kind      voiceKind
	carrier   float64
	beat      float64
	amplitude float64
	width     float64
}

type nameDef struct {
	name   string
	voices []voice
	line   int
}

type timelineEvent struct {
	at       time.Duration
	name     string
	initial  bool
	relative bool
	line     int
}

type sequence struct {
	source      string
	musicPath   string
	definitions []nameDef
	timeline    []timelineEvent
}
