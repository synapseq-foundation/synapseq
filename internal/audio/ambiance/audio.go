/*
 * SynapSeq - Text-Driven Audio Sequencer for Brainwave Entrainment
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

package ambiance

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gopxl/beep/v2"
	bwav "github.com/gopxl/beep/v2/wav"

	p "github.com/synapseq-foundation/synapseq/v4/internal/audio/pcm"
	r "github.com/synapseq-foundation/synapseq/v4/internal/resource"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const resampleQuality = 4

type Audio struct {
	filePaths    []string
	currentIndex int

	decoder       beep.StreamSeekCloser
	sampleRate    int
	channels      int
	bitDepth      int
	closed        bool
	hasReachedEOF bool

	cachedData [][]byte
	decoders   []beep.StreamSeekCloser

	buffer     []int
	bufferSize int
}

func NewAudio(filePaths []string, expectedSampleRate int) (*Audio, error) {
	if len(filePaths) == 0 {
		return &Audio{}, nil
	}

	aa := &Audio{
		filePaths:    filePaths,
		currentIndex: 0,
		bufferSize:   t.BufferSize * stereoChannels,
		cachedData:   make([][]byte, len(filePaths)),
		decoders:     make([]beep.StreamSeekCloser, len(filePaths)),
	}

	if err := aa.loadAndCacheAll(); err != nil {
		return nil, err
	}
	if err := aa.validateTracks(expectedSampleRate); err != nil {
		return nil, err
	}
	if err := aa.openFromCache(aa.currentIndex); err != nil {
		return nil, fmt.Errorf("failed to open ambiance file: %w", err)
	}
	aa.buffer = make([]int, aa.bufferSize)
	return aa, nil
}

func (aa *Audio) SampleRate() int {
	return aa.sampleRate
}

func (aa *Audio) Channels() int {
	return aa.channels
}

func (aa *Audio) BitDepth() int {
	return aa.bitDepth
}

func (aa *Audio) BufferSize() int {
	return aa.bufferSize
}

func (aa *Audio) CachedData() [][]byte {
	return aa.cachedData
}

func (aa *Audio) loadAndCacheAll() error {
	for i, path := range aa.filePaths {
		if aa.cachedData[i] != nil {
			continue
		}

		data, err := r.GetFile(path, t.FormatWAV)
		if err != nil {
			return fmt.Errorf("failed to load ambiance file [%d] (%s): %w", i, path, err)
		}
		aa.cachedData[i] = data
	}
	return nil
}

func (aa *Audio) validateTracks(expectedSampleRate int) error {
	if len(aa.cachedData) == 0 {
		return fmt.Errorf("no ambiance tracks loaded")
	}
	if expectedSampleRate <= 0 {
		return fmt.Errorf("invalid expected sample rate: %d", expectedSampleRate)
	}

	for i, data := range aa.cachedData {
		reader := bytes.NewReader(data)
		stream, format, err := bwav.Decode(reader)
		if err != nil {
			return fmt.Errorf("failed to decode ambiance file [%d] (%s): %w", i, aa.filePaths[i], err)
		}

		sr := int(format.SampleRate)
		ch := format.NumChannels

		if err := stream.Close(); err != nil {
			return fmt.Errorf("failed to close ambiance file [%d] (%s): %w", i, aa.filePaths[i], err)
		}

		if ch != stereoChannels {
			return fmt.Errorf("ambiance track [%d] (%s) has %d channel(s), expected %d", i, aa.filePaths[i], ch, stereoChannels)
		}

		if sr != expectedSampleRate {
			resampled, err := resampleWAVData(data, expectedSampleRate)
			if err != nil {
				return fmt.Errorf("failed to resample ambiance file [%d] (%s) from %d Hz to %d Hz: %w", i, aa.filePaths[i], sr, expectedSampleRate, err)
			}

			aa.cachedData[i] = resampled

			reader = bytes.NewReader(resampled)
			stream, format, err = bwav.Decode(reader)
			if err != nil {
				return fmt.Errorf("failed to decode resampled ambiance file [%d] (%s): %w", i, aa.filePaths[i], err)
			}

			sr = int(format.SampleRate)
			ch = format.NumChannels

			if err := stream.Close(); err != nil {
				return fmt.Errorf("failed to close resampled ambiance file [%d] (%s): %w", i, aa.filePaths[i], err)
			}

			if sr != expectedSampleRate {
				return fmt.Errorf("ambiance track [%d] (%s) has sample rate %d Hz after resample, expected %d Hz", i, aa.filePaths[i], sr, expectedSampleRate)
			}

			if ch != stereoChannels {
				return fmt.Errorf("ambiance track [%d] (%s) has %d channel(s) after resample, expected %d", i, aa.filePaths[i], ch, stereoChannels)
			}
		}
	}

	aa.sampleRate = expectedSampleRate
	aa.channels = stereoChannels
	return nil
}

func resampleWAVData(data []byte, expectedSampleRate int) ([]byte, error) {
	reader := bytes.NewReader(data)
	stream, format, err := bwav.Decode(reader)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	resampled := beep.Resample(
		resampleQuality,
		format.SampleRate,
		beep.SampleRate(expectedSampleRate),
		stream,
	)

	out := &memoryWriteSeeker{}
	outFormat := beep.Format{
		SampleRate:  beep.SampleRate(expectedSampleRate),
		NumChannels: format.NumChannels,
		Precision:   format.Precision,
	}
	if err := bwav.Encode(out, resampled, outFormat); err != nil {
		return nil, err
	}
	if err := resampled.Err(); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

type memoryWriteSeeker struct {
	data []byte
	pos  int
}

func (m *memoryWriteSeeker) Write(p []byte) (int, error) {
	end := m.pos + len(p)
	if end > len(m.data) {
		grown := make([]byte, end)
		copy(grown, m.data)
		m.data = grown
	}
	copy(m.data[m.pos:end], p)
	m.pos = end
	return len(p), nil
}

func (m *memoryWriteSeeker) Seek(offset int64, whence int) (int64, error) {
	var next int
	switch whence {
	case io.SeekStart:
		next = int(offset)
	case io.SeekCurrent:
		next = m.pos + int(offset)
	case io.SeekEnd:
		next = len(m.data) + int(offset)
	default:
		return 0, fmt.Errorf("invalid seek whence: %d", whence)
	}
	if next < 0 {
		return 0, fmt.Errorf("invalid seek position: %d", next)
	}
	m.pos = next
	return int64(next), nil
}

func (m *memoryWriteSeeker) Bytes() []byte {
	return append([]byte(nil), m.data...)
}

func (aa *Audio) openFromCache(index int) error {
	if index < 0 || index >= len(aa.cachedData) {
		return fmt.Errorf("invalid ambiance index: %d", index)
	}
	if aa.cachedData[index] == nil {
		return fmt.Errorf("no cached data available for index: %d", index)
	}

	if aa.decoders[index] != nil {
		aa.decoder = aa.decoders[index]
		aa.currentIndex = index
		aa.hasReachedEOF = false
		return nil
	}

	reader := bytes.NewReader(aa.cachedData[index])
	stream, format, err := bwav.Decode(reader)
	if err != nil {
		return err
	}

	aa.decoders[index] = stream
	aa.decoder = stream
	aa.currentIndex = index
	aa.sampleRate = int(format.SampleRate)
	aa.channels = format.NumChannels
	aa.bitDepth = format.Precision * 8
	aa.hasReachedEOF = false

	return nil
}

func (aa *Audio) restartAt(index int) error {
	if index < 0 || index >= len(aa.filePaths) {
		return fmt.Errorf("invalid ambiance index: %d", index)
	}

	if aa.decoders[index] != nil {
		_ = aa.decoders[index].Close()
		aa.decoders[index] = nil
	}

	return aa.openFromCache(index)
}

func (aa *Audio) ReadSamplesAt(index int, samples []int, numSamples int) (int, error) {
	if numSamples > len(samples) {
		numSamples = len(samples)
	}
	if aa.closed || len(aa.filePaths) == 0 {
		for i := 0; i < numSamples; i++ {
			samples[i] = 0
		}
		return numSamples, nil
	}
	if index < 0 || index >= len(aa.filePaths) {
		return 0, fmt.Errorf("invalid ambiance index: %d", index)
	}

	if aa.decoders[index] == nil {
		if err := aa.openFromCache(index); err != nil {
			return 0, err
		}
	}

	samplesRead := 0
	for samplesRead < numSamples {
		remaining := numSamples - samplesRead
		n, err := aa.readFromDecoderAt(index, samples[samplesRead:samplesRead+remaining], remaining)
		samplesRead += n

		if err == io.EOF || n < remaining {
			if restartErr := aa.restartAt(index); restartErr != nil {
				for i := samplesRead; i < numSamples; i++ {
					samples[i] = 0
				}
				return numSamples, nil
			}
			continue
		}
		if err != nil {
			return samplesRead, fmt.Errorf("error reading ambiance audio index %d: %w", index, err)
		}
	}

	return samplesRead, nil
}

func (aa *Audio) readFromDecoderAt(index int, samples []int, maxSamples int) (int, error) {
	if index < 0 || index >= len(aa.decoders) || aa.decoders[index] == nil {
		return 0, io.EOF
	}

	decoder := aa.decoders[index]
	framesToRead := maxSamples / aa.channels
	if framesToRead <= 0 {
		framesToRead = 1
	}

	buf := make([][2]float64, framesToRead)
	nFrames, ok := decoder.Stream(buf)
	if !ok || nFrames == 0 {
		if err := decoder.Err(); err != nil {
			return 0, err
		}
		return 0, io.EOF
	}

	outN := nFrames * aa.channels
	if outN > maxSamples {
		outN = maxSamples
	}
	framesOut := outN / aa.channels

	for i := 0; i < framesOut; i++ {
		l := p.FloatToSample16(buf[i][0])
		r := p.FloatToSample16(buf[i][1])

		samples[2*i] = l
		if 2*i+1 < outN {
			samples[2*i+1] = r
		}
	}

	return outN, nil
}

func (aa *Audio) Close() error {
	aa.closed = true
	var firstErr error

	for i := range aa.decoders {
		if aa.decoders[i] != nil {
			if err := aa.decoders[i].Close(); err != nil && firstErr == nil {
				firstErr = err
			}
			aa.decoders[i] = nil
		}
	}

	aa.decoder = nil
	return firstErr
}
