// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"testing"

	clistyle "github.com/synapseq-foundation/synapseq/v4/internal/cli"
)

func TestOutputStylingHelpers(ts *testing.T) {
	clistyle.SetColorEnabled(false)
	defer clistyle.SetColorEnabled(true)

	message := clistyle.SuccessText("Dump generated:") + " " + clistyle.Accent("\"out.json\"")
	if message != "Dump generated: \"out.json\"" {
		ts.Fatalf("unexpected dump message formatting: %q", message)
	}

	comment := clistyle.Label(">") + " " + clistyle.Muted("focus block")
	if comment != "> focus block" {
		ts.Fatalf("unexpected comment formatting: %q", comment)
	}

	runHint := clistyle.Label("Run:") + " " + clistyle.Command("synapseq session.spsq")
	if runHint != "Run: synapseq session.spsq" {
		ts.Fatalf("unexpected run hint formatting: %q", runHint)
	}

	windowsMessage := clistyle.SuccessText("Successfully installed .spsq file association with SynapSeq.")
	if windowsMessage != "Successfully installed .spsq file association with SynapSeq." {
		ts.Fatalf("unexpected windows success formatting: %q", windowsMessage)
	}
}
