// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sbg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/synapseq-foundation/synapseq/v4/spsq"
)

const (
	initialFadeDuration       = 30 * time.Second
	openEndedSequenceDuration = 30 * time.Minute
)

type parsedTime struct {
	duration time.Duration
	initial  bool
	relative bool
}

func appendTimeline(
	builder *spsq.Builder,
	parsedSequence *sequence,
	presets map[string]*spsq.Preset,
	silentDefinitions map[string]struct{},
) error {
	if len(parsedSequence.timeline) == 0 {
		return fmt.Errorf("convert %q: no timeline entries found", parsedSequence.source)
	}
	firstEvent := parsedSequence.timeline[0]
	baseTime := firstEvent.at
	if firstEvent.initial {
		baseTime = 0
		builder.SilenceAt(0).Steady()
	}
	events := parsedSequence.timeline
	if len(events) == 1 {
		endingEvent := events[0]
		endingEvent.at = baseTime + openEndedSequenceDuration
		events = []timelineEvent{events[0], endingEvent}
	}
	var previousAt time.Duration
	for index, event := range events {
		at := event.at - baseTime
		if index == 0 && event.initial && event.at == 0 {
			at = initialFadeDuration
		}
		if index > 0 && at <= previousAt {
			return lineError(parsedSequence.source, event.line, "timeline entries must be strictly increasing")
		}
		previousAt = at
		if _, silent := silentDefinitions[event.name]; silent {
			if index == 0 && firstEvent.initial {
				continue
			}
			builder.SilenceAt(at).Steady()
			continue
		}
		preset := presets[event.name]
		if preset == nil {
			message := fmt.Sprintf("NameDef %q has no convertible voices", event.name)
			return lineError(parsedSequence.source, event.line, message)
		}
		builder.PresetAt(at, preset).Steady()
	}
	return nil
}

func parseTimeline(fields []string) (timelineEvent, error) {
	if len(fields) < 2 || len(fields) > 4 {
		return timelineEvent{}, errors.New("timeline must contain time, optional fade marker, NameDef, and optional ->")
	}
	timestamp, err := parseTime(fields[0])
	if err != nil {
		return timelineEvent{}, err
	}
	index := 1
	if isFadeMarker(fields[index]) {
		index++
	}
	if index >= len(fields) || !validName(fields[index]) {
		return timelineEvent{}, errors.New("timeline has an invalid or missing NameDef")
	}
	event := timelineEvent{
		at:       timestamp.duration,
		initial:  timestamp.initial,
		relative: timestamp.relative,
		name:     fields[index],
	}
	index++
	if index == len(fields) {
		return event, nil
	}
	if fields[index] != "->" || index != len(fields)-1 {
		return timelineEvent{}, fmt.Errorf("unknown timeline modifier %q", fields[index])
	}
	return event, nil
}

func parseTime(value string) (parsedTime, error) {
	timestamp := parsedTime{initial: value == "NOW" || strings.HasPrefix(value, "NOW+")}
	switch {
	case value == "NOW":
		return timestamp, nil
	case strings.HasPrefix(value, "NOW+"):
		value = strings.TrimPrefix(value, "NOW+")
	case strings.HasPrefix(value, "+"):
		value = strings.TrimPrefix(value, "+")
		timestamp.relative = true
	}

	parts := strings.Split(value, ":")
	if len(parts) != 2 && len(parts) != 3 {
		return parsedTime{}, fmt.Errorf("invalid timeline time %q", value)
	}
	values := make([]int, len(parts))
	for index, part := range parts {
		if index == 0 {
			if len(part) < 1 || len(part) > 2 {
				return parsedTime{}, errors.New("timeline hour must have one or two digits")
			}
		} else if len(part) != 2 {
			return parsedTime{}, errors.New("timeline minutes and seconds must have two digits")
		}
		fieldValue, err := strconv.Atoi(part)
		if err != nil || fieldValue < 0 {
			return parsedTime{}, fmt.Errorf("invalid timeline time %q", value)
		}
		values[index] = fieldValue
	}
	hours := values[0]
	minutes := values[1]
	if minutes >= 60 {
		return parsedTime{}, errors.New("timeline minutes and seconds must be below 60")
	}
	duration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
	if len(values) == 3 {
		seconds := values[2]
		if seconds >= 60 {
			return parsedTime{}, errors.New("timeline minutes and seconds must be below 60")
		}
		duration += time.Duration(seconds) * time.Second
	}
	timestamp.duration = duration
	return timestamp, nil
}

func isTimelineTime(value string) bool {
	if value == "NOW" || strings.HasPrefix(value, "NOW+") || strings.HasPrefix(value, "+") {
		return true
	}
	if value == "" {
		return false
	}
	first := value[0]
	startsWithDigit := first >= '0' && first <= '9'
	startsWithNegativeDigit := first == '-' && len(value) > 1 && value[1] >= '0' && value[1] <= '9'
	return startsWithDigit || startsWithNegativeDigit
}

func isFadeMarker(value string) bool {
	if len(value) != 2 {
		return false
	}
	return strings.ContainsRune("<-=", rune(value[0])) && strings.ContainsRune(">-=", rune(value[1]))
}
