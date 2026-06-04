// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package pcm

const (
	Scale    = 32768.0
	MaxValue = 32767
	MinValue = -32768
)

func Clamp16(sample int) int {
	if sample > MaxValue {
		return MaxValue
	}
	if sample < MinValue {
		return MinValue
	}

	return sample
}

func SampleToInt16(sample int) int16 {
	return int16(Clamp16(sample))
}

func SampleToUnitFloat64(sample int) float64 {
	return float64(Clamp16(sample)) / Scale
}

func FloatToSample16(sample float64) int {
	return Clamp16(int(sample * Scale))
}

func EncodePCM16LE(dst []byte, samples []int) []byte {
	need := len(samples) * 2
	if cap(dst) < need {
		dst = make([]byte, need)
	}

	buf := dst[:need]
	writeIdx := 0
	for _, sample := range samples {
		value := SampleToInt16(sample)
		buf[writeIdx] = byte(value)
		buf[writeIdx+1] = byte(value >> 8)
		writeIdx += 2
	}

	return buf
}