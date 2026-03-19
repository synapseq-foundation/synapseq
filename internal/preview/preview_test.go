package preview

import (
	"math"
	"strings"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestGetPreviewContent(ts *testing.T) {
	periods := []t.Period{
		{
			Time: 0,
			TrackStart: [t.NumberOfChannels]t.Track{
				{
					Type:        t.TrackPinkNoise,
					NoiseSmooth: 35,
					Amplitude:   t.AmplitudePercentToRaw(35),
				},
				{
					Type:      t.TrackBinauralBeat,
					Carrier:   200,
					Resonance: 10,
					Amplitude: t.AmplitudePercentToRaw(15),
					Waveform:  t.WaveformSine,
					Effect: t.Effect{
						Type:      t.EffectPan,
						Value:     2.5,
						Intensity: t.IntensityPercentToRaw(75),
					},
				},
				{
					Type:         t.TrackAmbiance,
					AmbianceName: "river",
					Waveform:     t.WaveformSquare,
					Amplitude:    t.AmplitudePercentToRaw(40),
				},
			},
			TrackEnd: [t.NumberOfChannels]t.Track{
				{
					Type:        t.TrackPinkNoise,
					NoiseSmooth: 10,
					Amplitude:   0,
				},
				{
					Type:      t.TrackSilence,
					Carrier:   200,
					Resonance: 10,
					Amplitude: 0,
					Waveform:  t.WaveformSawtooth,
					Effect: t.Effect{
						Type:      t.EffectPan,
						Value:     2.5,
						Intensity: t.IntensityPercentToRaw(75),
					},
				},
				{
					Type:         t.TrackSilence,
					AmbianceName: "river",
					Waveform:     t.WaveformSquare,
					Amplitude:    0,
				},
			},
			Transition: t.TransitionSmooth,
		},
		{
			Time: 300000,
			TrackStart: [t.NumberOfChannels]t.Track{
				{
					Type:      t.TrackSilence,
					Amplitude: 0,
				},
				{
					Type:      t.TrackSilence,
					Carrier:   200,
					Resonance: 10,
					Amplitude: 0,
					Waveform:  t.WaveformSine,
					Effect: t.Effect{
						Type:      t.EffectPan,
						Value:     2.5,
						Intensity: t.IntensityPercentToRaw(75),
					},
				},
			},
			TrackEnd:   [t.NumberOfChannels]t.Track{},
			Transition: t.TransitionSteady,
		},
	}

	content, err := GetPreviewContent(periods)
	if err != nil {
		ts.Fatalf("unexpected error rendering preview: %v", err)
	}

	html := string(content)
	checks := []string{
		"SynapSeq Sequence Preview",
		"Frequency timeline",
		"00:05:00",
		"CH 01 Pink noise",
		"CH 02 Binaural beat",
		"CH 03 Ambiance",
		"Binaural beat",
		"Pink noise",
		"Ambiance",
		"Resonance",
		"Carrier",
		"Waveform",
		"Sine",
		"Square",
		"Sawtooth",
		"CH 02 Binaural beat • Sine",
		"CH 03 Ambiance • Square",
		"Carrier",
		"Beat",
		"Smooth",
		"Amplitude",
		"Effect",
		"Effect value",
		"2.50",
		"pan",
		"Intensity",
		"75.00%",
		"35.00%",
		"data-graph-target=\"resonance\"",
		"data-graph-target=\"carrier\"",
		"data-graph-target=\"waveform\"",
		"data-graph-target=\"amplitude\"",
		"data-graph-target=\"smooth\"",
		"data-graph-target=\"effect\"",
		"data-graph-target=\"intensity\"",
		"0.00%",
	}

	for _, expected := range checks {
		if !strings.Contains(html, expected) {
			ts.Fatalf("expected HTML preview to contain %q", expected)
		}
	}

	if strings.Contains(html, "CH 01 Silence") || strings.Contains(html, "CH 02 Silence") {
		ts.Fatalf("expected final node to reuse previous track labels instead of silence")
	}
}

func TestApplyTransitionAlpha(ts *testing.T) {
	steady := applyTransitionAlpha(0.5, t.TransitionSteady)
	if steady != 0.5 {
		ts.Fatalf("expected steady alpha to remain linear, got %.6f", steady)
	}

	easeOut := applyTransitionAlpha(0.5, t.TransitionEaseOut)
	if easeOut <= 0.5 {
		ts.Fatalf("expected ease-out alpha to be ahead of linear progress, got %.6f", easeOut)
	}

	easeIn := applyTransitionAlpha(0.5, t.TransitionEaseIn)
	if easeIn >= 0.5 {
		ts.Fatalf("expected ease-in alpha to lag behind linear progress, got %.6f", easeIn)
	}

	smooth := applyTransitionAlpha(0.5, t.TransitionSmooth)
	if math.Abs(smooth-0.5) > 0.000001 {
		ts.Fatalf("expected smooth alpha midpoint to stay centered, got %.6f", smooth)
	}
}
