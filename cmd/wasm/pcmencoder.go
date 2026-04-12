package main

import "fmt"

const (
	pcm16LEMin = -32768
	pcm16LEMax = 32767
)

type int16LEPCMEncoder struct{}

func encodePCM16LE(samples []int) ([]byte, error) {
	buffer := make([]byte, len(samples)*2)
	for index, sample := range samples {
		if sample < pcm16LEMin || sample > pcm16LEMax {
			return nil, fmt.Errorf("sample %d out of 16-bit range: %d", index, sample)
		}

		buffer[index*2] = byte(sample)
		buffer[index*2+1] = byte(sample >> 8)
	}

	return buffer, nil
}

func (int16LEPCMEncoder) Encode(samples []int) ([]byte, error) {
	return encodePCM16LE(samples)
}