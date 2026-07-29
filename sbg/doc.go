// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

/*
Package sbg converts supported SBaGen .sbg sequences into validated SynapSeq
sequences.

# Overview

New receives the caller's core.AppContext and returns a Converter. Its LoadFile
method converts a local SBaGen file, while LoadContent converts SBaGen text
already held in memory. Both methods return a core.LoadedContext, so callers
can inspect the generated sequence, retrieve its SPSQ representation through
RawContent, or render it with the regular core API.

# Example Usage

	converter, err := sbg.New(synapseq.NewAppContext())
	if err != nil {
		log.Fatal(err)
	}
	loaded, err := converter.LoadFile("session.sbg")
	if err != nil {
		log.Fatal(err)
	}
	if err := loaded.WAV("session.wav"); err != nil {
		log.Fatal(err)
	}

# Conversion Notes

Supported SBaGen voices are mapped to the corresponding SPSQ tracks. Binaural
sign orientation is not represented by SPSQ. Spin voices are approximated with
pink noise and a pan effect. An SBaGen -m source becomes the SPSQ music source
named music, and each mix voice becomes a music track. Music paths must remain
within the current directory and cannot contain parent traversal; their
extension is removed so SynapSeq can resolve an MP3 or WAV source. SBaGen fade
markers are accepted but converted to steady transitions.
*/
package sbg
