// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package core

import (
	"strings"
	"testing"
)

func TestResolveAIIntentRecognizesEnglishTerms(t *testing.T) {
	intent := resolveAIIntent("Create a gamma session")
	if intent.state != aiStateGamma || intent.profile != "" {
		t.Fatalf("unexpected gamma intent: %#v", intent)
	}

	intent = resolveAIIntent("Create a focus sequence")
	if intent.profile != "focus" {
		t.Fatalf("unexpected focus intent: %#v", intent)
	}

	intent = resolveAIIntent("Criar uma sequencia para foco")
	if intent.state != "" || intent.profile != "" {
		t.Fatalf("expected no non-English intent, got %#v", intent)
	}
}

func TestValidateAISequenceSemanticsRejectsLowGammaCarrier(t *testing.T) {
	loaded, err := NewAppContext().LoadContent(`
gamma
  tone 40 binaural 40 amplitude 10

00:00:00 gamma
00:05:00 gamma
`)
	if err != nil {
		t.Fatalf("LoadContent error: %v", err)
	}

	err = validateAISequenceSemantics(loaded, "Create a gamma session")
	if err == nil || !strings.Contains(err.Error(), "audible carrier between 100 and 600 Hz") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateAISequenceSemanticsAcceptsGammaCarrierAndBeat(t *testing.T) {
	loaded, err := NewAppContext().LoadContent(`
gamma
  tone 220 binaural 40 amplitude 10

00:00:00 gamma
00:05:00 gamma
`)
	if err != nil {
		t.Fatalf("LoadContent error: %v", err)
	}

	if err := validateAISequenceSemantics(loaded, "Create a gamma session"); err != nil {
		t.Fatalf("validateAISequenceSemantics error: %v", err)
	}
}

func TestValidateAISequenceSemanticsAcceptsSleepProgression(t *testing.T) {
	loaded, err := NewAppContext().LoadContent(`
sleep-entry
  tone 220 binaural 6 amplitude 10

sleep-deep
  tone 220 binaural 2 amplitude 10

00:00:00 sleep-entry
00:05:00 sleep-deep
00:10:00 sleep-deep
`)
	if err != nil {
		t.Fatalf("LoadContent error: %v", err)
	}

	if err := validateAISequenceSemantics(loaded, "Create a sleep sequence"); err != nil {
		t.Fatalf("validateAISequenceSemantics error: %v", err)
	}
}

func TestValidateAISequenceSemanticsRejectsIncompleteFocusProgression(t *testing.T) {
	loaded, err := NewAppContext().LoadContent(`
focus
  tone 220 binaural 10 amplitude 10

00:00:00 focus
00:10:00 focus
`)
	if err != nil {
		t.Fatalf("LoadContent error: %v", err)
	}

	err = validateAISequenceSemantics(loaded, "Create a focus sequence")
	if err == nil || !strings.Contains(err.Error(), "focus profile requires alpha to beta progression") {
		t.Fatalf("unexpected error: %v", err)
	}
}
