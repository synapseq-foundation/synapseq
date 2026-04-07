package main

import (
	"fmt"

	"github.com/synapseq-foundation/synapseq/v4/internal/audio"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type rendererBuilder struct{}

func buildWASMRendererOptions(sequence *t.Sequence) (*audio.AudioRendererOptions, error) {
	if sequence == nil {
		return nil, fmt.Errorf("sequence is nil")
	}
	if sequence.Options == nil {
		return nil, fmt.Errorf("sequence options are nil")
	}

	return &audio.AudioRendererOptions{
		SampleRate: sequence.Options.SampleRate,
		Volume:     sequence.Options.Volume,
		Ambiance:   sequence.Options.Ambiance,
		Colors:     false,
	}, nil
}

func (rendererBuilder) Build(sequence *t.Sequence) (renderableAudio, error) {
	options, err := buildWASMRendererOptions(sequence)
	if err != nil {
		return nil, err
	}

	return audio.NewAudioRenderer(sequence.Periods, options)
}