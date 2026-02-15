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

// BackgroundAudio handles background WAV playlist playback with looping
type BackgroundAudio struct {
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

func NewBackgroundAudio(filePaths []string, expectedSampleRate int) (*BackgroundAudio, error) {
	if len(filePaths) == 0 {
		return &BackgroundAudio{isEnabled: false}, nil
	}

	bg := &BackgroundAudio{
		filePaths:    filePaths,
		currentIndex: 0,
		bufferSize:   t.BufferSize * audioChannels,
		isEnabled:    true,
		cachedData:   make([][]byte, len(filePaths)),
		decoders:     make([]beep.StreamSeekCloser, len(filePaths)),
	}

	if err := bg.loadAndCacheAll(); err != nil {
		return nil, err
	}
	if err := bg.validateTracks(expectedSampleRate); err != nil {
		return nil, err
	}
	if err := bg.openFromCache(bg.currentIndex); err != nil {
		return nil, fmt.Errorf("failed to open background file: %w", err)
	}
	bg.buffer = make([]int, bg.bufferSize)
	return bg, nil
}

// loadAndCacheAll loads all files (local or remote) into memory cache
func (bg *BackgroundAudio) loadAndCacheAll() error {
	for i, p := range bg.filePaths {
		if bg.cachedData[i] != nil {
			continue
		}

		data, err := s.GetFile(p, t.FormatWAV)
		if err != nil {
			return fmt.Errorf("failed to load background file [%d] (%s): %w", i, p, err)
		}
		bg.cachedData[i] = data
	}
	return nil
}

// validateTracks validates all tracks for compatible format and expected output format
func (bg *BackgroundAudio) validateTracks(expectedSampleRate int) error {
	if len(bg.cachedData) == 0 {
		return fmt.Errorf("no background tracks loaded")
	}
	if expectedSampleRate <= 0 {
		return fmt.Errorf("invalid expected sample rate: %d", expectedSampleRate)
	}

	// baseDepth := 0

	for i, data := range bg.cachedData {
		reader := bytes.NewReader(data)
		stream, format, err := bwav.Decode(reader)
		if err != nil {
			return fmt.Errorf("failed to decode background file [%d] (%s): %w", i, bg.filePaths[i], err)
		}
		_ = stream.Close()

		sr := int(format.SampleRate)
		ch := format.NumChannels
		// depth := format.Precision * 8

		if sr != expectedSampleRate {
			return fmt.Errorf(
				"background track [%d] (%s) has sample rate %d Hz, expected %d Hz",
				i, bg.filePaths[i], sr, expectedSampleRate,
			)
		}

		if ch != audioChannels {
			return fmt.Errorf(
				"background track [%d] (%s) has %d channel(s), expected %d",
				i, bg.filePaths[i], ch, audioChannels,
			)
		}

		// if i == 0 {
		// 	baseDepth = depth
		// } else if depth != baseDepth {
		// 	return fmt.Errorf(
		// 		"background track format mismatch in [%d] (%s): bit depth %dbit, expected %dbit",
		// 		i, bg.filePaths[i], depth, baseDepth,
		// 	)
		// }
	}

	bg.sampleRate = expectedSampleRate
	bg.channels = audioChannels
	// bg.bitDepth = baseDepth

	return nil
}

// openFromCache opens a decoder from cached data at a given index
func (bg *BackgroundAudio) openFromCache(index int) error {
	if index < 0 || index >= len(bg.cachedData) {
		return fmt.Errorf("invalid background index: %d", index)
	}
	if bg.cachedData[index] == nil {
		return fmt.Errorf("no cached data available for index: %d", index)
	}

	// If decoder already exists for this index, reuse it
	if bg.decoders[index] != nil {
		bg.decoder = bg.decoders[index]
		bg.currentIndex = index
		bg.hasReachedEOF = false
		return nil
	}

	reader := bytes.NewReader(bg.cachedData[index])
	stream, format, err := bwav.Decode(reader)
	if err != nil {
		return err
	}

	bg.decoders[index] = stream
	bg.decoder = stream
	bg.currentIndex = index
	bg.sampleRate = int(format.SampleRate)
	bg.channels = format.NumChannels
	bg.bitDepth = format.Precision * 8
	bg.hasReachedEOF = false

	return nil
}

// restart advances to next track (looping playlist) and reopens decoder
func (bg *BackgroundAudio) restartAt(index int) error {
	if index < 0 || index >= len(bg.filePaths) {
		return fmt.Errorf("invalid background index: %d", index)
	}

	if bg.decoders[index] != nil {
		_ = bg.decoders[index].Close()
		bg.decoders[index] = nil
	}

	return bg.openFromCache(index)
}

// ReadSamples reads background audio samples with automatic looping
func (bg *BackgroundAudio) ReadSamplesAt(index int, samples []int, numSamples int) (int, error) {
	if numSamples > len(samples) {
		numSamples = len(samples)
	}
	if !bg.isEnabled {
		for i := 0; i < numSamples; i++ {
			samples[i] = 0
		}
		return numSamples, nil
	}
	if index < 0 || index >= len(bg.filePaths) {
		return 0, fmt.Errorf("invalid background index: %d", index)
	}

	if bg.decoders[index] == nil {
		if err := bg.openFromCache(index); err != nil {
			return 0, err
		}
	}

	samplesRead := 0
	for samplesRead < numSamples {
		remaining := numSamples - samplesRead
		n, err := bg.readFromDecoderAt(index, samples[samplesRead:samplesRead+remaining], remaining)
		samplesRead += n

		if err == io.EOF || n < remaining {
			if restartErr := bg.restartAt(index); restartErr != nil {
				for i := samplesRead; i < numSamples; i++ {
					samples[i] = 0
				}
				return numSamples, nil
			}
			continue
		}
		if err != nil {
			return samplesRead, fmt.Errorf("error reading background audio index %d: %w", index, err)
		}
	}

	return samplesRead, nil
}

// readFromDecoder reads raw samples from the WAV decoder
func (bg *BackgroundAudio) readFromDecoderAt(index int, samples []int, maxSamples int) (int, error) {
	if index < 0 || index >= len(bg.decoders) || bg.decoders[index] == nil {
		return 0, io.EOF
	}

	decoder := bg.decoders[index]
	framesToRead := maxSamples / bg.channels
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
	outN := nFrames * bg.channels
	if outN > maxSamples {
		outN = maxSamples
	}
	framesOut := outN / bg.channels

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

// Close closes the background audio decoder
func (bg *BackgroundAudio) Close() error {
	bg.isEnabled = false
	if bg.decoder != nil {
		return bg.decoder.Close()
	}
	return nil
}

// IsEnabled returns whether background audio is enabled
func (bg *BackgroundAudio) IsEnabled() bool {
	return bg.isEnabled
}
