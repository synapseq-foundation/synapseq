/*
 * SynapSeq - Text-Driven Audio Sequencer for Brainwave Entrainment
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

package manual

import (
	"strings"

	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
)

func writeTitle(b *strings.Builder, title, subtitle string) {
	b.WriteString(cli.Title(title))
	b.WriteString("\n")
	b.WriteString(cli.Muted(subtitle))
	b.WriteString("\n\n")
}

func writeSection(b *strings.Builder, title string) {
	b.WriteString(cli.Section(strings.ToUpper(title)))
	b.WriteString("\n")
}

func writeSubsection(b *strings.Builder, title string) {
	writeIndentedSubsection(b, 4, title)
}

func writeParagraph(b *strings.Builder, lines ...string) {
	for _, line := range lines {
		if line == "" {
			continue
		}
		writeWrappedLine(b, line, "    ", "    ")
	}
	b.WriteString("\n")
}

func writeBullet(b *strings.Builder, label, description string) {
	b.WriteString("    ")
	b.WriteString(cli.Label(label))
	b.WriteString("\n")
	for _, line := range strings.Split(description, "\n") {
		if line == "" {
			continue
		}
		writeWrappedLine(b, line, "        ", "        ")
	}
	b.WriteString("\n")
}

func writeLineBlock(b *strings.Builder, lines ...string) {
	writeIndentedLineBlock(b, 8, lines...)
}

func writeNestedSubsection(b *strings.Builder, title string) {
	writeIndentedSubsection(b, 8, title)
}

func writeNestedBullet(b *strings.Builder, label, description string) {
	writeIndentedBullet(b, 8, 12, label, description)
}

func writeDeepBullet(b *strings.Builder, label, description string) {
	writeIndentedBullet(b, 12, 16, label, description)
}

func writeNestedLineBlock(b *strings.Builder, lines ...string) {
	writeIndentedLineBlock(b, 12, lines...)
}

func writeNestedCodeBlock(b *strings.Builder, lines ...string) {
	writeIndentedCodeBlock(b, 12, lines...)
}

func writeIndentedSubsection(b *strings.Builder, indent int, title string) {
	b.WriteString(strings.Repeat(" ", indent))
	b.WriteString(cli.Label(title))
	b.WriteString("\n\n")
}

func writeIndentedBullet(b *strings.Builder, labelIndent, textIndent int, label, description string) {
	b.WriteString(strings.Repeat(" ", labelIndent))
	b.WriteString(cli.Label(label))
	b.WriteString("\n")
	prefix := strings.Repeat(" ", textIndent)
	for _, line := range strings.Split(description, "\n") {
		if line == "" {
			continue
		}
		writeWrappedLine(b, line, prefix, prefix)
	}
	b.WriteString("\n")
}

func writeIndentedLineBlock(b *strings.Builder, indent int, lines ...string) {
	prefix := strings.Repeat(" ", indent)
	for _, line := range lines {
		if line == "" {
			continue
		}
		writeWrappedLine(b, line, prefix, prefix)
	}
	b.WriteString("\n")
}

func writeCodeBlock(b *strings.Builder, lines ...string) {
	writeIndentedCodeBlock(b, 8, lines...)
}

func writeExample(b *strings.Builder, example string) {
	writeIndentedSubsection(b, 8, "Example")
	writeIndentedCodeBlock(b, 12, example)
}

func writeNestedExample(b *strings.Builder, example string) {
	writeIndentedSubsection(b, 12, "Example")
	writeIndentedCodeBlock(b, 16, example)
}

func writeIndentedCodeBlock(b *strings.Builder, indent int, lines ...string) {
	prefix := strings.Repeat(" ", indent)
	for _, line := range lines {
		if line == "" {
			b.WriteString("\n")
			continue
		}
		b.WriteString(prefix)
		b.WriteString(cli.Command(line))
		b.WriteString("\n")
	}
	b.WriteString("\n")
}

func writeWrappedLine(b *strings.Builder, text, firstPrefix, continuationPrefix string) {
	writeWrappedLineWithPrefixes(b, text, firstPrefix, continuationPrefix, len(firstPrefix), len(continuationPrefix))
}

func writeWrappedLineWithPrefixes(b *strings.Builder, text, firstPrefix, continuationPrefix string, firstWidth, continuationWidth int) {
	words := strings.Fields(text)
	if len(words) == 0 {
		b.WriteString("\n")
		return
	}

	prefix := firstPrefix
	available := manualWidth - firstWidth
	lineLen := 0

	b.WriteString(prefix)
	for index, word := range words {
		wordLen := len(word)
		separatorLen := 0
		if lineLen > 0 {
			separatorLen = 1
		}

		if lineLen > 0 && lineLen+separatorLen+wordLen > available {
			b.WriteString("\n")
			prefix = continuationPrefix
			available = manualWidth - continuationWidth
			b.WriteString(prefix)
			b.WriteString(word)
			lineLen = wordLen
			continue
		}

		if index > 0 && lineLen > 0 {
			b.WriteString(" ")
		}
		b.WriteString(word)
		lineLen += separatorLen + wordLen
	}
	if lineLen == 0 {
		b.WriteString(strings.Join(words, " "))
	}
	b.WriteString("\n")
}