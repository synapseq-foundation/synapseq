// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
	"github.com/synapseq-foundation/synapseq/v4/internal/remote"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const remoteIndexMissingError = "remote index not found. Please run 'synapseq -sync' to fetch the latest Remote index"

// remoteRunSync updates the local Remote index.
func remoteRunSync(quiet bool) error {
	if err := remote.RemoteSync(); err != nil {
		return fmt.Errorf("failed to sync remote. Error\n  %v", err)
	}
	index, err := loadRemoteIndex()
	if err != nil {
		return fmt.Errorf("failed to get remote index. Error\n  %v", err)
	}
	if !quiet {
		fmt.Printf("%s %s %s\n", cli.SuccessText("Fetched"), cli.Accent(fmt.Sprintf("%d", len(index.Entries))), cli.Muted(fmt.Sprintf("entries from Remote. Last update: %s", index.LastUpdated)))
	}
	return nil
}

// remoteRunClean cleans the local Remote cache.
func remoteRunClean(quiet bool) error {
	if err := remote.RemoteClean(); err != nil {
		return fmt.Errorf("failed to clean remote cache. Error\n  %v", err)
	}
	if !quiet {
		fmt.Println(cli.SuccessText("Remote cache cleaned successfully."))
	}
	return nil
}

// remoteRunGet retrieves and processes a sequence from Remote.
func remoteRunGet(sequenceID, outputFile string, opts *cli.CLIOptions) error {
	if err := ensureRemoteIndexExists(); err != nil {
		return err
	}

	entry, err := remote.RemoteGet(sequenceID)
	if err != nil {
		return fmt.Errorf("failed to load remote index. Error\n  %v", err)
	}
	if entry == nil {
		return fmt.Errorf("sequence not found in remote: %s", sequenceID)
	}

	inputFile, err := remote.RemoteDownload(entry)
	if err != nil {
		return fmt.Errorf("failed to download sequence from remote. Error\n  %v", err)
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

// remoteRunList prints all available sequences from the Remote index in a tabular format.
func remoteRunList() error {
	index, err := loadRemoteIndex()
	if err != nil {
		return fmt.Errorf("failed to load remote index. Error\n  %v", err)
	}

	entries := sortedRemoteEntries(index.Entries)
	printRemoteHeading(fmt.Sprintf("%d available sequences", len(entries)), fmt.Sprintf("(Last updated: %s)", index.LastUpdated))
	printRemoteEntriesTable(entries)
	return nil
}

// remoteRunSearch searches for sequences in Remote by keyword (case-insensitive).
func remoteRunSearch(query string) error {
	if strings.TrimSpace(query) == "" {
		return fmt.Errorf("missing search term")
	}

	index, err := loadRemoteIndex()
	if err != nil {
		return fmt.Errorf("failed to load remote index. Error\n  %v", err)
	}

	query = strings.ToLower(query)
	var results []t.RemoteEntry

	for _, e := range index.Entries {
		if strings.Contains(strings.ToLower(e.Name), query) ||
			strings.Contains(strings.ToLower(e.Description), query) ||
			strings.Contains(strings.ToLower(e.Category), query) {
			results = append(results, e)
		}
	}

	if len(results) == 0 {
		fmt.Printf("%s %s\n", cli.Muted("No matches found for"), cli.Accent(fmt.Sprintf("%q", query)))
		return nil
	}

	results = sortedRemoteEntries(results)
	printRemoteHeading(fmt.Sprintf("%d matching results", len(results)), fmt.Sprintf("for %q", query))
	printRemoteEntriesTable(results)
	return nil
}

// remoteRunDownload downloads a sequence into a given folder.
func remoteRunDownload(sequenceID, targetDir string, quiet bool) error {
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

	index, err := loadRemoteIndex()
	if err != nil {
		return fmt.Errorf("failed to load remote index. Error\n  %v", err)
	}

	entry := findRemoteEntryByID(index, sequenceID)
	if entry == nil {
		return fmt.Errorf("sequence not found: %s", sequenceID)
	}

	seqFile, err := remote.RemoteDownload(entry)
	if err != nil {
		return fmt.Errorf("failed to download sequence from remote. Error\n  %v", err)
	}

	if err := copyRemoteSequence(seqFile, filepath.Join(targetDir, entry.ID+".spsq")); err != nil {
		return fmt.Errorf("failed to copy file to target directory. Error\n  %v", err)
	}

	if !quiet {
		fmt.Printf("%s %s %s\n", cli.SuccessText("Downloaded"), cli.Accent(fmt.Sprintf("%q", entry.ID)), cli.Muted(fmt.Sprintf("to %s", targetDir)))
	}

	return nil
}

// remoteRunInfo shows information about a sequence from Remote.
func remoteRunInfo(sequenceID string) error {
	if strings.TrimSpace(sequenceID) == "" {
		return fmt.Errorf("missing sequence ID")
	}

	index, err := loadRemoteIndex()
	if err != nil {
		return fmt.Errorf("failed to load remote index. Error\n  %v", err)
	}

	entry := findRemoteEntryByID(index, sequenceID)
	if entry == nil {
		return fmt.Errorf("sequence not found: %s", sequenceID)
	}

	printRemoteInfoSummary(entry)
	fmt.Print(formatRemoteDescription(entry.Description))

	return nil
}

func printRemoteHeading(title, detail string) {
	fmt.Printf("%s %s %s\n\n", cli.Title("SynapSeq Remote"), cli.Accent(title), cli.Muted(detail))
}

func printRemoteEntriesTable(entries []t.RemoteEntry) {
	fmt.Print(renderRemoteEntriesTable(entries))
}

func renderRemoteEntriesTable(entries []t.RemoteEntry) string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tDURATION\tCATEGORY\tCREATED")

	for _, entry := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			entry.ID,
			fmt.Sprintf("%d min", entry.DurationMinutes),
			entry.Category,
			shortDate(entry.CreatedAt),
		)
	}

	w.Flush()
	return formatRemoteTable(buf.String())
}

func formatRemoteTable(table string) string {
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

func printRemoteInfoSummary(entry *t.RemoteEntry) {
	fmt.Printf("%s %s\n", cli.Label("Name:"), cli.Accent(entry.Name))
	fmt.Printf("%s %s\n", cli.Label("Category:"), cli.Accent(entry.Category))
	fmt.Printf("%s %s\n", cli.Label("Created At:"), cli.Accent(shortDate(entry.CreatedAt)))
}

func formatRemoteDescription(description string) string {
	if strings.TrimSpace(description) == "" {
		return "\nDescription: No description available.\n"
	}

	var builder strings.Builder
	builder.WriteString("\n")
	builder.WriteString(cli.Section("Description:"))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("  %s\n", description))

	return builder.String()
}

func ensureRemoteIndexExists() error {
	if !remote.IndexExists() {
		return fmt.Errorf(remoteIndexMissingError)
	}

	return nil
}

func loadRemoteIndex() (*t.RemoteIndex, error) {
	if err := ensureRemoteIndexExists(); err != nil {
		return nil, err
	}

	return remote.GetIndex()
}

func findRemoteEntryByID(index *t.RemoteIndex, sequenceID string) *t.RemoteEntry {
	if index == nil {
		return nil
	}

	for _, entry := range index.Entries {
		if entry.ID == sequenceID {
			match := entry
			return &match
		}
	}

	return nil
}

func sortedRemoteEntries(entries []t.RemoteEntry) []t.RemoteEntry {
	sorted := append([]t.RemoteEntry(nil), entries...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].CreatedAt > sorted[j].CreatedAt
	})
	return sorted
}

func shortDate(value string) string {
	if len(value) <= 10 {
		return value
	}

	return value[:10]
}

func copyRemoteSequence(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0644)
}
