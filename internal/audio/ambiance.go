package audio

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gopxl/beep/v2"
	bwav "github.com/gopxl/beep/v2/wav"

	s "github.com/synapseq-foundation/synapseq/v3/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
)

// AmbianceAudio handles ambiance WAV playlist playback with looping
type AmbianceAudio struct {
	filePaths    []string
	currentIndex int

	decoder       beep.StreamSeekCloser
	sampleRate    int
	channels      int
	bitDepth      int
	isEnabled     bool
	hasReachedEOF bool

	cachedData [][]byte
	decoders   []beep.StreamSeekCloser

	buffer     []int
	bufferSize int
}

// NewAmbianceAudio creates a new AmbianceAudio instance with the given file paths and expected sample rate,
// loading and validating the ambiance audio tracks
func NewAmbianceAudio(filePaths []string, expectedSampleRate int) (*AmbianceAudio, error) {
	if len(filePaths) == 0 {
		return &AmbianceAudio{isEnabled: false}, nil
	}

	aa := &AmbianceAudio{
		filePaths:    filePaths,
		currentIndex: 0,
		bufferSize:   t.BufferSize * audioChannels,
		isEnabled:    true,
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

// loadAndCacheAll loads all files (local or remote) into memory cache
func (aa *AmbianceAudio) loadAndCacheAll() error {
	for i, p := range aa.filePaths {
		if aa.cachedData[i] != nil {
			continue
		}

		data, err := s.GetFile(p, t.FormatWAV)
		if err != nil {
			return fmt.Errorf("failed to load ambiance file [%d] (%s): %w", i, p, err)
		}
		aa.cachedData[i] = data
	}
	return nil
}

// validateTracks validates all tracks for compatible format and expected output format
func (aa *AmbianceAudio) validateTracks(expectedSampleRate int) error {
	if len(aa.cachedData) == 0 {
		return fmt.Errorf("no ambiance tracks loaded")
	}
	if expectedSampleRate <= 0 {
		return fmt.Errorf("invalid expected sample rate: %d", expectedSampleRate)
	}

	// baseDepth := 0

	for i, data := range aa.cachedData {
		reader := bytes.NewReader(data)
		stream, format, err := bwav.Decode(reader)
		if err != nil {
			return fmt.Errorf("failed to decode ambiance file [%d] (%s): %w", i, aa.filePaths[i], err)
		}
		_ = stream.Close()

		sr := int(format.SampleRate)
		ch := format.NumChannels
		// depth := format.Precision * 8

		if sr != expectedSampleRate {
			return fmt.Errorf(
				"ambiance track [%d] (%s) has sample rate %d Hz, expected %d Hz",
				i, aa.filePaths[i], sr, expectedSampleRate,
			)
		}

		if ch != audioChannels {
			return fmt.Errorf(
				"ambiance track [%d] (%s) has %d channel(s), expected %d",
				i, aa.filePaths[i], ch, audioChannels,
			)
		}

		// if i == 0 {
		// 	baseDepth = depth
		// } else if depth != baseDepth {
		// 	return fmt.Errorf(
		// 		"ambiance track format mismatch in [%d] (%s): bit depth %dbit, expected %dbit",
		// 		i, aa.filePaths[i], depth, baseDepth,
		// 	)
		// }
	}

	aa.sampleRate = expectedSampleRate
	aa.channels = audioChannels
	// aa.bitDepth = baseDepth

	return nil
}

// openFromCache opens a decoder from cached data at a given index
func (aa *AmbianceAudio) openFromCache(index int) error {
	if index < 0 || index >= len(aa.cachedData) {
		return fmt.Errorf("invalid ambiance index: %d", index)
	}
	if aa.cachedData[index] == nil {
		return fmt.Errorf("no cached data available for index: %d", index)
	}

	// If decoder already exists for this index, reuse it
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

// restart advances to next track (looping playlist) and reopens decoder
func (aa *AmbianceAudio) restartAt(index int) error {
	if index < 0 || index >= len(aa.filePaths) {
		return fmt.Errorf("invalid ambiance index: %d", index)
	}

	if aa.decoders[index] != nil {
		_ = aa.decoders[index].Close()
		aa.decoders[index] = nil
	}

	return aa.openFromCache(index)
}

// ReadSamples reads ambiance audio samples with automatic looping
func (aa *AmbianceAudio) ReadSamplesAt(index int, samples []int, numSamples int) (int, error) {
	if numSamples > len(samples) {
		numSamples = len(samples)
	}
	if !aa.isEnabled {
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

// readFromDecoder reads raw samples from the WAV decoder
func (aa *AmbianceAudio) readFromDecoderAt(index int, samples []int, maxSamples int) (int, error) {
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

	const scale = 32768.0
	outN := nFrames * aa.channels
	if outN > maxSamples {
		outN = maxSamples
	}
	framesOut := outN / aa.channels

	for i := 0; i < framesOut; i++ {
		l := int(buf[i][0] * scale)
		r := int(buf[i][1] * scale)

		if l > audioMaxValue {
			l = audioMaxValue
		}
		if l < audioMinValue {
			l = audioMinValue
		}
		if r > audioMaxValue {
			r = audioMaxValue
		}
		if r < audioMinValue {
			r = audioMinValue
		}

		samples[2*i] = l
		if 2*i+1 < outN {
			samples[2*i+1] = r
		}
	}

	return outN, nil
}

// Close closes the ambiance audio decoder
func (aa *AmbianceAudio) Close() error {
	aa.isEnabled = false
	if aa.decoder != nil {
		return aa.decoder.Close()
	}
	return nil
}

// IsEnabled returns whether ambiance audio is enabled
func (aa *AmbianceAudio) IsEnabled() bool {
	return aa.isEnabled
}
