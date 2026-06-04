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

	message := clistyle.SuccessText("Preview generated:") + " " + clistyle.Accent("\"out.html\"")
	if message != "Preview generated: \"out.html\"" {
		ts.Fatalf("unexpected preview message formatting: %q", message)
	}

	comment := clistyle.Label(">") + " " + clistyle.Muted("focus block")
	if comment != "> focus block" {
		ts.Fatalf("unexpected comment formatting: %q", comment)
	}

	templateMessage := clistyle.SuccessText("Template generated:") + " " + clistyle.Accent("\"meditation\"") + " " + clistyle.Muted("as \"session.spsq\"")
	if templateMessage != "Template generated: \"meditation\" as \"session.spsq\"" {
		ts.Fatalf("unexpected template message formatting: %q", templateMessage)
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
