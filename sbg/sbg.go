// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sbg

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	"github.com/synapseq-foundation/synapseq/v4/spsq"
)

// Converter converts SBaGen content through the provided SynapSeq application context.
type Converter struct {
	ctx *synapseq.AppContext
}

// New creates an SBaGen converter using ctx.
func New(ctx *synapseq.AppContext) (*Converter, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}
	return &Converter{ctx: ctx}, nil
}

// LoadFile converts the SBaGen sequence at path into a validated SynapSeq sequence.
func (c *Converter) LoadFile(path string) (*synapseq.LoadedContext, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open SBaGen file %q: %w", path, err)
	}
	defer file.Close()

	return c.load(path, file)
}

// LoadContent converts in-memory SBaGen content into a validated SynapSeq sequence.
func (c *Converter) LoadContent(content string) (*synapseq.LoadedContext, error) {
	return c.load("<content>", strings.NewReader(content))
}

func (c *Converter) load(source string, reader io.Reader) (*synapseq.LoadedContext, error) {
	parsed, err := parse(source, reader)
	if err != nil {
		return nil, err
	}
	builder, err := build(parsed, source)
	if err != nil {
		return nil, err
	}
	loaded, err := builder.Load(c.ctx)
	if err != nil {
		return nil, fmt.Errorf("validate converted sequence: %w", err)
	}
	return loaded, nil
}

func build(parsed *sequence, inputPath string) (*spsq.Builder, error) {
	builder := spsq.New()

	presets := make(map[string]*spsq.Preset, len(parsed.definitions))
	names := convertedNames(parsed.definitions)
	allOff := make(map[string]bool)
	for _, definition := range parsed.definitions {
		active := 0
		for _, voice := range definition.voices {
			if voice.kind != voiceOff && voice.kind != voiceMix {
				active++
			}
		}
		if active == 0 {
			allOff[definition.name] = true
			continue
		}

		preset := builder.NewPreset(names[definition.name])
		for _, voice := range definition.voices {
			switch voice.kind {
			case voiceOff:
				continue
			case voicePink:
				preset.Pink(0).Amplitude(voice.amplitude)
			case voiceTone:
				preset.Tone(voice.carrier).Amplitude(voice.amplitude)
			case voiceBinaural:
				preset.Tone(voice.carrier).Binaural(voice.beat).Amplitude(voice.amplitude)
			case voiceSpin:
				preset.Pink(0).Pan(voice.beat).Intensity(100).Amplitude(voice.amplitude)
			case voiceMix:
				// SBaGen soundtrack input has no direct conversion target.
				continue
			}
		}
		presets[definition.name] = preset
	}

	if len(parsed.timeline) == 0 {
		return nil, fmt.Errorf("convert %q: no timeline entries found", inputPath)
	}
	firstEvent := parsed.timeline[0]
	baseTime := firstEvent.at
	if firstEvent.initial {
		baseTime = 0
		builder.SilenceAt(0).Steady()
	}
	for index, event := range parsed.timeline {
		at := event.at - baseTime
		if index == 0 && event.initial && event.at == 0 {
			at = 30 * time.Second
		}
		if index > 0 && at <= timelineAt(parsed.timeline[index-1], index-1, baseTime) {
			return nil, lineError(parsed.source, event.line, "timeline entries must be strictly increasing")
		}
		if allOff[event.name] {
			if index == 0 && firstEvent.initial {
				continue
			}
			builder.SilenceAt(at).Steady()
			continue
		}
		preset := presets[event.name]
		if preset == nil {
			return nil, lineError(parsed.source, event.line, fmt.Sprintf("NameDef %q has no convertible voices", event.name))
		}
		builder.PresetAt(at, preset).Steady()
	}
	return builder, nil
}

func timelineAt(event timelineEvent, index int, baseTime time.Duration) time.Duration {
	at := event.at - baseTime
	if index == 0 && event.initial && event.at == 0 {
		return 30 * time.Second
	}
	return at
}

func convertedNames(definitions []nameDef) map[string]string {
	result := make(map[string]string, len(definitions))
	used := map[string]struct{}{"silence": {}}
	for _, definition := range definitions {
		var name strings.Builder
		for _, char := range strings.ToLower(definition.name) {
			if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '_' || char == '-' {
				name.WriteRune(char)
			} else {
				name.WriteByte('-')
			}
		}
		base := name.String()
		if len(base) > 20 {
			base = base[:20]
		}
		candidate := base
		for suffix := 2; ; suffix++ {
			if _, exists := used[candidate]; !exists {
				break
			}
			ending := "-" + strconv.Itoa(suffix)
			candidate = base
			if len(candidate)+len(ending) > 20 {
				candidate = candidate[:20-len(ending)]
			}
			candidate += ending
		}
		used[candidate] = struct{}{}
		result[definition.name] = candidate
	}
	return result
}
