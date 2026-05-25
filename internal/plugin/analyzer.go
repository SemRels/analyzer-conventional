// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The analyzer-conventional Authors

package plugin

import (
	"fmt"
	"regexp"
	"strings"
)

type BumpLevel string

const (
	BumpNone  BumpLevel = "none"
	BumpPatch BumpLevel = "patch"
	BumpMinor BumpLevel = "minor"
	BumpMajor BumpLevel = "major"
)

type AnalysisResult struct {
	Bump   BumpLevel `json:"bump"`
	Reason string    `json:"reason"`
}

type Analyzer struct{}

var conventionalHeaderPattern = regexp.MustCompile(`^(\w+)(\([\w-]+\))?(!)?:(.+)$`)

func New() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) Analyze(commits []string) AnalysisResult {
	result := AnalysisResult{
		Bump:   BumpNone,
		Reason: "no conventional commits require a version bump",
	}

	if len(commits) == 0 {
		return result
	}

	majorCount := 0
	minorCount := 0
	patchCount := 0
	patchTypes := map[string]int{}

	for _, commit := range commits {
		commitType, breaking := parseCommit(commit)
		if commitType == "" && !breaking {
			continue
		}

		if breaking {
			majorCount++
			continue
		}

		switch commitType {
		case "feat":
			minorCount++
		case "fix", "perf", "revert":
			patchCount++
			patchTypes[commitType]++
		}
	}

	if majorCount > 0 {
		return AnalysisResult{
			Bump:   BumpMajor,
			Reason: fmt.Sprintf("%d breaking commit(s) require a major bump", majorCount),
		}
	}

	if minorCount > 0 {
		return AnalysisResult{
			Bump:   BumpMinor,
			Reason: fmt.Sprintf("%d feat commit(s) require a minor bump", minorCount),
		}
	}

	if patchCount > 0 {
		return AnalysisResult{
			Bump:   BumpPatch,
			Reason: fmt.Sprintf("%d %s commit(s) require a patch bump", patchCount, summarizePatchTypes(patchTypes)),
		}
	}

	return result
}

func parseCommit(rawMessage string) (commitType string, breaking bool) {
	trimmed := strings.TrimSpace(rawMessage)
	if trimmed == "" {
		return "", false
	}

	if hasBreakingChangeFooter(trimmed) {
		breaking = true
	}

	header := firstLine(trimmed)
	matches := conventionalHeaderPattern.FindStringSubmatch(header)
	if len(matches) == 0 {
		return "", breaking
	}

	commitType = strings.ToLower(matches[1])
	if matches[3] == "!" {
		breaking = true
	}

	return commitType, breaking
}

func hasBreakingChangeFooter(message string) bool {
	for _, line := range strings.Split(message, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "BREAKING CHANGE:") {
			return true
		}
	}

	return false
}

func firstLine(message string) string {
	message = strings.TrimSpace(message)
	if message == "" {
		return ""
	}

	parts := strings.SplitN(message, "\n", 2)
	return strings.TrimSpace(parts[0])
}

func summarizePatchTypes(counts map[string]int) string {
	if len(counts) == 1 {
		for commitType := range counts {
			return commitType
		}
	}

	return "patch-level"
}
