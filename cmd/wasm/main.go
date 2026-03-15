//go:build wasm

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

package main

import (
	"fmt"
	"syscall/js"

	"github.com/synapseq-foundation/synapseq/v4/internal/audio"
	"github.com/synapseq-foundation/synapseq/v4/internal/info"
	"github.com/synapseq-foundation/synapseq/v4/internal/sequence"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// streamPcm(onChunk, onDone, onError, spsqUint8Array)
func streamPcm(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return js.Global().Get("Promise").New(js.FuncOf(
			func(_ js.Value, pArgs []js.Value) interface{} {
				reject := pArgs[1]
				reject.Invoke("missing callbacks")
				return nil
			},
		))
	}

	onChunk := args[0]
	onDone := args[1]
	onError := args[2]

	return js.Global().Get("Promise").New(js.FuncOf(
		func(_ js.Value, pArgs []js.Value) interface{} {

			resolve := pArgs[0]
			reject := pArgs[1]

			go func() {
				defer func() {
					if r := recover(); r != nil {
						reject.Invoke(fmt.Sprintf("panic: %v", r))
					}
				}()

				if len(args) < 5 {
					reject.Invoke("missing file data")
					return
				}

				content := args[3]

				raw := make([]byte, content.Length())
				js.CopyBytesToGo(raw, content)

				var seq *t.Sequence
				var err error

				seq, err = sequence.LoadTextSequence(raw)
				if err != nil {
					onError.Invoke(err.Error())
					reject.Invoke(err.Error())
					return
				}

				renderer, err := audio.NewAudioRenderer(seq.Periods, &audio.AudioRendererOptions{
					SampleRate: seq.Options.SampleRate,
					Volume:     seq.Options.Volume,
					Ambiance:   seq.Options.Ambiance,
					Colors:     false,
				})
				if err != nil {
					onError.Invoke(err.Error())
					reject.Invoke(err.Error())
					return
				}

				err = renderer.Render(func(samples []int) error {
					buf := make([]byte, len(samples)*2)
					for i, v := range samples {
						buf[i*2] = byte(v)
						buf[i*2+1] = byte(v >> 8)
					}

					arr := js.Global().Get("Uint8Array").New(len(buf))
					js.CopyBytesToJS(arr, buf)

					onChunk.Invoke(arr)
					return nil
				})

				if err != nil {
					onError.Invoke(err.Error())
					reject.Invoke(err.Error())
					return
				}

				onDone.Invoke()
				resolve.Invoke(true)
			}()

			return nil
		},
	))
}

func main() {
	js.Global().Set("synapseqStreamPcm", js.FuncOf(streamPcm))
	js.Global().Set("synapseqVersion", js.ValueOf(info.VERSION))
	js.Global().Set("synapseqBuildDate", js.ValueOf(info.BUILD_DATE))
	js.Global().Set("synapseqHash", js.ValueOf(info.GIT_COMMIT))

	select {}
}
