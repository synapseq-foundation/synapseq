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
	for i := range t.BufferSize {
		var left, right int

		for ch := range t.NumberOfChannels {
			channel := &r.channels[ch]
			waveIdx := int(channel.Track.Waveform)

			switch channel.Track.Type {
			case t.TrackPureTone:
				inc0 := channel.Increment[0]

				if channel.Track.Effect.Type == t.EffectDoppler {
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					factor := r.calcDopplerFactor(channel.Effect.Offset, channel.Track.Effect.Intensity)
					inc0 = int(math.Round(float64(inc0) * factor))
				}

				channel.Offset[0] += inc0
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				sample := channel.Amplitude[0] * r.waveTables[waveIdx][channel.Offset[0]>>16]

				if channel.Track.Effect.Type == t.EffectModulation {
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					modFactor := r.calcModulationFactor(channel.Track.Waveform, channel.Effect.Offset) // 0..1
					effectIntensity := float64(channel.Track.Effect.Intensity) * 0.7                   // intensity já é 0..1

					gain := (1.0 - effectIntensity) + (effectIntensity * modFactor)
					sample = int(float64(sample) * gain)
				}

				if channel.Track.Effect.Type == t.EffectPan {
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					ll, rr := sample, sample
					ll, rr = r.applyPan(channel, ll, rr)
					left += ll
					right += rr
				} else {
					left += sample
					right += sample
				}
			case t.TrackBinauralBeat:
				inc0 := channel.Increment[0]
				inc1 := channel.Increment[1]

				if channel.Track.Effect.Type == t.EffectDoppler {
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					factor := r.calcDopplerFactor(channel.Effect.Offset, channel.Track.Effect.Intensity)
					inc0 = int(math.Round(float64(inc0) * factor))
					inc1 = int(math.Round(float64(inc1) * factor))
				}

				channel.Offset[0] += inc0
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				channel.Offset[1] += inc1
				channel.Offset[1] &= (t.SineTableSize << 16) - 1

				ll := channel.Amplitude[0] * r.waveTables[waveIdx][channel.Offset[0]>>16]
				rr := channel.Amplitude[1] * r.waveTables[waveIdx][channel.Offset[1]>>16]

				if channel.Track.Effect.Type == t.EffectModulation {
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					modFactor := r.calcModulationFactor(channel.Track.Waveform, channel.Effect.Offset) // 0..1
					effectIntensity := float64(channel.Track.Effect.Intensity) * 0.7                   // intensity já é 0..1

					gain := (1.0 - effectIntensity) + (effectIntensity * modFactor)
					ll = int(float64(ll) * gain)
					rr = int(float64(rr) * gain)
				}

				if channel.Track.Effect.Type == t.EffectPan {
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					ll, rr = r.applyPan(channel, ll, rr)
				}

				left += ll
				right += rr
			case t.TrackMonauralBeat:
				inc0 := channel.Increment[0]
				inc1 := channel.Increment[1]

				if channel.Track.Effect.Type == t.EffectDoppler {
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					factor := r.calcDopplerFactor(channel.Effect.Offset, channel.Track.Effect.Intensity)
					inc0 = int(math.Round(float64(inc0) * factor))
					inc1 = int(math.Round(float64(inc1) * factor))
				}

				channel.Offset[0] += inc0
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				channel.Offset[1] += inc1
				channel.Offset[1] &= (t.SineTableSize << 16) - 1

				freqHigh := r.waveTables[waveIdx][channel.Offset[0]>>16]
				freqLow := r.waveTables[waveIdx][channel.Offset[1]>>16]

				mixedSample := (channel.Amplitude[0] * (freqHigh + freqLow)) >> 1

				if channel.Track.Effect.Type == t.EffectModulation {
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					modFactor := r.calcModulationFactor(channel.Track.Waveform, channel.Effect.Offset) // 0..1
					effectIntensity := float64(channel.Track.Effect.Intensity) * 0.7
					gain := (1.0 - effectIntensity) + (effectIntensity * modFactor)

					mixedSample = int(float64(mixedSample) * gain)
				}

				if channel.Track.Effect.Type == t.EffectPan {
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					ll, rr := mixedSample, mixedSample
					ll, rr = r.applyPan(channel, ll, rr)

					left += ll
					right += rr
				} else {
					left += mixedSample
					right += mixedSample
				}
			case t.TrackIsochronicBeat:
				incCarrier := channel.Increment[0]

				if channel.Track.Effect.Type == t.EffectDoppler {
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					factor := r.calcDopplerFactor(channel.Effect.Offset, channel.Track.Effect.Intensity)
					incCarrier = int(math.Round(float64(incCarrier) * factor))
				}

				channel.Offset[0] += incCarrier
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				channel.Offset[1] += channel.Increment[1]
				channel.Offset[1] &= (t.SineTableSize << 16) - 1

				modFactor := r.calcModulationFactor(channel.Track.Waveform, channel.Offset[1])

				carrier := float64(r.waveTables[waveIdx][channel.Offset[0]>>16])
				amp := float64(channel.Amplitude[0])

				out := int(amp * carrier * modFactor)

				if channel.Track.Effect.Type == t.EffectModulation {
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					modFactor := r.calcModulationFactor(channel.Track.Waveform, channel.Effect.Offset) // 0..1
					effectIntensity := float64(channel.Track.Effect.Intensity) * 0.7
					gain := (1.0 - effectIntensity) + (effectIntensity * modFactor)

					out = int(float64(out) * gain)
				}

				if channel.Track.Effect.Type == t.EffectPan {
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					ll, rr := out, out
					ll, rr = r.applyPan(channel, ll, rr)

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
				case t.EffectModulation:
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					modFactor := r.calcModulationFactor(channel.Track.Waveform, channel.Effect.Offset) // 0..1

					// Intensity (0..100) -> 0..1
					effectIntensity := float64(channel.Track.Effect.Intensity) * 0.7

					gain := (1.0 - effectIntensity) + (effectIntensity * modFactor)
					sampleVal = int(float64(sampleVal) * gain)

					left += sampleVal
					right += sampleVal
				case t.EffectPan:
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					ll, rr := sampleVal, sampleVal
					ll, rr = r.applyPan(channel, ll, rr)

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

				idx := channel.Track.BackgroundIndex
				if idx < 0 || idx >= len(r.backgroundSamplesByIndex) {
					continue
				}
				bgBuf := r.backgroundSamplesByIndex[idx]
				if len(bgBuf) < i*2+2 {
					continue
				}

				bgLeft := bgBuf[i*2] * bgScaleFactor
				bgRight := bgBuf[i*2+1] * bgScaleFactor

				// Apply gain reduction if configured (default GainLevelVeryHigh = 0dB, no reduction)
				if r.GainLevel > 0 {
					dbValue := -float64(r.GainLevel)
					gainFactor := math.Pow(10, dbValue/20.0)
					bgLeft = int(float64(bgLeft) * gainFactor)
					bgRight = int(float64(bgRight) * gainFactor)
				}

				backgroundAmplitude := channel.Amplitude[0]

				switch channel.Track.Effect.Type {
				case t.EffectPan:
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					inL := bgLeft * backgroundAmplitude
					inR := bgRight * backgroundAmplitude

					outL, outR := r.applyPan(channel, inL, inR)
					left += outL
					right += outR
				case t.EffectModulation:
					// LFO for modulation
					channel.Effect.Offset += channel.Effect.Increment
					channel.Effect.Offset &= (t.SineTableSize << 16) - 1

					// 0..1
					modFactor := r.calcModulationFactor(channel.Track.Waveform, channel.Effect.Offset)

					// Mix the effect (0..1) weighted by intensity
					effectIntensity := float64(channel.Track.Effect.Intensity) * 0.7
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
