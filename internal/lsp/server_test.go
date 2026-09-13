// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package lsp

import (
	"slices"
	"strings"
	"testing"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	types "github.com/synapseq-foundation/synapseq/v4/internal/types"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

func TestSyntaxDiagnosticsContinueAfterInvalidLines(t *testing.T) {
	diagnostics := syntaxDiagnostics("@volum 80\n  tone\n00:00 invalid")
	if len(diagnostics) != 3 {
		t.Fatalf("diagnostics = %d, want 3", len(diagnostics))
	}
	for index, diagnostic := range diagnostics {
		if diagnostic.Range.Start.Line != uint32(index) {
			t.Fatalf("diagnostic %d line = %d, want %d", index, diagnostic.Range.Start.Line, index)
		}
	}
}

func TestValidateDocumentUsesSyntaxOnlyForUntitledBuffer(t *testing.T) {
	diagnostics := validateDocument(document{
		uri:  uri.MustParse("untitled:scratch.spsq"),
		text: "@ambiance rain audio/rain\n",
	})
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %v, want no filesystem diagnostics", diagnostics)
	}
}

func TestDiagnosticUsesUTF16Columns(t *testing.T) {
	diagnostic := toProtocolDiagnostic(
		parseDiagnostic(6, 7),
		"😀 bad",
	)
	if diagnostic.Range.Start.Character != 3 {
		t.Fatalf("start character = %d, want 3", diagnostic.Range.Start.Character)
	}
}

func TestCompletionIncludesDeclaredSymbols(t *testing.T) {
	items := complete("@ambiance rain audio/rain\nfocus\n00:00:00 ", protocol.Position{
		Line:      2,
		Character: 9,
	})
	labels := make([]string, 0, len(items))
	for _, item := range items {
		labels = append(labels, item.Label)
	}
	if !strings.Contains(strings.Join(labels, ","), "focus") {
		t.Fatalf("completion labels = %v, want focus", labels)
	}
}

func TestAmplitudeCompletionSupportsStereoForm(t *testing.T) {
	if got := trackCompletions([]string{types.KeywordTone, "220", types.KeywordAmplitude}, symbols{}); !slices.Equal(got, []string{types.KeywordLeft}) {
		t.Fatalf("amplitude completion = %v, want left", got)
	}
	if got := trackCompletions([]string{types.KeywordTone, "220", types.KeywordAmplitude, types.KeywordLeft, "20"}, symbols{}); !slices.Equal(got, []string{types.KeywordRight}) {
		t.Fatalf("left amplitude completion = %v, want right", got)
	}
}

func parseDiagnostic(column, endColumn int) *diag.Diagnostic {
	return diag.Parse("unexpected token").WithSpan(diag.Span{
		Column:    column,
		EndColumn: endColumn,
	})
}
