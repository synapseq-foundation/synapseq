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
)

func parseTimeline(fields []string) (timelineEvent, error) {
	if len(fields) < 2 || len(fields) > 4 {
		return timelineEvent{}, errors.New("timeline must contain time, optional fade marker, NameDef, and optional ->")
	}
	at, initial, relative, err := parseTime(fields[0])
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
	event := timelineEvent{at: at, initial: initial, relative: relative, name: fields[index]}
	index++
	if index < len(fields) && (fields[index] != "->" || index != len(fields)-1) {
		return timelineEvent{}, fmt.Errorf("unknown timeline modifier %q", fields[index])
	}
	return event, nil
}

func parseTime(value string) (time.Duration, bool, bool, error) {
	initial := value == "NOW" || strings.HasPrefix(value, "NOW+")
	relative := false
	switch {
	case value == "NOW":
		return 0, true, false, nil
	case strings.HasPrefix(value, "NOW+"):
		value = strings.TrimPrefix(value, "NOW+")
	case strings.HasPrefix(value, "+"):
		value = strings.TrimPrefix(value, "+")
		relative = true
	}

	parts := strings.Split(value, ":")
	if len(parts) != 2 && len(parts) != 3 {
		return 0, false, false, fmt.Errorf("invalid timeline time %q", value)
	}
	values := make([]int, len(parts))
	for index, part := range parts {
		if len(part) != 2 {
			return 0, false, false, errors.New("timeline time fields must have two digits")
		}
		parsed, err := strconv.Atoi(part)
		if err != nil {
			return 0, false, false, fmt.Errorf("invalid timeline time %q", value)
		}
		values[index] = parsed
	}
	minutes := values[len(values)-2]
	seconds := values[len(values)-1]
	if minutes >= 60 || seconds >= 60 {
		return 0, false, false, errors.New("timeline minutes and seconds must be below 60")
	}
	duration := time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
	if len(values) == 3 {
		duration += time.Duration(values[0]) * time.Hour
	}
	return duration, initial, relative, nil
}

func isTimelineTime(value string) bool {
	if value == "NOW" || strings.HasPrefix(value, "NOW+") || strings.HasPrefix(value, "+") {
		return true
	}
	_, _, _, err := parseTime(value)
	return err == nil
}

func isFadeMarker(value string) bool {
	return len(value) == 2 && strings.ContainsRune("<-=", rune(value[0])) && strings.ContainsRune(">-=", rune(value[1]))
}
