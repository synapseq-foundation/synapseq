// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package spsq

import t "github.com/synapseq-foundation/synapseq/v4/internal/types"

// Preset builds the tracks for one named .spsq preset.
type Preset struct {
	builder *Builder
	name    string
	index   int
}

// NewPreset creates a new preset with the given name.
func (b *Builder) NewPreset(name string) *Preset {
	if index := b.presetIndex(name); index >= 0 {
		b.presets[index].tracks = []t.Track{}
		return &Preset{builder: b, name: name, index: index}
	}

	b.presets = append(b.presets, presetEntry{name: name, tracks: []t.Track{}})
	return &Preset{builder: b, name: name, index: len(b.presets) - 1}
}

// Tone adds a pure tone track.
func (p *Preset) Tone(carrier float64) *Preset {
	return p.addTrack(t.Track{
		Type:    t.TrackPureTone,
		Carrier: carrier,
	})
}

// Ambiance adds an ambiance track.
func (p *Preset) Ambiance(name string) *Preset {
	return p.addTrack(t.Track{
		Type:       t.TrackAmbiance,
		SourceName: name,
		Waveform:   t.WaveformSine,
	})
}

// Music adds a music track.
func (p *Preset) Music(name string) *Preset {
	return p.addTrack(t.Track{
		Type:       t.TrackMusic,
		SourceName: name,
		Waveform:   t.WaveformSine,
	})
}

// WhiteNoise adds a white noise track.
func (p *Preset) WhiteNoise(smooth float64) *Preset {
	return p.addTrack(t.Track{
		Type:        t.TrackWhiteNoise,
		NoiseSmooth: smooth,
	})
}

// PinkNoise adds a pink noise track.
func (p *Preset) PinkNoise(smooth float64) *Preset {
	return p.addTrack(t.Track{
		Type:        t.TrackPinkNoise,
		NoiseSmooth: smooth,
	})
}

// BrownNoise adds a brown noise track.
func (p *Preset) BrownNoise(smooth float64) *Preset {
	return p.addTrack(t.Track{
		Type:        t.TrackBrownNoise,
		NoiseSmooth: smooth,
	})
}

// Sine sets the waveform of the last track to sine.
func (p *Preset) Sine() *Preset {
	return p.setWaveform(t.WaveformSine)
}

// Square sets the waveform of the last track to square.
func (p *Preset) Square() *Preset {
	return p.setWaveform(t.WaveformSquare)
}

// Triangle sets the waveform of the last track to triangle.
func (p *Preset) Triangle() *Preset {
	return p.setWaveform(t.WaveformTriangle)
}

// Sawtooth sets the waveform of the last track to sawtooth.
func (p *Preset) Sawtooth() *Preset {
	return p.setWaveform(t.WaveformSawtooth)
}

// Binaural converts the last tone track to a binaural beat.
func (p *Preset) Binaural(beat float64) *Preset {
	track := p.lastTrack()
	if track == nil {
		return p
	}

	track.Type = t.TrackBinauralBeat
	track.Resonance = beat
	return p
}

// Monaural converts the last tone track to a monaural beat.
func (p *Preset) Monaural(beat float64) *Preset {
	track := p.lastTrack()
	if track == nil {
		return p
	}

	track.Type = t.TrackMonauralBeat
	track.Resonance = beat
	return p
}

// Isochronic converts the last tone track to an isochronic beat.
func (p *Preset) Isochronic(beat float64) *Preset {
	track := p.lastTrack()
	if track == nil {
		return p
	}

	track.Type = t.TrackIsochronicBeat
	track.Resonance = beat
	return p
}

// Pan adds a pan effect to the last track.
func (p *Preset) Pan(value float64) *Preset {
	return p.setEffect(t.EffectPan, value)
}

// Modulation adds a modulation effect to the last track.
func (p *Preset) Modulation(value float64) *Preset {
	return p.setEffect(t.EffectModulation, value)
}

// Doppler adds a doppler effect to the last track.
func (p *Preset) Doppler(value float64) *Preset {
	return p.setEffect(t.EffectDoppler, value)
}

// Intensity sets the effect intensity on the last track.
func (p *Preset) Intensity(value float64) *Preset {
	track := p.lastTrack()
	if track == nil {
		return p
	}

	track.Effect.Intensity = t.IntensityPercentToRaw(value)
	return p
}

// Amplitude sets the amplitude of the last track.
func (p *Preset) Amplitude(amplitude float64) *Preset {
	track := p.lastTrack()
	if track == nil {
		return p
	}

	track.Amplitude = t.AmplitudePercentToRaw(amplitude)
	return p
}

// addTrack adds a track to the current preset.
func (p *Preset) addTrack(track t.Track) *Preset {
	if p == nil || p.builder == nil || !p.valid() {
		return p
	}

	p.builder.presets[p.index].tracks = append(p.builder.presets[p.index].tracks, track)
	return p
}

// setWaveform sets the waveform of the last track.
func (p *Preset) setWaveform(waveform t.WaveformType) *Preset {
	track := p.lastTrack()
	if track == nil {
		return p
	}

	track.Waveform = waveform
	return p
}

// setEffect sets the effect type and value on the last track.
func (p *Preset) setEffect(effect t.EffectType, value float64) *Preset {
	track := p.lastTrack()
	if track == nil {
		return p
	}

	track.Effect.Type = effect
	track.Effect.Value = value
	return p
}

// lastTrack returns the last track in the current preset.
func (p *Preset) lastTrack() *t.Track {
	if p == nil || p.builder == nil || !p.valid() {
		return nil
	}

	tracks := p.builder.presets[p.index].tracks
	if len(tracks) == 0 {
		return nil
	}

	return &p.builder.presets[p.index].tracks[len(tracks)-1]
}

// valid returns whether the preset is valid (i.e. the index is within bounds and the name matches).
func (p *Preset) valid() bool {
	return p.index >= 0 && p.index < len(p.builder.presets) && p.builder.presets[p.index].name == p.name
}

// presetIndex returns the index of the given preset name, or -1 if not found.
func (b *Builder) presetIndex(name string) int {
	for i, preset := range b.presets {
		if preset.name == name {
			return i
		}
	}

	return -1
}
