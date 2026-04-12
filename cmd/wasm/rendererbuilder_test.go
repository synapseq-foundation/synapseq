package main

import (
	"testing"

	types "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestBuildWASMRendererOptions(t *testing.T) {
	sequence := &types.Sequence{
		Options: &types.SequenceOptions{
			SampleRate: 48000,
			Volume:     75,
			Ambiance: map[string]string{
				"rain": "rain.wav",
			},
		},
	}

	options, err := buildWASMRendererOptions(sequence)
	if err != nil {
		t.Fatalf("buildWASMRendererOptions returned error: %v", err)
	}
	if options.SampleRate != 48000 || options.Volume != 75 {
		t.Fatalf("unexpected renderer options: %#v", options)
	}
	if options.Ambiance["rain"] != "rain.wav" {
		t.Fatalf("expected ambiance to be preserved, got %#v", options.Ambiance)
	}
	if options.Colors {
		t.Fatal("expected colors to be disabled in wasm renderer options")
	}
}

func TestBuildWASMRendererOptionsErrors(t *testing.T) {
	if _, err := buildWASMRendererOptions(nil); err == nil {
		t.Fatal("expected nil sequence error")
	}
	if _, err := buildWASMRendererOptions(&types.Sequence{}); err == nil {
		t.Fatal("expected nil sequence options error")
	}
}