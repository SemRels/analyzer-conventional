// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The go-semrel Authors

package plugin

import (
	"fmt"
	"regexp"
	"strings"

	semrelv1 "github.com/SemRels/analyzer-conventional/api/gen/v1"
)

var (
	conventionalCommitPattern = regexp.MustCompile(`^([a-zA-Z][a-zA-Z0-9-]*)(?:\([^)]+\))?(!)?:\s+(.+)$`)
	breakingFooterPattern     = regexp.MustCompile(`(?mi)^BREAKING(?: |-)?CHANGE:\s+.+$`)
	breakingBodyPattern       = regexp.MustCompile(`(?mi)\bBREAKING CHANGE\b`)
)

var defaultBumpMapping = map[string]semrelv1.BumpLevel{
	"build":    semrelv1.BumpLevel_BUMP_LEVEL_NONE,
	"chore":    semrelv1.BumpLevel_BUMP_LEVEL_NONE,
	"ci":       semrelv1.BumpLevel_BUMP_LEVEL_NONE,
	"docs":     semrelv1.BumpLevel_BUMP_LEVEL_NONE,
	"feat":     semrelv1.BumpLevel_BUMP_LEVEL_MINOR,
	"fix":      semrelv1.BumpLevel_BUMP_LEVEL_PATCH,
	"perf":     semrelv1.BumpLevel_BUMP_LEVEL_NONE,
	"refactor": semrelv1.BumpLevel_BUMP_LEVEL_NONE,
	"style":    semrelv1.BumpLevel_BUMP_LEVEL_NONE,
	"test":     semrelv1.BumpLevel_BUMP_LEVEL_NONE,
}

// AnalysisResult describes the highest bump required by a commit set.
type AnalysisResult struct {
	Bump   semrelv1.BumpLevel
	Reason string
}

// Analyzer inspects commit messages and maps them to semantic version bumps.
type Analyzer struct {
	config map[string]string
}

// NewAnalyzer creates an analyzer with optional default configuration.
func NewAnalyzer(config ...map[string]string) *Analyzer {
	merged := make(map[string]string)
	if len(config) > 0 {
		for key, value := range config[0] {
			merged[key] = value
		}
	}

	return &Analyzer{config: merged}
}

// AnalyzeCommits determines the highest required bump across all commits.
func (a *Analyzer) AnalyzeCommits(commits []*semrelv1.Commit, config map[string]string) (*AnalysisResult, error) {
	bumpMapping, err := a.bumpMapping(config)
	if err != nil {
		return nil, err
	}

	result := &AnalysisResult{
		Bump:   semrelv1.BumpLevel_BUMP_LEVEL_NONE,
		Reason: "no release-triggering conventional commits found",
	}

	if len(commits) == 0 {
		return result, nil
	}

	for _, commit := range commits {
		if commit == nil || strings.TrimSpace(commit.GetRawMessage()) == "" {
			continue
		}

		commitResult, matched := analyzeCommit(commit, bumpMapping)
		if !matched {
			continue
		}

		if result.Bump == semrelv1.BumpLevel_BUMP_LEVEL_NONE && commitResult.Bump == semrelv1.BumpLevel_BUMP_LEVEL_NONE {
			result.Reason = commitResult.Reason
		}

		if commitResult.Bump > result.Bump {
			result = commitResult
		}
	}

	return result, nil
}

func (a *Analyzer) bumpMapping(runtimeConfig map[string]string) (map[string]semrelv1.BumpLevel, error) {
	mapping := make(map[string]semrelv1.BumpLevel, len(defaultBumpMapping))
	for commitType, bump := range defaultBumpMapping {
		mapping[commitType] = bump
	}

	for _, config := range []map[string]string{a.config, runtimeConfig} {
		for key, value := range config {
			if !strings.HasPrefix(key, "bump.") {
				continue
			}

			commitType := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(key, "bump.")))
			if commitType == "" {
				return nil, fmt.Errorf("invalid bump override key %q", key)
			}

			bump, ok := parseBumpLevel(value)
			if !ok {
				return nil, fmt.Errorf("invalid bump override %q for type %q", value, commitType)
			}

			mapping[commitType] = bump
		}
	}

	return mapping, nil
}

func analyzeCommit(commit *semrelv1.Commit, bumpMapping map[string]semrelv1.BumpLevel) (*AnalysisResult, bool) {
	subject, body := splitCommitMessage(commit.GetRawMessage())
	matches := conventionalCommitPattern.FindStringSubmatch(subject)
	if matches == nil {
		return nil, false
	}

	commitType := strings.ToLower(matches[1])
	if matches[2] == "!" || hasBreakingChange(body) {
		return &AnalysisResult{
			Bump:   semrelv1.BumpLevel_BUMP_LEVEL_MAJOR,
			Reason: fmt.Sprintf("commit %q introduces a breaking change", subject),
		}, true
	}

	bump, ok := bumpMapping[commitType]
	if !ok {
		bump = semrelv1.BumpLevel_BUMP_LEVEL_NONE
	}

	if bump == semrelv1.BumpLevel_BUMP_LEVEL_NONE {
		return &AnalysisResult{
			Bump:   semrelv1.BumpLevel_BUMP_LEVEL_NONE,
			Reason: fmt.Sprintf("commit %q does not require a release", subject),
		}, true
	}

	return &AnalysisResult{
		Bump:   bump,
		Reason: fmt.Sprintf("commit %q maps %s to a %s bump", subject, commitType, strings.ToLower(strings.TrimPrefix(bump.String(), "BUMP_LEVEL_"))),
	}, true
}

func splitCommitMessage(message string) (string, string) {
	normalized := strings.ReplaceAll(message, "\r\n", "\n")
	parts := strings.SplitN(normalized, "\n", 2)
	subject := strings.TrimSpace(parts[0])
	if len(parts) == 1 {
		return subject, ""
	}

	return subject, strings.TrimSpace(parts[1])
}

func hasBreakingChange(body string) bool {
	if body == "" {
		return false
	}

	return breakingFooterPattern.MatchString(body) || breakingBodyPattern.MatchString(body)
}

func parseBumpLevel(value string) (semrelv1.BumpLevel, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "none":
		return semrelv1.BumpLevel_BUMP_LEVEL_NONE, true
	case "patch":
		return semrelv1.BumpLevel_BUMP_LEVEL_PATCH, true
	case "minor":
		return semrelv1.BumpLevel_BUMP_LEVEL_MINOR, true
	case "major":
		return semrelv1.BumpLevel_BUMP_LEVEL_MAJOR, true
	default:
		return semrelv1.BumpLevel_BUMP_LEVEL_UNSPECIFIED, false
	}
}
