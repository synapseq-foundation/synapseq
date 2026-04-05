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

package textstyle

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/synapseq-foundation/synapseq/v4/internal/palette"
)

// SetColorEnabled enables or disables ANSI colors for text output.
func SetColorEnabled(enabled bool) {
	color.NoColor = !enabled
}

func warmRGB(token palette.RGBColor, attrs ...color.Attribute) *color.Color {
	styled := color.RGB(token.R(), token.G(), token.B())
	if len(attrs) > 0 {
		styled.Add(attrs...)
	}
	return styled
}

// Title formats top-level headings.
func Title(text string) string {
	return warmRGB(palette.Terracotta, color.Bold).Sprint(text)
}

// Section formats section headings.
func Section(text string) string {
	return warmRGB(palette.Ochre, color.Bold).Sprint(text)
}

// Command formats shell commands and paths.
func Command(text string) string {
	return warmRGB(palette.Green).Sprint(text)
}

// Flag formats command-line flags.
func Flag(text string) string {
	return warmRGB(palette.Terracotta, color.Bold).Sprint(text)
}

// FlagColumn formats a left-aligned flag label.
func FlagColumn(text string, width int) string {
	return Flag(fmt.Sprintf("%-*s", width, text))
}

// Muted formats secondary explanatory text.
func Muted(text string) string {
	return warmRGB(palette.MutedWarm).Sprint(text)
}

// ErrorText formats error text.
func ErrorText(text string) string {
	return warmRGB(palette.DangerRed, color.Bold).Sprint(text)
}

// SuccessText formats success text.
func SuccessText(text string) string {
	return warmRGB(palette.Green, color.Bold).Sprint(text)
}

// Label formats field labels.
func Label(text string) string {
	return warmRGB(palette.TerracottaDark, color.Bold).Sprint(text)
}

// Accent formats highlighted values.
func Accent(text string) string {
	return warmRGB(palette.Terracotta).Sprint(text)
}