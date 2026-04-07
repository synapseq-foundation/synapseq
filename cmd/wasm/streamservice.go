package main

import types "github.com/synapseq-foundation/synapseq/v4/internal/types"

type renderableAudio interface {
	Render(func(samples []int) error) error
}

type wasmSequenceLoader interface {
	Load(rawContent []byte) (*types.Sequence, error)
}

type wasmRendererBuilder interface {
	Build(sequence *types.Sequence) (renderableAudio, error)
}

type pcmEncoder interface {
	Encode(samples []int) ([]byte, error)
}

type pcmChunkSink interface {
	OnChunk(buffer []byte) error
	OnDone() error
}

type streamService struct {
	loader  wasmSequenceLoader
	builder wasmRendererBuilder
	encoder pcmEncoder
}

func (s *streamService) Stream(rawContent []byte, sink pcmChunkSink) error {
	sequence, err := s.loader.Load(rawContent)
	if err != nil {
		return err
	}

	renderer, err := s.builder.Build(sequence)
	if err != nil {
		return err
	}

	err = renderer.Render(func(samples []int) error {
		buffer, encodeErr := s.encoder.Encode(samples)
		if encodeErr != nil {
			return encodeErr
		}

		return sink.OnChunk(buffer)
	})
	if err != nil {
		return err
	}

	return sink.OnDone()
}
