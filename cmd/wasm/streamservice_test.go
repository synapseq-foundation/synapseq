package main

import (
	"errors"
	"strings"
	"testing"

	types "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type fakeLoader struct {
	sequence *types.Sequence
	err      error
}

func (f fakeLoader) Load(_ []byte) (*types.Sequence, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.sequence, nil
}

type fakeBuilder struct {
	renderer renderableAudio
	err      error
}

func (f fakeBuilder) Build(_ *types.Sequence) (renderableAudio, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.renderer, nil
}

type fakeRenderer struct {
	samples [][]int
	err     error
}

func (f fakeRenderer) Render(consume func(samples []int) error) error {
	if f.err != nil {
		return f.err
	}
	for _, sampleChunk := range f.samples {
		if err := consume(sampleChunk); err != nil {
			return err
		}
	}
	return nil
}

type fakeEncoder struct {
	buffer []byte
	err    error
}

func (f fakeEncoder) Encode(_ []int) ([]byte, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.buffer, nil
}

type fakeSink struct {
	chunks [][]byte
	done   bool
	err    error
}

func (f *fakeSink) OnChunk(buffer []byte) error {
	copyBuffer := make([]byte, len(buffer))
	copy(copyBuffer, buffer)
	f.chunks = append(f.chunks, copyBuffer)
	return f.err
}

func (f *fakeSink) OnDone() error {
	f.done = true
	return nil
}

func TestStreamServiceStreamSuccess(t *testing.T) {
	service := &streamService{
		loader:  fakeLoader{sequence: &types.Sequence{Options: &types.SequenceOptions{}}},
		builder: fakeBuilder{renderer: fakeRenderer{samples: [][]int{{1, 2, 3}}}},
		encoder: fakeEncoder{buffer: []byte{1, 2, 3, 4}},
	}
	sink := &fakeSink{}

	if err := service.Stream([]byte("sequence"), sink); err != nil {
		t.Fatalf("Stream returned error: %v", err)
	}
	if !sink.done {
		t.Fatal("expected sink OnDone to be called")
	}
	if len(sink.chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(sink.chunks))
	}
}

func TestStreamServiceStreamPropagatesLoaderError(t *testing.T) {
	service := &streamService{
		loader:  fakeLoader{err: errors.New("load failed")},
		builder: fakeBuilder{},
		encoder: fakeEncoder{},
	}
	sink := &fakeSink{}

	err := service.Stream([]byte("sequence"), sink)
	if err == nil || !strings.Contains(err.Error(), "load failed") {
		t.Fatalf("expected loader error, got %v", err)
	}
}

func TestStreamServiceStreamPropagatesRenderError(t *testing.T) {
	service := &streamService{
		loader:  fakeLoader{sequence: &types.Sequence{Options: &types.SequenceOptions{}}},
		builder: fakeBuilder{renderer: fakeRenderer{err: errors.New("render failed")}},
		encoder: fakeEncoder{},
	}
	sink := &fakeSink{}

	err := service.Stream([]byte("sequence"), sink)
	if err == nil || !strings.Contains(err.Error(), "render failed") {
		t.Fatalf("expected render error, got %v", err)
	}
	if sink.done {
		t.Fatal("did not expect sink OnDone to be called on render failure")
	}
}

func TestStreamServiceStreamPropagatesSinkChunkError(t *testing.T) {
	service := &streamService{
		loader:  fakeLoader{sequence: &types.Sequence{Options: &types.SequenceOptions{}}},
		builder: fakeBuilder{renderer: fakeRenderer{samples: [][]int{{1, 2, 3}}}},
		encoder: fakeEncoder{buffer: []byte{1, 2, 3, 4}},
	}
	sink := &fakeSink{err: errors.New("chunk failed")}

	err := service.Stream([]byte("sequence"), sink)
	if err == nil || !strings.Contains(err.Error(), "chunk failed") {
		t.Fatalf("expected chunk error, got %v", err)
	}
	if sink.done {
		t.Fatal("did not expect sink OnDone to be called when chunk delivery fails")
	}
}