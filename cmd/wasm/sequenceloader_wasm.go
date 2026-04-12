//go:build wasm

package main

import (
	sequencepkg "github.com/synapseq-foundation/synapseq/v4/internal/sequence"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type sequenceLoader struct{}

func (sequenceLoader) Load(rawContent []byte) (*t.Sequence, error) {
	return sequencepkg.LoadTextSequence(rawContent)
}