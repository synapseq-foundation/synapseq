/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

package main

import (
	"fmt"
	"os"

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

	fmt.Printf("Template %q has been generated as %q\n", templateName, outputFile)
	fmt.Printf("You can edit the file and run it with: synapseq %s\n", outputFile)
	return nil
}
