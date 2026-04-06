//go:build !js && !wasm

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

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
	"github.com/synapseq-foundation/synapseq/v4/internal/hub"
	r "github.com/synapseq-foundation/synapseq/v4/internal/resource"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const hubManifestMissingError = "hub manifest not found. Please run 'synapseq -hub-update' to fetch the latest Hub manifest"

// hubRunUpdate updates the local Hub manifest
func hubRunUpdate(quiet bool) error {
	if err := hub.HubUpdate(); err != nil {
		return fmt.Errorf("failed to update hub. Error\n  %v", err)
	}
	manifest, err := loadHubManifest()
	if err != nil {
		return fmt.Errorf("failed to get hub manifest. Error\n  %v", err)
	}
	if !quiet {
		fmt.Printf("%s %s %s\n", cli.SuccessText("Fetched"), cli.Accent(fmt.Sprintf("%d", len(manifest.Entries))), cli.Muted(fmt.Sprintf("entries from the Hub. Last update: %s", manifest.LastUpdated)))
	}
	return nil
}

// hubRunClean cleans the local Hub cache
func hubRunClean(quiet bool) error {
	if err := hub.HubClean(); err != nil {
		return fmt.Errorf("failed to clean hub cache. Error\n  %v", err)
	}
	if !quiet {
		fmt.Println(cli.SuccessText("Hub cache cleaned successfully."))
	}
	return nil
}

// hubRunGet retrieves and processes a sequence from the Hub
func hubRunGet(sequenceId, outputFile string, opts *cli.CLIOptions) error {
	if err := ensureHubManifestExists(); err != nil {
		return err
	}

	entry, err := hub.HubGet(sequenceId)
	if err != nil {
		return fmt.Errorf("failed to load hub manifest. Error\n  %v", err)
	}
	if entry == nil {
		return fmt.Errorf("sequence not found in hub: %s", sequenceId)
	}

	inputFile, err := hub.HubDownload(entry)
	if err != nil {
		return fmt.Errorf("failed to download sequence from hub. Error\n  %v", err)
	}

	outputFile, outputFormat := resolveOutputTarget(entry.Name, outputFile, opts)

	loadedCtx, err := loadSequenceContext(inputFile, outputFile, os.Stdout, opts)
	if err != nil {
		return fmt.Errorf("failed to load sequence. Error\n  %v", err)
	}

	if err := runLoadedSequence(loadedCtx, outputFile, outputFormat, opts); err != nil {
		return fmt.Errorf("failed to process sequence output. Error\n  %v", err)
	}

	return nil
}

// hubRunList prints all available sequences from the Hub manifest in a tabular format
func hubRunList() error {
	manifest, err := loadHubManifest()
	if err != nil {
		return fmt.Errorf("failed to load hub manifest. Error\n  %v", err)
	}

	printHubHeading(fmt.Sprintf("%d available sequences", len(manifest.Entries)), fmt.Sprintf("(Last updated: %s)", manifest.LastUpdated))
	printHubEntriesTable(manifest.Entries)
	return nil
}

// hubRunSearch searches for sequences in the Hub by keyword (case-insensitive)
func hubRunSearch(query string) error {
	if strings.TrimSpace(query) == "" {
		return fmt.Errorf("missing search term")
	}

	manifest, err := loadHubManifest()
	if err != nil {
		return fmt.Errorf("failed to load hub manifest. Error\n  %v", err)
	}

	query = strings.ToLower(query)
	var results []t.HubEntry

	for _, e := range manifest.Entries {
		if strings.Contains(strings.ToLower(e.ID), query) ||
			strings.Contains(strings.ToLower(e.Name), query) ||
			strings.Contains(strings.ToLower(e.Author), query) ||
			strings.Contains(strings.ToLower(e.Category), query) {
			results = append(results, e)
		}
	}

	if len(results) == 0 {
		fmt.Printf("%s %s\n", cli.Muted("No matches found for"), cli.Accent(fmt.Sprintf("%q", query)))
		return nil
	}

	printHubHeading(fmt.Sprintf("%d matching results", len(results)), fmt.Sprintf("for %q", query))
	printHubEntriesTable(results)
	return nil
}

// hubRunDownload downloads a sequence and all its dependencies into a given folder
func hubRunDownload(sequenceID, targetDir string, quiet bool) error {
	if strings.TrimSpace(sequenceID) == "" {
		return fmt.Errorf("missing sequence ID")
	}

	if targetDir == "" || targetDir == "." {
		var err error
		targetDir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	manifest, err := loadHubManifest()
	if err != nil {
		return fmt.Errorf("failed to load hub manifest. Error\n  %v", err)
	}

	entry := findHubEntryByID(manifest, sequenceID)
	if entry == nil {
		return fmt.Errorf("sequence not found: %s", sequenceID)
	}

	seqFile, err := hub.HubDownload(entry)
	if err != nil {
		return fmt.Errorf("failed to download sequence from hub. Error\n  %v", err)
	}

	if err := r.CopyDir(filepath.Dir(seqFile), filepath.Join(targetDir, entry.ID)); err != nil {
		return fmt.Errorf("failed to copy files to target directory. Error\n  %v", err)
	}

	if !quiet {
		fmt.Printf("%s %s %s\n", cli.SuccessText("Downloaded"), cli.Accent(fmt.Sprintf("%q", entry.ID)), cli.Muted(fmt.Sprintf("and its dependencies to %s", targetDir)))
	}

	return nil
}

// hubRunInfo shows information about a sequence from the Hub
func hubRunInfo(sequenceID string) error {
	if strings.TrimSpace(sequenceID) == "" {
		return fmt.Errorf("missing sequence ID")
	}

	manifest, err := loadHubManifest()
	if err != nil {
		return fmt.Errorf("failed to load hub manifest. Error\n  %v", err)
	}

	entry := findHubEntryByID(manifest, sequenceID)
	if entry == nil {
		return fmt.Errorf("sequence not found: %s", sequenceID)
	}

	seqFile, err := hub.HubDownload(entry)
	if err != nil {
		return fmt.Errorf("failed to download sequence from hub. Error\n  %v", err)
	}

	loadedCtx, err := loadSequenceContext(seqFile, "", nil, &cli.CLIOptions{Quiet: true})
	if err != nil {
		return fmt.Errorf("failed to load sequence. Error\n  %v", err)
	}

	printHubInfoSummary(entry)
	fmt.Print(formatHubDependencies(entry.Dependencies))
	fmt.Print(formatHubDescription(loadedCtx.Comments()))

	return nil
}

func printHubHeading(title, detail string) {
	fmt.Printf("%s %s %s\n\n", cli.Title("SynapSeq Hub"), cli.Accent(title), cli.Muted(detail))
}

func printHubEntriesTable(entries []t.HubEntry) {
	fmt.Print(renderHubEntriesTable(entries))
}

func renderHubEntriesTable(entries []t.HubEntry) string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tAUTHOR\tCATEGORY\tUPDATED")

	for _, entry := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			entry.ID,
			entry.Author,
			entry.Category,
			entry.UpdatedAt[:10],
		)
	}

	w.Flush()
	return formatHubTable(buf.String())
}

func formatHubTable(table string) string {
	lines := strings.Split(strings.TrimRight(table, "\n"), "\n")
	if len(lines) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString(cli.Section(lines[0]))
	builder.WriteString("\n")
	for _, line := range lines[1:] {
		builder.WriteString(line)
		builder.WriteString("\n")
	}

	return builder.String()
}

func printHubInfoSummary(entry *t.HubEntry) {
	fmt.Printf("%s %s\n", cli.Label("Name:"), cli.Accent(entry.Name))
	fmt.Printf("%s %s\n", cli.Label("Author:"), cli.Accent(entry.Author))
	fmt.Printf("%s %s\n", cli.Label("Category:"), cli.Accent(entry.Category))
	fmt.Printf("%s %s\n", cli.Label("Updated At:"), cli.Accent(entry.UpdatedAt[:10]))
}

func formatHubDependencies(dependencies []t.HubDependency) string {
	if len(dependencies) == 0 {
		return "\nDependencies: None\n"
	}

	var builder strings.Builder
	builder.WriteString("\n")
	builder.WriteString(cli.Section("Dependencies:"))
	builder.WriteString("\n")
	for _, dependency := range dependencies {
		builder.WriteString(fmt.Sprintf("  - %s %s\n", cli.Accent(dependency.ID), cli.Muted(fmt.Sprintf("(%s)", dependency.Type.String()))))
	}

	return builder.String()
}

func formatHubDescription(comments []string) string {
	if len(comments) == 0 {
		return "\nDescription: No description available.\n"
	}

	var builder strings.Builder
	builder.WriteString("\n")
	builder.WriteString(cli.Section("Description:"))
	builder.WriteString("\n")
	for _, comment := range comments {
		builder.WriteString(fmt.Sprintf("  %s\n", comment))
	}

	return builder.String()
}

func ensureHubManifestExists() error {
	if !hub.ManifestExists() {
		return fmt.Errorf(hubManifestMissingError)
	}

	return nil
}

func loadHubManifest() (*t.HubManifest, error) {
	if err := ensureHubManifestExists(); err != nil {
		return nil, err
	}

	return hub.GetManifest()
}

func findHubEntryByID(manifest *t.HubManifest, sequenceID string) *t.HubEntry {
	if manifest == nil {
		return nil
	}

	for _, entry := range manifest.Entries {
		if entry.ID == sequenceID {
			match := entry
			return &match
		}
	}

	return nil
}
