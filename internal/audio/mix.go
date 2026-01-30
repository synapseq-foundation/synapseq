/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

package audio

import (
	"math"

	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
)

// mix generates a stereo audio sample by mixing all channels
func (r *AudioRenderer) mix(samples []int) []int {
	// Read background audio samples if enabled
	var backgroundSamples []int

	if r.backgroundAudio.IsEnabled() {
		// Buffer for background audio
		backgroundSamples = make([]int, t.BufferSize*audioChannels) // Stereo
		r.backgroundAudio.ReadSamples(backgroundSamples, t.BufferSize*audioChannels)
	}

	for i := range t.BufferSize {
		var left, right int

		for ch := range t.NumberOfChannels {
			channel := &r.channels[ch]
			waveIdx := int(channel.Track.Waveform)

			switch channel.Track.Type {
			case t.TrackPureTone:
				channel.Offset[0] += channel.Increment[0]
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				left += channel.Amplitude[0] * r.waveTables[waveIdx][channel.Offset[0]>>16]
				right += channel.Amplitude[0] * r.waveTables[waveIdx][channel.Offset[0]>>16]
			case t.TrackBinauralBeat:
				channel.Offset[0] += channel.Increment[0]
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				channel.Offset[1] += channel.Increment[1]
				channel.Offset[1] &= (t.SineTableSize << 16) - 1

				left += channel.Amplitude[0] * r.waveTables[waveIdx][channel.Offset[0]>>16]
				right += channel.Amplitude[1] * r.waveTables[waveIdx][channel.Offset[1]>>16]
			case t.TrackMonauralBeat:
				channel.Offset[0] += channel.Increment[0]
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				channel.Offset[1] += channel.Increment[1]
				channel.Offset[1] &= (t.SineTableSize << 16) - 1

				freqHigh := r.waveTables[waveIdx][channel.Offset[0]>>16]
				freqLow := r.waveTables[waveIdx][channel.Offset[1]>>16]

				// halfAmp := channel.Amplitude[0] / 2
				// mixedSample := halfAmp * (freqHigh + freqLow)
				mixedSample := (channel.Amplitude[0] * (freqHigh + freqLow)) >> 1

				left += mixedSample
				right += mixedSample
			case t.TrackIsochronicBeat:
				channel.Offset[0] += channel.Increment[0]
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				channel.Offset[1] += channel.Increment[1]
				channel.Offset[1] &= (t.SineTableSize << 16) - 1

				modFactor := r.calcPulseFactor(channel)

				carrier := float64(r.waveTables[waveIdx][channel.Offset[0]>>16])
				amp := float64(channel.Amplitude[0])

				out := int(amp * carrier * modFactor)

				left += out
				right += out
			case t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise:
				// Use pre-generated pink noise sample for efficiency
				noiseVal := r.noiseGenerator.Generate(t.TrackPinkNoise)
				if channel.Track.Type != t.TrackPinkNoise {
					noiseVal = r.noiseGenerator.Generate(channel.Track.Type)
				}

				// Scale noise by amplitude
				sampleVal := channel.Amplitude[0] * noiseVal
				left += sampleVal
				right += sampleVal
			case t.TrackBackground:
				// Scale factor to match wavetable amplitude range
				// WaveTableAmplitude (0x7FFFF = 524287) vs 16-bit samples (32768)
				// Scale: 524287 / 32768 ≈ 16
				const bgScaleFactor = 16

				bgLeft := backgroundSamples[i*2] * bgScaleFactor
				bgRight := backgroundSamples[i*2+1] * bgScaleFactor

				// Apply gain reduction if configured (default GainLevelVeryHigh = 0dB, no reduction)
				if r.GainLevel > 0 {
					dbValue := -float64(r.GainLevel)
					gainFactor := math.Pow(10, dbValue/20.0)
					bgLeft = int(float64(bgLeft) * gainFactor)
					bgRight = int(float64(bgRight) * gainFactor)
				}

				backgroundAmplitude := channel.Amplitude[0]

				switch channel.Track.Effect.Type {
				case t.EffectSpin:
					channel.Offset[0] += channel.Increment[0]
					channel.Offset[0] &= (t.SineTableSize << 16) - 1

					spinPos := (channel.Increment[1] * r.waveTables[waveIdx][channel.Offset[0]>>16]) >> 24

					effectIntensity := float64(channel.Track.Intensity) * 0.7
					spinGain := 0.5 + effectIntensity*3.5

					ampSpin := int(float64(spinPos) * spinGain)
					if ampSpin > 127 {
						ampSpin = 127
					}
					if ampSpin < -128 {
						ampSpin = -128
					}

					posVal := ampSpin
					if posVal < 0 {
						posVal = -posVal
					}
					if posVal > 128 {
						posVal = 128
					}

					var spinLeft, spinRight int
					if ampSpin >= 0 {
						spinLeft = (bgLeft * backgroundAmplitude * (128 - posVal)) >> 7
						spinRight = bgRight*backgroundAmplitude + ((bgLeft * backgroundAmplitude * posVal) >> 7)
					} else {
						spinLeft = bgLeft*backgroundAmplitude + ((bgRight * backgroundAmplitude * posVal) >> 7)
						spinRight = (bgRight * backgroundAmplitude * (128 - posVal)) >> 7
					}

					left += spinLeft
					right += spinRight
				case t.EffectPulse:
					// LFO for pulse modulation
					channel.Offset[1] += channel.Increment[1]
					channel.Offset[1] &= (t.SineTableSize << 16) - 1

					// 0..1
					modFactor := r.calcPulseFactor(channel)

					// Mix the effect (0..1) weighted by intensity
					effectIntensity := float64(channel.Track.Intensity) * 0.7
					gain := (1.0 - effectIntensity) + (effectIntensity * modFactor)

					left += int(float64(bgLeft*backgroundAmplitude) * gain)
					right += int(float64(bgRight*backgroundAmplitude) * gain)
				default:
					// BG without effect
					left += bgLeft * backgroundAmplitude
					right += bgRight * backgroundAmplitude
				}
			}
		}

		if r.Volume != 100 {
			left = left * r.Volume / 100
			right = right * r.Volume / 100
		}

		// Scale down to 24-bit range
		left >>= audioBitShift
		right >>= audioBitShift

		// Clipping to 24-bit range
		if left > audioMaxValue {
			left = audioMaxValue
		}
		if left < audioMinValue {
			left = audioMinValue
		}
		if right > audioMaxValue {
			right = audioMaxValue
		}
		if right < audioMinValue {
			right = audioMinValue
		}

		samples[i*2] = int(left)
		samples[i*2+1] = int(right)
	}

	return samples
}
