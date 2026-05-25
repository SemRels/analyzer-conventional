// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The analyzer-conventional Authors

package plugin

import (
	"context"
	"testing"

	semrelv1 "github.com/SemRels/semrel-api/api/gen/v1"
	"github.com/stretchr/testify/require"
)

func TestAnalyzerAnalyzeCommitsMajorTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		message    string
		wantBump   semrelv1.BumpLevel
		wantReason bool
	}{
		{name: "feature", message: "feat: add search", wantBump: semrelv1.BumpLevel_BUMP_LEVEL_MINOR, wantReason: true},
		{name: "fix", message: "fix: resolve panic", wantBump: semrelv1.BumpLevel_BUMP_LEVEL_PATCH, wantReason: true},
		{name: "perf", message: "perf: speed up diff", wantBump: semrelv1.BumpLevel_BUMP_LEVEL_PATCH, wantReason: true},
		{name: "revert", message: "revert: feat: add search", wantBump: semrelv1.BumpLevel_BUMP_LEVEL_PATCH, wantReason: true},
		{name: "docs", message: "docs: update usage", wantBump: semrelv1.BumpLevel_BUMP_LEVEL_NONE, wantReason: false},
	}

	analyzer := New()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resp, err := analyzer.AnalyzeCommits(context.Background(), requestWithMessages(tt.message))
			require.NoError(t, err)
			require.Equal(t, tt.wantBump, resp.GetBump())
			if tt.wantReason {
				require.NotEmpty(t, resp.GetReason())
			}
		})
	}
}

func TestAnalyzerAnalyzeCommitsBreakingChangeFooter(t *testing.T) {
	t.Parallel()

	resp, err := New().AnalyzeCommits(context.Background(), requestWithMessages("feat: add config\n\nBREAKING CHANGE: config shape changed"))
	require.NoError(t, err)
	require.Equal(t, semrelv1.BumpLevel_BUMP_LEVEL_MAJOR, resp.GetBump())
	require.NotEmpty(t, resp.GetReason())
}

func TestAnalyzerAnalyzeCommitsBangHeaderIsMajor(t *testing.T) {
	t.Parallel()

	resp, err := New().AnalyzeCommits(context.Background(), requestWithMessages("feat!: remove deprecated endpoint"))
	require.NoError(t, err)
	require.Equal(t, semrelv1.BumpLevel_BUMP_LEVEL_MAJOR, resp.GetBump())
	require.NotEmpty(t, resp.GetReason())
}

func TestAnalyzerAnalyzeCommitsEmptyCommits(t *testing.T) {
	t.Parallel()

	resp, err := New().AnalyzeCommits(context.Background(), &semrelv1.AnalyzeCommitsRequest{Ctx: &semrelv1.ReleaseContext{}})
	require.NoError(t, err)
	require.Equal(t, semrelv1.BumpLevel_BUMP_LEVEL_NONE, resp.GetBump())
}

func TestAnalyzerAnalyzeCommitsNonConventionalCommits(t *testing.T) {
	t.Parallel()

	resp, err := New().AnalyzeCommits(context.Background(), requestWithMessages("merge branch 'main'", "update changelog manually"))
	require.NoError(t, err)
	require.Equal(t, semrelv1.BumpLevel_BUMP_LEVEL_NONE, resp.GetBump())
}

func TestAnalyzerAnalyzeCommitsHighestBumpWins(t *testing.T) {
	t.Parallel()

	resp, err := New().AnalyzeCommits(context.Background(), requestWithMessages("fix: resolve nil pointer", "feat(api): add pagination"))
	require.NoError(t, err)
	require.Equal(t, semrelv1.BumpLevel_BUMP_LEVEL_MINOR, resp.GetBump())
	require.NotEmpty(t, resp.GetReason())
}

func TestAnalyzerAnalyzeCommitsReasonPresentWhenBumpRequired(t *testing.T) {
	t.Parallel()

	resp, err := New().AnalyzeCommits(context.Background(), requestWithMessages("fix: resolve nil pointer"))
	require.NoError(t, err)
	require.Equal(t, semrelv1.BumpLevel_BUMP_LEVEL_PATCH, resp.GetBump())
	require.NotEmpty(t, resp.GetReason())
}

func requestWithMessages(messages ...string) *semrelv1.AnalyzeCommitsRequest {
	commits := make([]*semrelv1.Commit, 0, len(messages))
	for _, message := range messages {
		commits = append(commits, &semrelv1.Commit{RawMessage: message})
	}

	return &semrelv1.AnalyzeCommitsRequest{
		Ctx: &semrelv1.ReleaseContext{Commits: commits},
	}
}
