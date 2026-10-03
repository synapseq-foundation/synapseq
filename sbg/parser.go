// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sbg

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

func parse(source string, reader io.Reader) (*sequence, error) {
	result := &sequence{source: source, sampleRate: defaultSampleRate}
	definitions := make(map[string]struct{})
	// Relative entries are offsets from the last absolute entry, not the previous entry.
	var relativeBase time.Duration
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := strings.TrimSuffix(scanner.Text(), "\r")
		if strings.IndexByte(line, 0) >= 0 {
			return nil, lineError(source, lineNumber, "NUL byte is not allowed")
		}
		if comment := strings.TrimLeft(line, " \t"); strings.HasPrefix(comment, "#") {
			result.comments = append(result.comments, comment)
			continue
		}
		if before, _, found := strings.Cut(line, "#"); found {
			line = before
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		musicPath, hasMusic, err := parseMusicOption(fields)
		if err != nil {
			return nil, lineError(source, lineNumber, err.Error())
		}
		sampleRate, hasSampleRate, err := parseSampleRateOption(fields)
		if err != nil {
			return nil, lineError(source, lineNumber, err.Error())
		}
		if hasMusic {
			result.musicPath = musicPath
			result.musicLine = lineNumber
		}
		if hasSampleRate {
			result.sampleRate = sampleRate
		}
		if hasMusic || hasSampleRate {
			continue
		}

		if isTimelineTime(fields[0]) {
			event, err := parseTimeline(fields)
			if err != nil {
				return nil, lineError(source, lineNumber, err.Error())
			}
			event.line = lineNumber
			if event.relative {
				event.at += relativeBase
			} else {
				relativeBase = event.at
			}
			result.timeline = append(result.timeline, event)
			continue
		}

		nameText, voicesText, found := strings.Cut(line, ":")
		if !found {
			continue
		}
		name := strings.TrimSpace(nameText)
		if !validName(name) {
			continue
		}
		if strings.TrimSpace(voicesText) == "{" {
			return nil, lineError(source, lineNumber, "block definitions are not supported by SynapSeq")
		}
		if _, exists := definitions[name]; exists {
			return nil, lineError(source, lineNumber, fmt.Sprintf("duplicate NameDef %q", name))
		}

		definition, err := parseDefinition(name, voicesText)
		if err != nil {
			return nil, lineError(source, lineNumber, err.Error())
		}
		definition.line = lineNumber
		definitions[name] = struct{}{}
		result.definitions = append(result.definitions, definition)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan %q: %w", source, err)
	}
	if len(result.definitions) == 0 {
		return nil, fmt.Errorf("parse %q: no NameDefs found", source)
	}
	if len(result.timeline) == 0 {
		return nil, fmt.Errorf("parse %q: no timeline entries found", source)
	}
	for _, event := range result.timeline {
		if _, exists := definitions[event.name]; !exists {
			message := fmt.Sprintf("timeline references undefined NameDef %q", event.name)
			return nil, lineError(source, event.line, message)
		}
	}
	return result, nil
}

func lineError(source string, line int, message string) error {
	return fmt.Errorf("parse %q line %d: %s", source, line, message)
}
