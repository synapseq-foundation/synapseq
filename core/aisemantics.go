// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package core

import (
	"fmt"
	"strings"
)

type aiState string

const (
	aiStateDelta aiState = "delta"
	aiStateTheta aiState = "theta"
	aiStateAlpha aiState = "alpha"
	aiStateBeta  aiState = "beta"
	aiStateGamma aiState = "gamma"
)

type aiIntent struct {
	state   aiState
	profile string
}

type aiRange struct {
	minimum float64
	maximum float64
}

type aiStage []aiState

var aiStateRanges = map[aiState]aiRange{
	aiStateDelta: {minimum: 0.5, maximum: 4},
	aiStateTheta: {minimum: 4, maximum: 8},
	aiStateAlpha: {minimum: 8, maximum: 13},
	aiStateBeta:  {minimum: 13, maximum: 30},
	aiStateGamma: {minimum: 30, maximum: 45},
}

func validateAISequenceSemantics(loaded *LoadedContext, prompt string) error {
	if err := validateAIBeatCarriers(loaded); err != nil {
		return err
	}

	intent := resolveAIIntent(prompt)
	if intent.profile != "" {
		return validateAIProfile(loaded, intent.profile)
	}
	if intent.state != "" {
		return validateAIState(loaded, intent.state)
	}

	return nil
}

func resolveAIIntent(prompt string) aiIntent {
	words := aiPromptWords(prompt)

	for _, profile := range []string{"sleep", "meditation", "focus", "relaxation"} {
		if words[profile] {
			return aiIntent{profile: profile}
		}
	}
	for _, state := range []aiState{aiStateDelta, aiStateTheta, aiStateAlpha, aiStateBeta, aiStateGamma} {
		if words[string(state)] {
			return aiIntent{state: state}
		}
	}

	return aiIntent{}
}

func aiPromptWords(prompt string) map[string]bool {
	words := map[string]bool{}
	for _, word := range strings.FieldsFunc(strings.ToLower(prompt), func(r rune) bool {
		return r < 'a' || r > 'z'
	}) {
		words[word] = true
	}

	return words
}

func validateAIBeatCarriers(loaded *LoadedContext) error {
	for _, preset := range loaded.Presets() {
		for _, track := range preset.Tracks {
			if !isAIBeatTrack(track.Type) {
				continue
			}
			if track.Carrier < 100 || track.Carrier > 600 {
				return fmt.Errorf("AI beat tracks require an audible carrier between 100 and 600 Hz; preset %q has %.2f Hz", preset.Name, track.Carrier)
			}
			if track.Resonance <= 0 {
				return fmt.Errorf("AI beat tracks require a beat frequency greater than zero; preset %q has %.2f Hz", preset.Name, track.Resonance)
			}
		}
	}

	return nil
}

func validateAIState(loaded *LoadedContext, target aiState) error {
	if hasAIState(loaded, target) {
		return nil
	}

	rangeValue := aiStateRanges[target]
	return fmt.Errorf("%s intent requires a beat frequency between %.1f and %.1f Hz", target, rangeValue.minimum, rangeValue.maximum)
}

func validateAIProfile(loaded *LoadedContext, profile string) error {
	stages := map[string][]aiStage{
		"sleep":      {{aiStateAlpha, aiStateTheta}, {aiStateDelta}},
		"meditation": {{aiStateAlpha}, {aiStateTheta}},
		"focus":      {{aiStateAlpha}, {aiStateBeta}},
		"relaxation": {{aiStateBeta}, {aiStateAlpha}},
	}[profile]

	if followsAIStages(loaded, stages) {
		return nil
	}

	return fmt.Errorf("%s profile requires %s progression", profile, aiStageDescription(stages))
}

func followsAIStages(loaded *LoadedContext, stages []aiStage) bool {
	stageIndex := 0
	for _, entry := range loaded.Timeline() {
		if entry.PresetName == "silence" {
			continue
		}
		states := aiPresetStates(loaded, entry.PresetName)
		for stageIndex < len(stages) && containsAnyAIState(states, stages[stageIndex]) {
			stageIndex++
		}
	}

	return stageIndex == len(stages)
}

func aiPresetStates(loaded *LoadedContext, name string) []aiState {
	states := []aiState{}
	for _, preset := range loaded.Presets() {
		if preset.Name != name {
			continue
		}
		for _, track := range preset.Tracks {
			if !isAIBeatTrack(track.Type) {
				continue
			}
			if state, ok := aiStateForBeat(track.Resonance); ok && !containsAIState(states, state) {
				states = append(states, state)
			}
		}
	}

	return states
}

func hasAIState(loaded *LoadedContext, target aiState) bool {
	for _, preset := range loaded.Presets() {
		for _, track := range preset.Tracks {
			if isAIBeatTrack(track.Type) && beatMatchesAIState(track.Resonance, target) {
				return true
			}
		}
	}

	return false
}

func isAIBeatTrack(trackType string) bool {
	return trackType == "binaural" || trackType == "monaural" || trackType == "isochronic"
}

func aiStateForBeat(beat float64) (aiState, bool) {
	for _, state := range []aiState{aiStateDelta, aiStateTheta, aiStateAlpha, aiStateBeta, aiStateGamma} {
		if beatMatchesAIState(beat, state) {
			return state, true
		}
	}

	return "", false
}

func beatMatchesAIState(beat float64, state aiState) bool {
	rangeValue := aiStateRanges[state]
	if state == aiStateGamma {
		return beat >= rangeValue.minimum && beat <= rangeValue.maximum
	}

	return beat >= rangeValue.minimum && beat < rangeValue.maximum
}

func containsAIState(states []aiState, target aiState) bool {
	for _, state := range states {
		if state == target {
			return true
		}
	}

	return false
}

func containsAnyAIState(states []aiState, targets []aiState) bool {
	for _, target := range targets {
		if containsAIState(states, target) {
			return true
		}
	}

	return false
}

func aiStageDescription(stages []aiStage) string {
	parts := make([]string, len(stages))
	for index, stage := range stages {
		stateNames := make([]string, len(stage))
		for stateIndex, state := range stage {
			stateNames[stateIndex] = string(state)
		}
		parts[index] = strings.Join(stateNames, " or ")
	}

	return strings.Join(parts, " to ")
}
