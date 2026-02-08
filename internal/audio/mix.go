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

	if r.backgroundAudio != nil && r.backgroundAudio.IsEnabled() {
		need := t.BufferSize * audioChannels // Stereo
		if len(r.backgroundSamples) != need {
			r.backgroundSamples = make([]int, need)
		}

		backgroundSamples = r.backgroundSamples
		r.backgroundAudio.ReadSamples(backgroundSamples, need)
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

				ll := channel.Amplitude[0] * r.waveTables[waveIdx][channel.Offset[0]>>16]
				rr := channel.Amplitude[1] * r.waveTables[waveIdx][channel.Offset[1]>>16]

				if channel.Track.Effect.Type == t.EffectPulse {
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					modFactor := r.calcPulseFactor(channel.Track.Waveform, channel.Effect.Offset) // 0..1
					effectIntensity := float64(channel.Track.Intensity) * 0.7

					gain := (1.0 - effectIntensity) + (effectIntensity * modFactor)
					ll = int(float64(ll) * gain)
					rr = int(float64(rr) * gain)
				}

				left += ll
				right += rr
			case t.TrackMonauralBeat:
				channel.Offset[0] += channel.Increment[0]
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				channel.Offset[1] += channel.Increment[1]
				channel.Offset[1] &= (t.SineTableSize << 16) - 1

				freqHigh := r.waveTables[waveIdx][channel.Offset[0]>>16]
				freqLow := r.waveTables[waveIdx][channel.Offset[1]>>16]

				mixedSample := (channel.Amplitude[0] * (freqHigh + freqLow)) >> 1

				if channel.Track.Effect.Type == t.EffectSpin {
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					ll := mixedSample
					rr := mixedSample

					ll, rr = r.applySpin(channel, ll, rr)

					left += ll
					right += rr
				} else {
					left += mixedSample
					right += mixedSample
				}
			case t.TrackIsochronicBeat:
				channel.Offset[0] += channel.Increment[0]
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				channel.Offset[1] += channel.Increment[1]
				channel.Offset[1] &= (t.SineTableSize << 16) - 1

				modFactor := r.calcPulseFactor(channel.Track.Waveform, channel.Offset[1])

				carrier := float64(r.waveTables[waveIdx][channel.Offset[0]>>16])
				amp := float64(channel.Amplitude[0])

				out := int(amp * carrier * modFactor)

				// Optional Spin effect (pan/crossmix)
				if channel.Track.Effect.Type == t.EffectSpin {
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					ll := out
					rr := out

					ll, rr = r.applySpin(channel, ll, rr)

					left += ll
					right += rr
				} else {
					left += out
					right += out
				}
			case t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise:
				// Use pre-generated pink noise sample for efficiency
				noiseVal := r.noiseGenerator.Generate(t.TrackPinkNoise)
				if channel.Track.Type != t.TrackPinkNoise {
					noiseVal = r.noiseGenerator.Generate(channel.Track.Type)
				}

				// Scale noise by amplitude
				sampleVal := channel.Amplitude[0] * noiseVal

				switch channel.Track.Effect.Type {
				case t.EffectPulse:
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					modFactor := r.calcPulseFactor(channel.Track.Waveform, channel.Effect.Offset) // 0..1

					// Intensity (0..100) -> 0..1
					effectIntensity := float64(channel.Track.Intensity) * 0.7

					gain := (1.0 - effectIntensity) + (effectIntensity * modFactor)
					sampleVal = int(float64(sampleVal) * gain)

					left += sampleVal
					right += sampleVal
				case t.EffectSpin:
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					ll, rr := sampleVal, sampleVal
					ll, rr = r.applySpin(channel, ll, rr)

					left += ll
					right += rr

				default:
					left += sampleVal
					right += sampleVal
				}
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
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					inL := bgLeft * backgroundAmplitude
					inR := bgRight * backgroundAmplitude

					outL, outR := r.applySpin(channel, inL, inR)
					left += outL
					right += outR
				case t.EffectPulse:
					// LFO for pulse modulation
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					// 0..1
					modFactor := r.calcPulseFactor(channel.Track.Waveform, channel.Effect.Offset)

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
