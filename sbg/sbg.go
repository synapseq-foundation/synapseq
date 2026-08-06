// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sbg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	builder, err := build(c.ctx, parsed)
	if err != nil {
		return nil, err
	}
	loaded, err := builder.Load()
	if err != nil {
		return nil, fmt.Errorf("validate converted sequence: %w", err)
	}
	return loaded, nil
}

func build(ctx *synapseq.AppContext, parsed *sequence) (*spsq.Builder, error) {
	builder, err := spsq.New(ctx)
	if err != nil {
		return nil, err
	}
	builder.SampleRate(parsed.sampleRate)
	var musicName string
	if parsed.musicPath != "" {
		musicPath, err := currentRelativeMusicPath(parsed.musicPath)
		if err != nil {
			return nil, lineError(parsed.source, parsed.musicLine, err.Error())
		}
		musicName = filepath.Base(musicPath)
		builder.Music(musicName, musicPath)
	}

	presets, allOff, err := convertDefinitions(builder, parsed, musicName)
	if err != nil {
		return nil, err
	}
	if err := appendTimeline(builder, parsed, presets, allOff); err != nil {
		return nil, err
	}

	return builder, nil
}

func convertDefinitions(
	builder *spsq.Builder,
	parsed *sequence,
	musicName string,
) (map[string]*spsq.Preset, map[string]struct{}, error) {
	presets := make(map[string]*spsq.Preset, len(parsed.definitions))
	names := convertedNames(parsed.definitions)
	allOff := make(map[string]struct{})
	for _, definition := range parsed.definitions {
		hasTrack := false
		hasOff := false
		for _, voice := range definition.voices {
			if voice.kind == voiceMix && parsed.musicPath == "" {
				return nil, nil, lineError(parsed.source, definition.line, "mix voice requires a -m music source")
			}
			if voice.kind == voiceOff {
				hasOff = true
				continue
			}
			hasTrack = true
		}
		if !hasTrack {
			if hasOff {
				allOff[definition.name] = struct{}{}
				continue
			}
			return nil, nil, lineError(parsed.source, definition.line, fmt.Sprintf("NameDef %q has no convertible voices", definition.name))
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
				preset.Music(musicName).Amplitude(voice.amplitude)
			}
		}
		presets[definition.name] = preset
	}
	return presets, allOff, nil
}

func appendTimeline(
	builder *spsq.Builder,
	parsed *sequence,
	presets map[string]*spsq.Preset,
	allOff map[string]struct{},
) error {
	if len(parsed.timeline) == 0 {
		return fmt.Errorf("convert %q: no timeline entries found", parsed.source)
	}
	firstEvent := parsed.timeline[0]
	baseTime := firstEvent.at
	if firstEvent.initial {
		baseTime = 0
		builder.SilenceAt(0).Steady()
	}
	var previousAt time.Duration
	for index, event := range parsed.timeline {
		at := event.at - baseTime
		if index == 0 && event.initial && event.at == 0 {
			at = 30 * time.Second
		}
		if index > 0 && at <= previousAt {
			return lineError(parsed.source, event.line, "timeline entries must be strictly increasing")
		}
		previousAt = at
		if _, off := allOff[event.name]; off {
			if index == 0 && firstEvent.initial {
				continue
			}
			builder.SilenceAt(at).Steady()
			continue
		}
		preset := presets[event.name]
		if preset == nil {
			return lineError(parsed.source, event.line, fmt.Sprintf("NameDef %q has no convertible voices", event.name))
		}
		builder.PresetAt(at, preset).Steady()
	}
	return nil
}

func currentRelativeMusicPath(value string) (string, error) {
	if strings.Contains(value, "://") {
		return "", fmt.Errorf("music source must be a local path")
	}
	if strings.Contains(value, "\\") {
		return "", fmt.Errorf("music source must use '/' path separators")
	}
	if hasParentTraversal(value) {
		return "", fmt.Errorf("music source must not contain parent directory traversal")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get current directory: %w", err)
	}
	absolutePath, err := filepath.Abs(value)
	if err != nil {
		return "", fmt.Errorf("resolve music source %q: %w", value, err)
	}
	relativePath, err := filepath.Rel(cwd, absolutePath)
	if err != nil {
		return "", fmt.Errorf("relativize music source %q: %w", value, err)
	}
	if relativePath == ".." || strings.HasPrefix(relativePath, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("music source must be inside the current directory")
	}

	basePath := strings.TrimSuffix(relativePath, filepath.Ext(relativePath))
	if basePath == "." || basePath == "" {
		return "", fmt.Errorf("music source must name a file")
	}
	return filepath.ToSlash(basePath), nil
}

func hasParentTraversal(path string) bool {
	for _, part := range strings.Split(path, "/") {
		if part == ".." {
			return true
		}
	}
	return false
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
