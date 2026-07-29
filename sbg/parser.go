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
	result := &sequence{source: source}
	definitions := make(map[string]struct{})
	var relativeBase time.Duration
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := strings.TrimSuffix(scanner.Text(), "\r")
		if strings.IndexByte(line, 0) >= 0 {
			return nil, lineError(source, lineNumber, "NUL byte is not allowed")
		}
		if before, _, found := strings.Cut(line, "#"); found {
			line = before
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		musicPath, ok, err := parseMusicOption(fields)
		if err != nil {
			return nil, lineError(source, lineNumber, err.Error())
		}
		if ok {
			result.musicPath = musicPath
			result.musicLine = lineNumber
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
		if _, exists := definitions[name]; exists {
			return nil, lineError(source, lineNumber, fmt.Sprintf("duplicate NameDef %q", name))
		}

		definition := nameDef{name: name, line: lineNumber}
		for _, token := range strings.Fields(voicesText) {
			parsed, err := parseVoice(token)
			if err != nil {
				return nil, lineError(source, lineNumber, fmt.Sprintf("voice %q: %v", token, err))
			}
			definition.voices = append(definition.voices, parsed)
		}
		if len(definition.voices) == 0 {
			return nil, lineError(source, lineNumber, fmt.Sprintf("NameDef %q has no voices", name))
		}
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
		if _, ok := definitions[event.name]; !ok {
			return nil, lineError(source, event.line, fmt.Sprintf("timeline references undefined NameDef %q", event.name))
		}
	}
	return result, nil
}
