/*
Package external provides integrations with external audio utilities
such as ffplay and ffmpeg to extend SynapSeq output and playback workflows.

# Overview

The package uses SynapSeq's real-time PCM streaming (AppContext.Stream)
and sends it directly to the stdin of ffplay/ffmpeg. This avoids temporary
files, reduces memory usage, and keeps startup fast.

# External Utilities

The following utilities are supported:

  - ffplay – real-time playback of SynapSeq-generated audio
  - ffmpeg – MP3 encoding (CBR 320kbps) from streamed PCM input

Custom paths may be provided when constructing FFplay or FFmpeg.
If no path is given, the package attempts to locate the utility
using the system PATH.

# Example: Real-Time Playback

This example shows how to play a SynapSeq sequence directly through
ffplay using streaming PCM audio.

	package main

	import (
	    "log"

	    synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	    "github.com/synapseq-foundation/synapseq/v4/external"
	)

	func main() {
	    ctx, err := synapseq.NewAppContext("input.spsq", "")
	    if err != nil {
	        log.Fatal(err)
	    }

	    if err := ctx.LoadSequence(); err != nil {
	        log.Fatal(err)
	    }

	    player, err := external.NewFFPlay("")
	    if err != nil {
	        log.Fatal(err)
	    }

	    if err := player.Play(ctx); err != nil {
	        log.Fatal(err)
	    }
	}

# Example: MP3 Encoding

This example streams PCM audio to ffmpeg and saves it as an MP3 file
using constant bit rate (CBR) at 320 kbps.

	package main

	import (
	    "log"

	    synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	    "github.com/synapseq-foundation/synapseq/v4/external"
	)

	func main() {
	    ctx, err := synapseq.NewAppContext("input.spsq", "output.mp3")
	    if err != nil {
	        log.Fatal(err)
	    }

	    if err := ctx.LoadSequence(); err != nil {
	        log.Fatal(err)
	    }

	    encoder, err := external.NewFFmpeg("")
	    if err != nil {
	        log.Fatal(err)
	    }

	    if err := encoder.Convert(ctx, "mp3"); err != nil {
	        log.Fatal(err)
	    }
	}

# Audio Format Support

The Convert method currently supports:

  - MP3: uses libmp3lame encoder with CBR at 320 kbps

# Error Handling

If an external tool does not exist or is not executable, constructors
(NewFFPlay, NewFFmpeg) return an error. If the tool exits with a non-zero
status code, the returned error includes stream/process failure details.

# Platform Notes

  - On Linux/macOS, executable permission bits are checked.
  - On Windows, lookups rely on PATH and associated .exe resolution.
  - Streaming uses stdin pipes and does not rely on temporary files.

# More Information

Full documentation and examples are available at:
https://github.com/synapseq-foundation/synapseq
*/
package external
