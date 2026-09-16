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

func TestInitializeAdvertisesSemanticTokens(t *testing.T) {
	result, err := (&server{}).Initialize(t.Context(), &protocol.InitializeParams{})
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	provider, ok := result.Capabilities.SemanticTokensProvider.(*protocol.SemanticTokensOptions)
	if !ok {
		t.Fatalf("semantic token provider = %T, want *protocol.SemanticTokensOptions", result.Capabilities.SemanticTokensProvider)
	}
	if !slices.Equal(provider.Legend.TokenTypes, semanticTokenLegend) {
		t.Fatalf("semantic token legend = %v, want %v", provider.Legend.TokenTypes, semanticTokenLegend)
	}
}

func TestSemanticTokensClassifySPSQSyntax(t *testing.T) {
	data := semanticTokenData("# sessão 😀\n@ambiance rain audio/rain\nfocus\n  tone 220 binaural 8 amplitude left 20 right 30\n00:00:00 focus smooth 4")
	got := decodeSemanticTokenData(data)
	want := []semanticToken{
		{line: 0, start: 0, length: 11, tokenType: semanticTokenComment},
		{line: 1, start: 0, length: 9, tokenType: semanticTokenKeyword},
		{line: 1, start: 10, length: 4, tokenType: semanticTokenVariable},
		{line: 1, start: 15, length: 10, tokenType: semanticTokenString},
		{line: 2, start: 0, length: 5, tokenType: semanticTokenVariable},
		{line: 3, start: 2, length: 4, tokenType: semanticTokenKeyword},
		{line: 3, start: 7, length: 3, tokenType: semanticTokenNumber},
		{line: 3, start: 11, length: 8, tokenType: semanticTokenKeyword},
		{line: 3, start: 20, length: 1, tokenType: semanticTokenNumber},
		{line: 3, start: 22, length: 9, tokenType: semanticTokenParameter},
		{line: 3, start: 32, length: 4, tokenType: semanticTokenParameter},
		{line: 3, start: 37, length: 2, tokenType: semanticTokenNumber},
		{line: 3, start: 40, length: 5, tokenType: semanticTokenParameter},
		{line: 3, start: 46, length: 2, tokenType: semanticTokenNumber},
		{line: 4, start: 0, length: 8, tokenType: semanticTokenNumber},
		{line: 4, start: 9, length: 5, tokenType: semanticTokenVariable},
		{line: 4, start: 15, length: 6, tokenType: semanticTokenKeyword},
		{line: 4, start: 22, length: 1, tokenType: semanticTokenNumber},
	}
	if !slices.Equal(got, want) {
		t.Fatalf("semantic tokens = %#v, want %#v", got, want)
	}
}

func TestSemanticTokensUseUTF16Columns(t *testing.T) {
	tokens := semanticTokens("  😀 tone 220")
	if tokens[0].start != 2 || tokens[0].length != 2 {
		t.Fatalf("emoji token = %#v, want UTF-16 width 2", tokens[0])
	}
	if tokens[1].start != 5 {
		t.Fatalf("tone start = %d, want 5", tokens[1].start)
	}
}

func TestSemanticTokensClassifyOptionPathsWithoutSeparators(t *testing.T) {
	tokens := semanticTokens("@extends common\n@music ocean library")
	if tokens[1].tokenType != semanticTokenString {
		t.Fatalf("extends path token type = %d, want string", tokens[1].tokenType)
	}
	if tokens[4].tokenType != semanticTokenString {
		t.Fatalf("music path token type = %d, want string", tokens[4].tokenType)
	}
}

func TestSemanticTokensFullUsesOpenDocument(t *testing.T) {
	documentURI := uri.MustParse("untitled:semantic.spsq")
	server := &server{
		documents: map[uri.URI]document{
			documentURI: {
				uri:  documentURI,
				text: "focus",
			},
		},
	}
	result, err := server.SemanticTokensFull(t.Context(), &protocol.SemanticTokensParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: documentURI},
	})
	if err != nil {
		t.Fatalf("SemanticTokensFull() error = %v", err)
	}
	if !slices.Equal(result.Data, []uint32{0, 0, 5, semanticTokenVariable, 0}) {
		t.Fatalf("semantic token data = %v", result.Data)
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

func decodeSemanticTokenData(data []uint32) []semanticToken {
	result := make([]semanticToken, 0, len(data)/5)
	var line uint32
	var start uint32
	for index := 0; index < len(data); index += 5 {
		line += data[index]
		if data[index] == 0 {
			start += data[index+1]
		} else {
			start = data[index+1]
		}
		result = append(result, semanticToken{
			line:      line,
			start:     start,
			length:    data[index+2],
			tokenType: data[index+3],
		})
	}
	return result
}
