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
	"path/filepath"
	"strings"
	"sync"
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
		fmt.Printf("Fetched %d entries from the Hub. Last update: %s\n", len(manifest.Entries), manifest.LastUpdated)
	}
	return nil
}

// hubRunClean cleans the local Hub cache
func hubRunClean(quiet bool) error {
	if err := hub.HubClean(); err != nil {
		return fmt.Errorf("failed to clean hub cache. Error\n  %v", err)
	}
	if !quiet {
		fmt.Println("Hub cache cleaned successfully.")
	}
	return nil
}

// hubRunGet retrieves and processes a sequence from the Hub
func hubRunGet(sequenceId, outputFile string, opts *cli.CLIOptions) error {
	var wg sync.WaitGroup

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

	inputFile, err := hub.HubDownload(entry, t.HubActionTrackingGet, &wg)
	if err != nil {
		return fmt.Errorf("failed to download sequence from hub. Error\n  %v", err)
	}

	if outputFile == "" {
		if opts.Mp3 {
			outputFile = entry.Name + ".mp3"
		} else {
			outputFile = entry.Name + ".wav"
		}
	}

	appCtx, err := synapseq.NewAppContext(inputFile, outputFile)
	if err != nil {
		return fmt.Errorf("failed to create application context. Error\n  %v", err)
	}

	if !opts.Quiet && outputFile != "-" {
		appCtx = appCtx.WithVerbose(os.Stdout)
	}

	if err := appCtx.LoadSequence(); err != nil {
		return fmt.Errorf("failed to load sequence. Error\n  %v", err)
	}

	outputOpts := &outputOptions{
		OutputFile: outputFile,
		Quiet:      opts.Quiet,
		Play:       opts.Play,
		Mp3:        opts.Mp3,
		FFplayPath: opts.FFplayPath,
		FFmpegPath: opts.FFmpegPath,
	}

	if err := processSequenceOutput(appCtx, outputOpts); err != nil {
		return fmt.Errorf("failed to process sequence output. Error\n  %v", err)
	}

	wg.Wait()
	return nil
}

// / hubRunList prints all available sequences from the Hub manifest in a tabular format
func hubRunList() error {
	if !hub.ManifestExists() {
		return fmt.Errorf("hub manifest not found. Please run 'synapseq -hub-update' to fetch the latest Hub manifest")
	}

	manifest, err := hub.GetManifest()
	if err != nil {
		return fmt.Errorf("failed to load hub manifest. Error\n  %v", err)
	}

	fmt.Printf("SynapSeq Hub — %d available sequences  (Last updated: %s)\n\n",
		len(manifest.Entries), manifest.LastUpdated)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
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
		fmt.Printf("No matches found for %q\n", query)
		return nil
	}

	fmt.Printf("SynapSeq Hub - %d matching results for %q\n\n", len(results), query)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
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
	return nil
}

// hubRunDownload downloads a sequence and all its dependencies into a given folder
func hubRunDownload(sequenceID, targetDir string, quiet bool) error {
	var wg sync.WaitGroup

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

	seqFile, err := hub.HubDownload(entry, t.HubActionTrackingDownload, &wg)
	if err != nil {
		return fmt.Errorf("failed to download sequence from hub. Error\n  %v", err)
	}

	if err := s.CopyDir(filepath.Dir(seqFile), filepath.Join(targetDir, entry.Name)); err != nil {
		return fmt.Errorf("failed to copy files to target directory. Error\n  %v", err)
	}

	wg.Wait()
	if !quiet {
		fmt.Printf("Sequence %q and its dependencies have been downloaded to %s\n", entry.Name, targetDir)
	}

	return nil
}

// hubRunInfo shows information about a sequence from the Hub
func hubRunInfo(sequenceID string) error {
	var wg sync.WaitGroup

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

	seqFile, err := hub.HubDownload(entry, t.HubActionTrackingInfo, &wg)
	if err != nil {
		return fmt.Errorf("failed to download sequence from hub. Error\n  %v", err)
	}

	appCtx, err := synapseq.NewAppContext(seqFile, "")
	if err != nil {
		return fmt.Errorf("failed to create application context. Error\n  %v", err)
	}

	if err := appCtx.LoadSequence(); err != nil {
		return fmt.Errorf("failed to load sequence. Error\n  %v", err)
	}

	fmt.Printf("Name:        %s\n", entry.Name)
	fmt.Printf("Author:      %s\n", entry.Author)
	fmt.Printf("Category:    %s\n", entry.Category)
	fmt.Printf("Updated At:  %s\n", entry.UpdatedAt[:10])

	dependencies := "\nDependencies: None\n"
	if len(entry.Dependencies) > 0 {
		dependencies = "\nDependencies:\n"
		for _, dep := range entry.Dependencies {
			dependencies += fmt.Sprintf("  - %s (%s)\n", dep.Name, dep.Type.String())
		}
	}
	fmt.Printf("%s", dependencies)

	description := "\nDescription: No description available.\n"
	comments := appCtx.Comments()
	if len(comments) > 0 {
		description = "\nDescription:\n"
		for _, comment := range comments {
			description += fmt.Sprintf("  %s\n", comment)
		}
	}
	fmt.Printf("%s", description)

	wg.Wait()
	return nil
}
