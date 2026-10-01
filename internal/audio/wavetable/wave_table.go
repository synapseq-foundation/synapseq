// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package wavetable

import (
	"fmt"
	"math"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// ID is a dense waveform table index used only by the audio pipeline.
type ID int

const (
	SineID ID = iota
	SquareID
	TriangleID
	SawtoothID
)

// Registry owns compiled waveform tables and their pre-resolved IDs.
type Registry struct {
	Tables [][]int
	ids    map[t.WaveformName]ID
}

// Init initializes the waveform lookup tables used during rendering.
func Init() [][]int {
	waveTables := make([][]int, 4)
	for i := range waveTables {
		waveformTable := make([]int, t.SineTableSize)

		for j := range t.SineTableSize {
			phase := float64(j) * 2.0 * float64(math.Pi) / float64(t.SineTableSize)
			var val float64

			switch i {
			case int(SineID):
				val = math.Sin(phase)
			case int(SquareID):
				if math.Sin(phase) > 0 {
					val = 1.0
				} else {
					val = -1.0
				}
			case int(TriangleID):
				if phase < math.Pi {
					val = (2.0 * phase / math.Pi) - 1.0
				} else {
					val = 3.0 - (2.0 * phase / math.Pi)
				}
			case int(SawtoothID):
				val = 2.0 * (phase/(2.0*math.Pi) - math.Floor(phase/(2.0*math.Pi)+0.5))
			default:
				val = math.Sin(phase)
			}

			waveformTable[j] = int(t.WaveTableAmplitude * val)
		}

		waveTables[i] = waveformTable
	}

	return waveTables
}

// Compile builds a rendering registry containing built-in and custom waveforms.
func Compile(definitions []t.WaveformDefinition) (*Registry, error) {
	registry := &Registry{
		Tables: Init(),
		ids: map[t.WaveformName]ID{
			t.WaveformSine:     SineID,
			t.WaveformSquare:   SquareID,
			t.WaveformTriangle: TriangleID,
			t.WaveformSawtooth: SawtoothID,
		},
	}

	for _, definition := range definitions {
		if _, exists := registry.ids[definition.Name]; exists {
			return nil, fmt.Errorf("duplicate waveform definition: %s", definition.Name)
		}
		if len(definition.Points) < t.MinWaveformPoints || len(definition.Points) > t.MaxWaveformPoints {
			return nil, fmt.Errorf("waveform %q must contain between %d and %d points", definition.Name, t.MinWaveformPoints, t.MaxWaveformPoints)
		}
		for _, point := range definition.Points {
			if math.IsNaN(point) || math.IsInf(point, 0) || point < -1 || point > 1 {
				return nil, fmt.Errorf("waveform %q contains an invalid normalized point", definition.Name)
			}
		}

		id := ID(len(registry.Tables))
		registry.ids[definition.Name] = id
		registry.Tables = append(registry.Tables, compileCustom(definition.Points))
	}

	return registry, nil
}

// Lookup resolves a domain waveform name to a dense table ID.
func (registry *Registry) Lookup(name t.WaveformName) (ID, bool) {
	if registry == nil {
		return 0, false
	}
	id, ok := registry.ids[name.Effective()]
	return id, ok
}

// BuiltinID resolves a built-in domain name without a registry.
func BuiltinID(name t.WaveformName) (ID, bool) {
	switch name.Effective() {
	case t.WaveformSine:
		return SineID, true
	case t.WaveformSquare:
		return SquareID, true
	case t.WaveformTriangle:
		return TriangleID, true
	case t.WaveformSawtooth:
		return SawtoothID, true
	default:
		return 0, false
	}
}

func compileCustom(points []float64) []int {
	table := make([]int, t.SineTableSize)
	pointCount := len(points)
	for index := range table {
		position := float64(index) * float64(pointCount) / float64(t.SineTableSize)
		left := int(position)
		alpha := position - float64(left)
		right := (left + 1) % pointCount
		value := points[left] + (points[right]-points[left])*alpha
		table[index] = int(float64(t.WaveTableAmplitude) * value)
	}
	return table
}
