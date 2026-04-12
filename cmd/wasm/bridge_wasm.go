//go:build wasm

package main

import (
	"fmt"
	"syscall/js"

	"github.com/synapseq-foundation/synapseq/v4/internal/info"
)

type jsCallbackSink struct {
	onChunk js.Value
	onDone  js.Value
	onError js.Value
}

type wasmStreamHandler struct {
	service *streamService
}

func newWASMStreamHandler() *wasmStreamHandler {
	return &wasmStreamHandler{service: newWASMStreamService()}
}

func newWASMStreamService() *streamService {
	return &streamService{
		loader:  sequenceLoader{},
		builder: rendererBuilder{},
		encoder: int16LEPCMEncoder{},
	}
}

func registerWASMGlobals(global js.Value) {
	handler := newWASMStreamHandler()
	global.Set("synapseqStreamPcm", js.FuncOf(handler.Handle))
	global.Set("synapseqVersion", js.ValueOf(info.VERSION))
	global.Set("synapseqBuildDate", js.ValueOf(info.BUILD_DATE))
	global.Set("synapseqHash", js.ValueOf(info.GIT_COMMIT))
}

func (h *wasmStreamHandler) Handle(_ js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return rejectedPromise("missing callbacks")
	}

	return js.Global().Get("Promise").New(js.FuncOf(func(_ js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]

		go h.handleAsync(args, resolve, reject)

		return nil
	}))
}

func (h *wasmStreamHandler) handleAsync(args []js.Value, resolve js.Value, reject js.Value) {
	sink := jsCallbackSink{
		onChunk: args[0],
		onDone:  args[1],
		onError: args[2],
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			message := fmt.Sprintf("panic: %v", recovered)
			_ = sink.OnError(message)
			reject.Invoke(message)
		}
	}()

	if len(args) < 5 {
		_ = sink.OnError("missing file data")
		reject.Invoke("missing file data")
		return
	}

	content := args[3]
	rawContent := make([]byte, content.Length())
	js.CopyBytesToGo(rawContent, content)

	if err := h.service.Stream(rawContent, sink); err != nil {
		_ = sink.OnError(err.Error())
		reject.Invoke(err.Error())
		return
	}

	resolve.Invoke(true)
}

func rejectedPromise(message string) js.Value {
	return js.Global().Get("Promise").New(js.FuncOf(func(_ js.Value, promiseArgs []js.Value) interface{} {
		promiseArgs[1].Invoke(message)
		return nil
	}))
}

func (s jsCallbackSink) OnChunk(buffer []byte) error {
	array := js.Global().Get("Uint8Array").New(len(buffer))
	js.CopyBytesToJS(array, buffer)
	s.onChunk.Invoke(array)
	return nil
}

func (s jsCallbackSink) OnDone() error {
	s.onDone.Invoke()
	return nil
}

func (s jsCallbackSink) OnError(message string) error {
	s.onError.Invoke(message)
	return nil
}