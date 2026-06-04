// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"os"

	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
	template "github.com/synapseq-foundation/synapseq/v4/internal/template"
)

// generateTemplate generates a new sequence file from a template
func generateTemplate(templateName, outputFile string) error {
	content, err := template.GetTemplateContent(templateName)
	if err != nil {
		return err
	}

	if outputFile == "-" {
		fmt.Println(string(content))
		return nil
	}

	if err := os.WriteFile(outputFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write template to file: %v", err)
	}

	fmt.Printf("%s %s %s\n", cli.SuccessText("Template generated:"), cli.Accent(fmt.Sprintf("%q", templateName)), cli.Muted(fmt.Sprintf("as %q", outputFile)))
	fmt.Printf("%s %s\n", cli.Label("Run:"), cli.Command(fmt.Sprintf("synapseq %s", outputFile)))
	return nil
}
