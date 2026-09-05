// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package audiosource

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gopxl/beep/v2"
	bmp3 "github.com/gopxl/beep/v2/mp3"
	bwav "github.com/gopxl/beep/v2/wav"

	p "github.com/synapseq-foundation/synapseq/v4/internal/audio/pcm"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const resampleQuality = 4

type PlaybackMode int

const (
	PlaybackLoop PlaybackMode = iota
	PlaybackFinite
)

type FileLoader func(string) ([]byte, t.AmbianceAudioFormat, error)

type Options struct {
	PlaybackMode PlaybackMode
	LoadFile     FileLoader
	SourceKind   string
	LazyLoad     bool
}

type Audio struct {
	filePaths    []string
	currentIndex int
	playbackMode PlaybackMode
	sourceKind   string
	loadFile     FileLoader
	expectedRate int

	decoder       beep.StreamSeekCloser
	sampleRate    int
	channels      int
	bitDepth      int
	closed        bool
	hasReachedEOF bool

	cachedData [][]byte
	formats    []t.AmbianceAudioFormat
	validated  []bool
	decoders   []beep.StreamSeekCloser

	decodeBuffer [][2]float64
	bufferSize   int
}

func New(filePaths []string, expectedSampleRate int, opts Options) (*Audio, error) {
	return newAudio(filePaths, expectedSampleRate, opts)
}

// NewLazy defers file reads and format validation until a source is first played.
func NewLazy(filePaths []string, expectedSampleRate int, opts Options) (*Audio, error) {
	opts.LazyLoad = true
	return newAudio(filePaths, expectedSampleRate, opts)
}

func newAudio(filePaths []string, expectedSampleRate int, opts Options) (*Audio, error) {
	sourceKind := opts.SourceKind
	if sourceKind == "" {
		sourceKind = "external audio"
	}
	if len(filePaths) == 0 {
		return &Audio{playbackMode: opts.PlaybackMode, sourceKind: sourceKind}, nil
	}
	if opts.LoadFile == nil {
		return nil, fmt.Errorf("%s file loader cannot be nil", sourceKind)
	}
	if expectedSampleRate <= 0 {
		return nil, fmt.Errorf("invalid expected sample rate: %d", expectedSampleRate)
	}

	aa := &Audio{
		filePaths:    filePaths,
		currentIndex: 0,
		playbackMode: opts.PlaybackMode,
		sourceKind:   sourceKind,
		loadFile:     opts.LoadFile,
		expectedRate: expectedSampleRate,
		sampleRate:   expectedSampleRate,
		channels:     stereoChannels,
		bufferSize:   t.BufferSize * stereoChannels,
		cachedData:   make([][]byte, len(filePaths)),
		formats:      make([]t.AmbianceAudioFormat, len(filePaths)),
		validated:    make([]bool, len(filePaths)),
		decoders:     make([]beep.StreamSeekCloser, len(filePaths)),
	}

	if !opts.LazyLoad {
		if err := aa.loadAndCacheAll(); err != nil {
			return nil, err
		}
		if err := aa.validateTracks(); err != nil {
			return nil, err
		}
		if err := aa.openFromCache(aa.currentIndex); err != nil {
			return nil, fmt.Errorf("failed to open %s file: %w", aa.sourceKind, err)
		}
	}
	aa.decodeBuffer = make([][2]float64, aa.bufferSize/stereoChannels)
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

// CloneAudio creates an independent playback cursor over the immutable cached assets.
func (aa *Audio) CloneAudio() (SampleAudio, error) {
	if aa == nil {
		return nil, fmt.Errorf("audio is nil")
	}

	return &Audio{
		filePaths:    aa.filePaths,
		playbackMode: aa.playbackMode,
		sourceKind:   aa.sourceKind,
		loadFile:     aa.loadFile,
		expectedRate: aa.expectedRate,
		sampleRate:   aa.sampleRate,
		channels:     aa.channels,
		bitDepth:     aa.bitDepth,
		cachedData:   aa.cachedData,
		formats:      aa.formats,
		validated:    aa.validated,
		decoders:     make([]beep.StreamSeekCloser, len(aa.decoders)),
		decodeBuffer: make([][2]float64, aa.bufferSize/stereoChannels),
		bufferSize:   aa.bufferSize,
	}, nil
}

// SelectAudio resets playback at index while retaining its decoder when possible.
func (aa *Audio) SelectAudio(index int) error {
	if aa == nil {
		return fmt.Errorf("audio is nil")
	}
	if index < 0 || index >= len(aa.filePaths) {
		return fmt.Errorf("invalid %s index: %d", aa.sourceKind, index)
	}
	if aa.decoders[index] == nil {
		return aa.openFromCache(index)
	}
	if err := aa.decoders[index].Seek(0); err != nil {
		return aa.restartAt(index)
	}

	aa.decoder = aa.decoders[index]
	aa.currentIndex = index
	aa.hasReachedEOF = false
	return nil
}

// RewindBufferedFrames moves playback back to the first frame retained by a
// fractional-rate buffer so fixed-rate playback can resume without a jump.
func (aa *Audio) RewindBufferedFrames(index, frames int) error {
	if aa == nil || frames <= 0 {
		return nil
	}
	if index < 0 || index >= len(aa.filePaths) {
		return fmt.Errorf("invalid %s index: %d", aa.sourceKind, index)
	}
	if aa.decoders[index] == nil {
		return nil
	}

	position := aa.decoders[index].Position() - frames
	if aa.playbackMode == PlaybackLoop {
		length := aa.decoders[index].Len()
		if length > 0 {
			position %= length
			if position < 0 {
				position += length
			}
		}
	} else if position < 0 {
		position = 0
	}
	return aa.decoders[index].Seek(position)
}

func (aa *Audio) loadAndCacheAll() error {
	for i, path := range aa.filePaths {
		if err := aa.loadAndCache(i, path); err != nil {
			return err
		}
	}
	return nil
}

func (aa *Audio) loadAndCache(index int, path string) error {
	if aa.cachedData[index] != nil {
		return nil
	}

	data, format, err := aa.loadFile(path)
	if err != nil {
		return fmt.Errorf("failed to load %s file [%d] (%s): %w", aa.sourceKind, index, path, err)
	}
	aa.cachedData[index] = data
	aa.formats[index] = format
	return nil
}

func (aa *Audio) validateTracks() error {
	if len(aa.cachedData) == 0 {
		return fmt.Errorf("no %s tracks loaded", aa.sourceKind)
	}

	for i := range aa.cachedData {
		if err := aa.validateTrack(i); err != nil {
			return err
		}
	}
	return nil
}

func (aa *Audio) validateTrack(i int) error {
	data := aa.cachedData[i]
	stream, format, err := decodeAudioData(data, aa.formats[i])
	if err != nil {
		return fmt.Errorf("failed to decode %s file [%d] (%s): %w", aa.sourceKind, i, aa.filePaths[i], err)
	}

	sr := int(format.SampleRate)
	ch := format.NumChannels

	if err := stream.Close(); err != nil {
		return fmt.Errorf("failed to close %s file [%d] (%s): %w", aa.sourceKind, i, aa.filePaths[i], err)
	}

	if ch != stereoChannels {
		return fmt.Errorf("%s track [%d] (%s) has %d channel(s), expected %d", aa.sourceKind, i, aa.filePaths[i], ch, stereoChannels)
	}

	if sr != aa.expectedRate {
		resampled, err := resampleAudioData(data, aa.formats[i], aa.expectedRate)
		if err != nil {
			return fmt.Errorf("failed to resample %s file [%d] (%s) from %d Hz to %d Hz: %w", aa.sourceKind, i, aa.filePaths[i], sr, aa.expectedRate, err)
		}

		aa.cachedData[i] = resampled
		aa.formats[i] = t.AmbianceAudioWAV

		stream, format, err = decodeAudioData(resampled, aa.formats[i])
		if err != nil {
			return fmt.Errorf("failed to decode resampled %s file [%d] (%s): %w", aa.sourceKind, i, aa.filePaths[i], err)
		}

		sr = int(format.SampleRate)
		ch = format.NumChannels

		if err := stream.Close(); err != nil {
			return fmt.Errorf("failed to close resampled %s file [%d] (%s): %w", aa.sourceKind, i, aa.filePaths[i], err)
		}

		if sr != aa.expectedRate {
			return fmt.Errorf("%s track [%d] (%s) has sample rate %d Hz after resample, expected %d Hz", aa.sourceKind, i, aa.filePaths[i], sr, aa.expectedRate)
		}

		if ch != stereoChannels {
			return fmt.Errorf("%s track [%d] (%s) has %d channel(s) after resample, expected %d", aa.sourceKind, i, aa.filePaths[i], ch, stereoChannels)
		}
	}
	aa.validated[i] = true
	return nil
}

func decodeAudioData(data []byte, audioFormat t.AmbianceAudioFormat) (beep.StreamSeekCloser, beep.Format, error) {
	reader := &bytesReadSeekCloser{Reader: bytes.NewReader(data)}
	switch audioFormat {
	case t.AmbianceAudioWAV:
		return bwav.Decode(reader)
	case t.AmbianceAudioMP3:
		return bmp3.Decode(reader)
	default:
		return nil, beep.Format{}, fmt.Errorf("unsupported external audio format %q", audioFormat.String())
	}
}

func resampleAudioData(data []byte, audioFormat t.AmbianceAudioFormat, expectedSampleRate int) ([]byte, error) {
	stream, format, err := decodeAudioData(data, audioFormat)
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

type bytesReadSeekCloser struct {
	*bytes.Reader
}

func (b *bytesReadSeekCloser) Close() error {
	return nil
}

type memoryWriteSeeker struct {
	data []byte
	pos  int
}

func (m *memoryWriteSeeker) Write(p []byte) (int, error) {
	oldLen := len(m.data)
	end := m.pos + len(p)
	if end > cap(m.data) {
		capacity := cap(m.data)
		if capacity == 0 {
			capacity = 4096
		}
		for capacity < end {
			capacity *= 2
		}
		grown := make([]byte, len(m.data), capacity)
		copy(grown, m.data)
		m.data = grown
	}
	if end > len(m.data) {
		m.data = m.data[:end]
	}
	if m.pos > oldLen {
		clear(m.data[oldLen:m.pos])
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
	return m.data
}

func (aa *Audio) openFromCache(index int) error {
	if index < 0 || index >= len(aa.cachedData) {
		return fmt.Errorf("invalid %s index: %d", aa.sourceKind, index)
	}
	if aa.cachedData[index] == nil {
		if err := aa.loadAndCache(index, aa.filePaths[index]); err != nil {
			return err
		}
	}
	if !aa.validated[index] {
		if err := aa.validateTrack(index); err != nil {
			return err
		}
	}

	if aa.decoders[index] != nil {
		aa.decoder = aa.decoders[index]
		aa.currentIndex = index
		aa.hasReachedEOF = false
		return nil
	}

	stream, format, err := decodeAudioData(aa.cachedData[index], aa.formats[index])
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
		return fmt.Errorf("invalid %s index: %d", aa.sourceKind, index)
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
		return 0, fmt.Errorf("invalid %s index: %d", aa.sourceKind, index)
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
			if aa.playbackMode == PlaybackFinite {
				for i := samplesRead; i < numSamples; i++ {
					samples[i] = 0
				}
				return numSamples, nil
			}
			if restartErr := aa.restartAt(index); restartErr != nil {
				return samplesRead, fmt.Errorf("failed to restart %s audio index %d: %w", aa.sourceKind, index, restartErr)
			}
			continue
		}
		if err != nil {
			return samplesRead, fmt.Errorf("error reading %s audio index %d: %w", aa.sourceKind, index, err)
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

	if cap(aa.decodeBuffer) < framesToRead {
		aa.decodeBuffer = make([][2]float64, framesToRead)
	}
	buf := aa.decodeBuffer[:framesToRead]
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
