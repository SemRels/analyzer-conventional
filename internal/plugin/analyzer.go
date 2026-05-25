// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The analyzer-conventional Authors

// Package plugin implements CommitAnalyzerPlugin for Conventional Commits.
package plugin

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	semrelv1 "github.com/SemRels/semrel-api/api/gen/v1"
)

var conventionalHeaderPattern = regexp.MustCompile(`^(\w+)(\([\w-]+\))?(!)?:(.+)$`)

// Analyzer implements the SemRels CommitAnalyzerPlugin gRPC service.
type Analyzer struct {
	semrelv1.UnimplementedCommitAnalyzerPluginServer
}

// New returns a Conventional Commits analyzer plugin implementation.
func New() *Analyzer {
	return &Analyzer{}
}

// AnalyzeCommits determines the highest semantic version bump required by the release context commits.
func (a *Analyzer) AnalyzeCommits(_ context.Context, req *semrelv1.AnalyzeCommitsRequest) (*semrelv1.AnalyzeCommitsResponse, error) {
	response := &semrelv1.AnalyzeCommitsResponse{
		Bump:   semrelv1.BumpLevel_BUMP_LEVEL_NONE,
		Reason: "no conventional commits require a version bump",
	}

	if req == nil || req.GetCtx() == nil || len(req.GetCtx().GetCommits()) == 0 {
		return response, nil
	}

	majorCount := 0
	minorCount := 0
	patchCount := 0
	patchTypes := map[string]int{}

	for _, commit := range req.GetCtx().GetCommits() {
		if commit == nil {
			continue
		}

		commitType, breaking := parseCommit(commit.GetRawMessage())
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
		response.Bump = semrelv1.BumpLevel_BUMP_LEVEL_MAJOR
		response.Reason = fmt.Sprintf("%d breaking commit(s) require a MAJOR bump", majorCount)
		return response, nil
	}

	if minorCount > 0 {
		response.Bump = semrelv1.BumpLevel_BUMP_LEVEL_MINOR
		response.Reason = fmt.Sprintf("%d feat commit(s) require a MINOR bump", minorCount)
		return response, nil
	}

	if patchCount > 0 {
		response.Bump = semrelv1.BumpLevel_BUMP_LEVEL_PATCH
		response.Reason = fmt.Sprintf("%d %s commit(s) require a PATCH bump", patchCount, summarizePatchTypes(patchTypes))
	}

	return response, nil
}

func parseCommit(rawMessage string) (commitType string, breaking bool) {
	trimmed := strings.TrimSpace(rawMessage)
	if trimmed == "" {
		return "", false
	}

	if hasBreakingChangeFooter(trimmed) {
		breaking = true
	}

	header := trimmed
	if idx := strings.Index(header, "\n"); idx >= 0 {
		header = header[:idx]
	}

	matches := conventionalHeaderPattern.FindStringSubmatch(strings.TrimSpace(header))
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

func summarizePatchTypes(counts map[string]int) string {
	if len(counts) == 1 {
		for commitType := range counts {
			return commitType
		}
	}

	return "patch-level"
}
