// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

// Package lsp implements the SynapSeq language server.
package lsp

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	"github.com/synapseq-foundation/synapseq/v4/internal/parser"
	"github.com/synapseq-foundation/synapseq/v4/internal/sequence"
	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

const serverName = "synapseq"

type stdio struct {
	in  io.Reader
	out io.Writer
}

func (s stdio) Read(p []byte) (int, error) {
	return s.in.Read(p)
}

func (s stdio) Write(p []byte) (int, error) {
	return s.out.Write(p)
}

func (stdio) Close() error {
	return nil
}

type document struct {
	uri     uri.URI
	version int32
	text    string
}

type server struct {
	protocol.UnimplementedServer

	mu        sync.RWMutex
	documents map[uri.URI]document
	client    protocol.Client
	stop      context.CancelFunc
}

// Run serves LSP traffic over the supplied standard input and output streams.
func Run(ctx context.Context, input io.Reader, output io.Writer) error {
	ctx, stop := context.WithCancel(ctx)
	defer stop()

	server := &server{
		documents: map[uri.URI]document{},
		stop:      stop,
	}
	_, connection, client := protocol.NewServer(ctx, server, jsonrpc2.NewStream(stdio{
		in:  input,
		out: output,
	}))
	server.client = client

	<-connection.Done()
	return connection.Err()
}

func (s *server) Initialize(context.Context, *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			PositionEncoding: protocol.PositionEncodingKindUTF16,
			TextDocumentSync: protocol.TextDocumentSyncKindFull,
			CompletionProvider: &protocol.CompletionOptions{
				TriggerCharacters: []string{"@"},
			},
		},
		ServerInfo: protocol.ServerInfo{
			Name: serverName,
		},
	}, nil
}

func (s *server) Shutdown(context.Context) error {
	return nil
}

func (s *server) Exit(context.Context) error {
	s.stop()
	return nil
}

func (s *server) DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error {
	doc := document{
		uri:     params.TextDocument.URI,
		version: params.TextDocument.Version,
		text:    params.TextDocument.Text,
	}
	s.store(doc)
	return s.publishDiagnostics(ctx, doc)
}

func (s *server) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	if len(params.ContentChanges) == 0 {
		return nil
	}

	change, ok := params.ContentChanges[len(params.ContentChanges)-1].(*protocol.TextDocumentContentChangeWholeDocument)
	if !ok {
		return nil
	}
	doc := document{
		uri:     params.TextDocument.URI,
		version: params.TextDocument.Version,
		text:    change.Text,
	}
	s.store(doc)
	return s.publishDiagnostics(ctx, doc)
}

func (s *server) DidSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) error {
	doc, ok := s.load(params.TextDocument.URI)
	if !ok {
		return nil
	}
	if params.Text != nil {
		doc.text = *params.Text
		s.store(doc)
	}
	return s.publishDiagnostics(ctx, doc)
}

func (s *server) DidClose(ctx context.Context, params *protocol.DidCloseTextDocumentParams) error {
	s.mu.Lock()
	delete(s.documents, params.TextDocument.URI)
	s.mu.Unlock()

	return s.client.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: []protocol.Diagnostic{},
	})
}

func (s *server) Completion(_ context.Context, params *protocol.CompletionParams) (protocol.CompletionResult, error) {
	doc, ok := s.load(params.TextDocument.URI)
	if !ok {
		return protocol.CompletionItemSlice{}, nil
	}
	return complete(doc.text, params.Position), nil
}

func (s *server) store(doc document) {
	s.mu.Lock()
	s.documents[doc.uri] = doc
	s.mu.Unlock()
}

func (s *server) load(documentURI uri.URI) (document, bool) {
	s.mu.RLock()
	doc, ok := s.documents[documentURI]
	s.mu.RUnlock()
	return doc, ok
}

func (s *server) publishDiagnostics(ctx context.Context, doc document) error {
	return s.client.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{
		URI:         doc.uri,
		Version:     protocol.NewOptional(doc.version),
		Diagnostics: validateDocument(doc),
	})
}

func validateDocument(doc document) []protocol.Diagnostic {
	if doc.uri.IsFile() && !hasRemoteExtends(doc.text) {
		path := doc.uri.Path()
		_, err := sequence.LoadTextSequence([]byte(doc.text), path, filepath.Dir(path))
		return diagnosticsFromError(err, doc.text)
	}

	return syntaxDiagnostics(doc.text)
}

func hasRemoteExtends(text string) bool {
	for _, line := range strings.Split(text, "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 || fields[0] != "@extends" {
			continue
		}
		if strings.HasPrefix(fields[1], "http://") || strings.HasPrefix(fields[1], "https://") {
			return true
		}
	}
	return false
}

func syntaxDiagnostics(text string) []protocol.Diagnostic {
	result := []protocol.Diagnostic{}
	lines := strings.Split(text, "\n")
	for index, line := range lines {
		ctx := parser.NewTextParser(strings.TrimSuffix(line, "\r"))
		if len(ctx.Line.Tokens) == 0 || ctx.HasComment() {
			continue
		}

		var err error
		switch {
		case ctx.HasOption():
			_, err = ctx.ParseOption("")
		case ctx.HasPreset():
			_, err = ctx.ParsePresetDeclaration()
		case ctx.HasTrack():
			_, err = ctx.ParseTrackDeclaration()
		case ctx.HasTrackOverride():
			_, err = ctx.ParseTrackOverrideDeclaration()
		case ctx.HasTimeline():
			_, err = ctx.ParseTimelineDeclaration()
		default:
			err = diag.Parse("invalid syntax").WithSpan(diag.Span{
				Column:    1,
				EndColumn: 2,
				LineText:  line,
			})
		}
		if err == nil {
			continue
		}

		diagnostics := diagnosticsFromError(err, line)
		for item := range diagnostics {
			diagnostics[item].Range.Start.Line = uint32(index)
			diagnostics[item].Range.End.Line = uint32(index)
		}
		result = append(result, diagnostics...)
	}
	return result
}

func diagnosticsFromError(err error, text string) []protocol.Diagnostic {
	if err == nil {
		return []protocol.Diagnostic{}
	}

	diagnostic, ok := diag.As(err)
	if !ok {
		diagnostic = diag.New(diag.KindInternal, err.Error())
	}
	return []protocol.Diagnostic{toProtocolDiagnostic(diagnostic, text)}
}

func toProtocolDiagnostic(diagnostic *diag.Diagnostic, text string) protocol.Diagnostic {
	line := diagnostic.Span.Line - 1
	if line < 0 {
		line = 0
	}
	lines := strings.Split(text, "\n")
	lineText := diagnostic.Span.LineText
	if line < len(lines) {
		lineText = strings.TrimSuffix(lines[line], "\r")
	}

	start := utf16Offset(lineText, diagnostic.Span.Column-1)
	end := utf16Offset(lineText, diagnostic.Span.EndColumn-1)
	if end <= start {
		end = start + 1
	}

	return protocol.Diagnostic{
		Range: protocol.Range{
			Start: protocol.Position{Line: uint32(line), Character: uint32(start)},
			End:   protocol.Position{Line: uint32(line), Character: uint32(end)},
		},
		Severity: protocol.DiagnosticSeverityError,
		Source:   protocol.NewOptional(serverName),
		Message:  protocol.String(diagnosticMessage(diagnostic)),
	}
}

func diagnosticMessage(diagnostic *diag.Diagnostic) string {
	parts := []string{diagnostic.Message}
	if diagnostic.Found != "" {
		parts = append(parts, fmt.Sprintf("found %q", diagnostic.Found))
	}
	if len(diagnostic.Expected) > 0 {
		parts = append(parts, "expected "+strings.Join(diagnostic.Expected, ", "))
	}
	if diagnostic.Suggestion != "" {
		parts = append(parts, diagnostic.Suggestion)
	}
	if diagnostic.Hint != "" {
		parts = append(parts, diagnostic.Hint)
	}
	return strings.Join(parts, "; ")
}

func utf16Offset(text string, byteOffset int) int {
	if byteOffset <= 0 {
		return 0
	}
	if byteOffset > len(text) {
		byteOffset = len(text)
	}

	units := 0
	for index := 0; index < byteOffset; {
		runeValue, size := utf8.DecodeRuneInString(text[index:])
		if size == 0 {
			break
		}
		if index+size > byteOffset {
			break
		}
		if runeValue > 0xFFFF {
			units += 2
		} else {
			units++
		}
		index += size
	}
	return units
}
