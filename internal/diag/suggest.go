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

package diag

import "strings"

// DefaultSuggestionDistance returns a conservative max edit distance for typo hints.
func DefaultSuggestionDistance(input string) int {
	switch n := len(input); {
	case n <= 4:
		return 1
	case n <= 8:
		return 2
	default:
		return 3
	}
}

// ClosestMatch returns the best candidate within maxDistance.
func ClosestMatch(input string, candidates []string, maxDistance int) (string, bool) {
	if input == "" || len(candidates) == 0 || maxDistance < 0 {
		return "", false
	}

	input = strings.ToLower(input)
	best := ""
	bestDistance := maxDistance + 1
	tied := false

	for _, candidate := range candidates {
		distance := levenshtein(input, strings.ToLower(candidate))
		if distance > maxDistance {
			continue
		}
		if distance < bestDistance {
			best = candidate
			bestDistance = distance
			tied = false
			continue
		}
		if distance == bestDistance {
			tied = true
		}
	}

	if best == "" || tied {
		return "", false
	}

	return best, true
}

func levenshtein(a, b string) int {
	if a == b {
		return 0
	}
	if a == "" {
		return len(b)
	}
	if b == "" {
		return len(a)
	}

	prev := make([]int, len(b)+1)
	current := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}

	for i := 1; i <= len(a); i++ {
		current[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}

			deletion := prev[j] + 1
			insertion := current[j-1] + 1
			substitution := prev[j-1] + cost

			current[j] = min(deletion, insertion, substitution)
		}
		prev, current = current, prev
	}

	return prev[len(b)]
}

func min(values ...int) int {
	best := values[0]
	for _, value := range values[1:] {
		if value < best {
			best = value
		}
	}
	return best
}
