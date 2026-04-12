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
			Steps: 1,
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
		"Transition smooth - 1 step",
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

func TestStepAlphaForPreview(ts *testing.T) {
	period := t.Period{Transition: t.TransitionSteady, Steps: 1}
	if got := stepAlphaForPreview(1.0/3.0, period); math.Abs(got-1.0) > 0.000001 {
		ts.Fatalf("expected first leg to reach end value, got %.6f", got)
	}
	if got := stepAlphaForPreview(2.0/3.0, period); math.Abs(got-0.0) > 0.000001 {
		ts.Fatalf("expected second leg to return to start value, got %.6f", got)
	}
	if got := stepAlphaForPreview(1.0, period); math.Abs(got-1.0) > 0.000001 {
		ts.Fatalf("expected final point to end on the next preset, got %.6f", got)
	}
}

func TestBuildGraphMetrics(ts *testing.T) {
	periods := []t.Period{
		{
			Time: 0,
			TrackStart: [t.NumberOfChannels]t.Track{
				{Type: t.TrackBinauralBeat, Carrier: 220, Resonance: 8, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformSine},
				{Type: t.TrackPinkNoise, NoiseSmooth: 30, Amplitude: t.AmplitudePercentToRaw(35)},
			},
			TrackEnd: [t.NumberOfChannels]t.Track{
				{Type: t.TrackBinauralBeat, Carrier: 240, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(25), Waveform: t.WaveformSquare},
				{Type: t.TrackPinkNoise, NoiseSmooth: 40, Amplitude: t.AmplitudePercentToRaw(20)},
			},
			Transition: t.TransitionSmooth,
			Steps:      1,
		},
		{
			Time:       60000,
			TrackStart: [t.NumberOfChannels]t.Track{{Type: t.TrackBinauralBeat, Carrier: 240, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(25), Waveform: t.WaveformSquare}},
		},
	}

	metrics := buildGraphMetrics(periods, 60000)
	if len(metrics) != 7 {
		ts.Fatalf("expected 7 metrics, got %d", len(metrics))
	}

	keys := make(map[string]previewGraphMetricView, len(metrics))
	for _, metric := range metrics {
		keys[metric.Key] = metric
	}

	for _, key := range []string{"resonance", "carrier", "waveform", "amplitude", "smooth", "effect", "intensity"} {
		if _, ok := keys[key]; !ok {
			ts.Fatalf("expected graph metric %q", key)
		}
	}

	if !keys["carrier"].HasData {
		ts.Fatalf("expected carrier metric to have data")
	}
	if !keys["waveform"].HasData {
		ts.Fatalf("expected waveform metric to have data")
	}
	if !keys["smooth"].HasData {
		ts.Fatalf("expected smooth metric to have data")
	}
}

func TestResolveGraphTrackFromSilenceUsesPreviousVisibleTrack(ts *testing.T) {
	periods := []t.Period{
		{
			Time: 0,
			TrackStart: [t.NumberOfChannels]t.Track{
				{Type: t.TrackBinauralBeat, Carrier: 220, Resonance: 8, Amplitude: t.AmplitudePercentToRaw(20)},
			},
		},
		{
			Time: 60000,
			TrackStart: [t.NumberOfChannels]t.Track{
				{Type: t.TrackSilence, Carrier: 180, Resonance: 6, Amplitude: t.AmplitudePercentToRaw(10)},
			},
		},
	}

	track, ok := resolveGraphTrack(periods, 1, 0)
	if !ok {
		ts.Fatalf("expected silent graph track to resolve from previous visible track")
	}
	if track.Type != t.TrackBinauralBeat {
		ts.Fatalf("expected inherited type %v, got %v", t.TrackBinauralBeat, track.Type)
	}
	if track.Carrier != 180 {
		ts.Fatalf("expected carrier to use silent boundary value, got %.2f", track.Carrier)
	}
	if track.Resonance != 6 {
		ts.Fatalf("expected resonance to use silent boundary value, got %.2f", track.Resonance)
	}
}

func TestBuildTrackViewAndSegmentItemView(ts *testing.T) {
	track := t.Track{
		Type:         t.TrackAmbiance,
		Carrier:      180,
		Resonance:    7,
		Waveform:     t.WaveformTriangle,
		AmbianceName: "rain",
		Amplitude:    t.AmplitudePercentToRaw(45),
		Effect: t.Effect{
			Type:      t.EffectPan,
			Value:     1.5,
			Intensity: t.IntensityPercentToRaw(60),
		},
	}

	view := buildTrackView(2, track)
	if view.ChannelLabel != "CH 03" {
		ts.Fatalf("expected channel label CH 03, got %q", view.ChannelLabel)
	}
	if view.TypeLabel != "Ambiance" {
		ts.Fatalf("expected type label Ambiance, got %q", view.TypeLabel)
	}
	if view.Summary != "Ambiance layer \"rain\"" {
		ts.Fatalf("unexpected track summary %q", view.Summary)
	}

	labels := make(map[string]string, len(view.Meta))
	for _, item := range view.Meta {
		labels[item.Label] = item.Value
	}

	for _, expected := range []string{"Carrier", "Beat", "Waveform", "Ambiance", "Amplitude", "Effect", "Effect value", "Intensity"} {
		if _, ok := labels[expected]; !ok {
			ts.Fatalf("expected metadata item %q", expected)
		}
	}

	segment := buildSegmentItemView(2, t.Track{Type: t.TrackSilence}, track)
	if segment.ChannelLabel != "CH 03" {
		ts.Fatalf("expected segment channel label CH 03, got %q", segment.ChannelLabel)
	}
	if segment.Label != "Ambiance" {
		ts.Fatalf("expected segment label Ambiance, got %q", segment.Label)
	}
	if segment.Class != "track-ambiance" {
		ts.Fatalf("expected segment class track-ambiance, got %q", segment.Class)
	}
}
