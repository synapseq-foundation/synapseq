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

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
	"github.com/synapseq-foundation/synapseq/v4/internal/hub"
	s "github.com/synapseq-foundation/synapseq/v4/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// hubRunUpdate updates the local Hub manifest
func hubRunUpdate(quiet bool) error {
	if err := hub.HubUpdate(); err != nil {
		return fmt.Errorf("failed to update hub. Error\n  %v", err)
	}
	manifest, err := hub.GetManifest()
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
	if !hub.ManifestExists() {
		return fmt.Errorf("hub manifest not found. Please run 'synapseq -hub-update' to fetch the latest Hub manifest")
	}

	entry, err := hub.HubGet(sequenceId)
	if err != nil {
		return fmt.Errorf("failed to load hub manifest. Error\n  %v", err)
	}
	if entry == nil {
		return fmt.Errorf("sequence not found in hub: %s", sequenceId)
	}

	inputFile, err := hub.HubDownload(entry, t.HubActionTrackingGet)
	if err != nil {
		return fmt.Errorf("failed to download sequence from hub. Error\n  %v", err)
	}

	outputFormat := ".wav"
	if opts.Preview {
		outputFormat = ".html"
	}
	if opts.Mp3 {
		outputFormat = ".mp3"
	}
	if outputFile == "" {
		outputFile = entry.Name + outputFormat
	} else {
		outputFormat = strings.ToLower(filepath.Ext(outputFile))
	}

	appCtx := synapseq.NewAppContext()

	if !opts.Quiet && outputFile != "-" {
		appCtx = appCtx.WithVerbose(os.Stdout, !opts.NoColor)
	}

	loadedCtx, err := appCtx.Load(inputFile)
	if err != nil {
		return fmt.Errorf("failed to load sequence. Error\n  %v", err)
	}

	outputOpts := &outputOptions{
		OutputFile: outputFile,
		Quiet:      opts.Quiet,
		Preview:    opts.Preview,
		Play:       opts.Play,
		Mp3:        outputFormat == ".mp3" || opts.Mp3,
		FFplayPath: opts.FFplayPath,
		FFmpegPath: opts.FFmpegPath,
	}

	if err := processSequenceOutput(loadedCtx, outputOpts); err != nil {
		return fmt.Errorf("failed to process sequence output. Error\n  %v", err)
	}

	return nil
}

// hubRunList prints all available sequences from the Hub manifest in a tabular format
func hubRunList() error {
	if !hub.ManifestExists() {
		return fmt.Errorf("hub manifest not found. Please run 'synapseq -hub-update' to fetch the latest Hub manifest")
	}

	manifest, err := hub.GetManifest()
	if err != nil {
		return fmt.Errorf("failed to load hub manifest. Error\n  %v", err)
	}

	fmt.Printf("%s %s %s\n\n",
		cli.Title("SynapSeq Hub"),
		cli.Accent(fmt.Sprintf("%d available sequences", len(manifest.Entries))),
		cli.Muted(fmt.Sprintf("(Last updated: %s)", manifest.LastUpdated)))

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tAUTHOR\tCATEGORY\tUPDATED")

	for _, e := range manifest.Entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			e.ID,
			e.Author,
			e.Category,
			e.UpdatedAt[:10],
		)
	}

	w.Flush()
	printHubTable(buf.String())
	return nil
}

// hubRunSearch searches for sequences in the Hub by keyword (case-insensitive)
func hubRunSearch(query string) error {
	if !hub.ManifestExists() {
		return fmt.Errorf("hub manifest not found. Please run 'synapseq -hub-update' to fetch the latest Hub manifest")
	}

	if strings.TrimSpace(query) == "" {
		return fmt.Errorf("missing search term")
	}

	manifest, err := hub.GetManifest()
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

	fmt.Printf("%s %s %s\n\n", cli.Title("SynapSeq Hub"), cli.Accent(fmt.Sprintf("%d matching results", len(results))), cli.Muted(fmt.Sprintf("for %q", query)))

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tAUTHOR\tCATEGORY\tUPDATED")

	for _, e := range results {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			e.ID,
			e.Author,
			e.Category,
			e.UpdatedAt[:10],
		)
	}

	w.Flush()
	printHubTable(buf.String())
	return nil
}

// hubRunDownload downloads a sequence and all its dependencies into a given folder
func hubRunDownload(sequenceID, targetDir string, quiet bool) error {
	if !hub.ManifestExists() {
		return fmt.Errorf("hub manifest not found. Please run 'synapseq -hub-update' to fetch the latest Hub manifest")
	}

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

	manifest, err := hub.GetManifest()
	if err != nil {
		return fmt.Errorf("failed to load hub manifest. Error\n  %v", err)
	}

	var entry *t.HubEntry
	for _, e := range manifest.Entries {
		if e.ID == sequenceID {
			entry = &e
			break
		}
	}
	if entry == nil {
		return fmt.Errorf("sequence not found: %s", sequenceID)
	}

	seqFile, err := hub.HubDownload(entry, t.HubActionTrackingDownload)
	if err != nil {
		return fmt.Errorf("failed to download sequence from hub. Error\n  %v", err)
	}

	if err := s.CopyDir(filepath.Dir(seqFile), filepath.Join(targetDir, entry.Name)); err != nil {
		return fmt.Errorf("failed to copy files to target directory. Error\n  %v", err)
	}

	if !quiet {
		fmt.Printf("%s %s %s\n", cli.SuccessText("Downloaded"), cli.Accent(fmt.Sprintf("%q", entry.Name)), cli.Muted(fmt.Sprintf("and its dependencies to %s", targetDir)))
	}

	return nil
}

// hubRunInfo shows information about a sequence from the Hub
func hubRunInfo(sequenceID string) error {
	if !hub.ManifestExists() {
		return fmt.Errorf("hub manifest not found. Please run 'synapseq -hub-update' to fetch the latest Hub manifest")
	}

	if strings.TrimSpace(sequenceID) == "" {
		return fmt.Errorf("missing sequence ID")
	}

	manifest, err := hub.GetManifest()
	if err != nil {
		return fmt.Errorf("failed to load hub manifest. Error\n  %v", err)
	}

	var entry *t.HubEntry
	for _, e := range manifest.Entries {
		if e.ID == sequenceID {
			entry = &e
			break
		}
	}
	if entry == nil {
		return fmt.Errorf("sequence not found: %s", sequenceID)
	}

	seqFile, err := hub.HubDownload(entry, t.HubActionTrackingInfo)
	if err != nil {
		return fmt.Errorf("failed to download sequence from hub. Error\n  %v", err)
	}

	appCtx := synapseq.NewAppContext()

	loadedCtx, err := appCtx.Load(seqFile)
	if err != nil {
		return fmt.Errorf("failed to load sequence. Error\n  %v", err)
	}

	fmt.Printf("%s %s\n", cli.Label("Name:"), cli.Accent(entry.Name))
	fmt.Printf("%s %s\n", cli.Label("Author:"), cli.Accent(entry.Author))
	fmt.Printf("%s %s\n", cli.Label("Category:"), cli.Accent(entry.Category))
	fmt.Printf("%s %s\n", cli.Label("Updated At:"), cli.Accent(entry.UpdatedAt[:10]))

	dependencies := "\nDependencies: None\n"
	if len(entry.Dependencies) > 0 {
		dependencies = "\n" + cli.Section("Dependencies:") + "\n"
		for _, dep := range entry.Dependencies {
			dependencies += fmt.Sprintf("  - %s %s\n", cli.Accent(dep.Name), cli.Muted(fmt.Sprintf("(%s)", dep.Type.String())))
		}
	}
	fmt.Printf("%s", dependencies)

	description := "\nDescription: No description available.\n"
	comments := loadedCtx.Comments()
	if len(comments) > 0 {
		description = "\n" + cli.Section("Description:") + "\n"
		for _, comment := range comments {
			description += fmt.Sprintf("  %s\n", comment)
		}
	}
	fmt.Printf("%s", description)

	return nil
}

func printHubTable(table string) {
	lines := strings.Split(strings.TrimRight(table, "\n"), "\n")
	if len(lines) == 0 {
		return
	}
	fmt.Println(cli.Section(lines[0]))
	for _, line := range lines[1:] {
		fmt.Println(line)
	}
}
