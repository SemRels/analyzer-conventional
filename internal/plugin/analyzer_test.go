// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The go-semrel Authors

package plugin

import (
	"testing"

	semrelv1 "github.com/SemRels/analyzer-conventional/api/gen/v1"
	"github.com/stretchr/testify/require"
)

func TestAnalyzerAnalyzeCommits(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		analyzerConfig map[string]string
		runtimeConfig  map[string]string
		commits        []*semrelv1.Commit
		wantBump       semrelv1.BumpLevel
		wantReason     string
	}{
		{
			name:       "basic feat commit produces minor",
			commits:    []*semrelv1.Commit{{RawMessage: "feat: add plugin support"}},
			wantBump:   semrelv1.BumpLevel_BUMP_LEVEL_MINOR,
			wantReason: "maps feat to a minor bump",
		},
		{
			name:       "basic fix commit produces patch",
			commits:    []*semrelv1.Commit{{RawMessage: "fix: handle nil context"}},
			wantBump:   semrelv1.BumpLevel_BUMP_LEVEL_PATCH,
			wantReason: "maps fix to a patch bump",
		},
		{
			name:       "breaking change bang produces major",
			commits:    []*semrelv1.Commit{{RawMessage: "feat!: remove deprecated endpoint"}},
			wantBump:   semrelv1.BumpLevel_BUMP_LEVEL_MAJOR,
			wantReason: "breaking change",
		},
		{
			name:       "breaking change footer produces major",
			commits:    []*semrelv1.Commit{{RawMessage: "feat(api): extend payload\n\nBREAKING CHANGE: payload schema changed"}},
			wantBump:   semrelv1.BumpLevel_BUMP_LEVEL_MAJOR,
			wantReason: "breaking change",
		},
		{
			name: "multiple commits use highest bump",
			commits: []*semrelv1.Commit{
				{RawMessage: "fix: patch bug"},
				{RawMessage: "feat(ui): add dashboard"},
			},
			wantBump:   semrelv1.BumpLevel_BUMP_LEVEL_MINOR,
			wantReason: "maps feat to a minor bump",
		},
		{
			name: "chore and docs produce none",
			commits: []*semrelv1.Commit{
				{RawMessage: "chore: update lockfile"},
				{RawMessage: "docs(readme): clarify usage"},
			},
			wantBump:   semrelv1.BumpLevel_BUMP_LEVEL_NONE,
			wantReason: "does not require a release",
		},
		{
			name:           "config override promotes chore to patch",
			analyzerConfig: map[string]string{"bump.chore": "patch"},
			commits:        []*semrelv1.Commit{{RawMessage: "chore: ship maintenance release"}},
			wantBump:       semrelv1.BumpLevel_BUMP_LEVEL_PATCH,
			wantReason:     "maps chore to a patch bump",
		},
		{
			name:       "empty commits produce none",
			commits:    nil,
			wantBump:   semrelv1.BumpLevel_BUMP_LEVEL_NONE,
			wantReason: "no release-triggering conventional commits found",
		},
		{
			name:       "invalid commit is ignored",
			commits:    []*semrelv1.Commit{{RawMessage: "merge branch feature/foo"}},
			wantBump:   semrelv1.BumpLevel_BUMP_LEVEL_NONE,
			wantReason: "no release-triggering conventional commits found",
		},
		{
			name: "scope and no scope variations are supported",
			commits: []*semrelv1.Commit{
				{RawMessage: "fix(parser): trim spaces"},
				{RawMessage: "feat: add report output"},
				{RawMessage: "refactor(core)!: drop legacy adapter"},
			},
			wantBump:   semrelv1.BumpLevel_BUMP_LEVEL_MAJOR,
			wantReason: "breaking change",
		},
		{
			name:          "runtime config overrides analyzer defaults",
			runtimeConfig: map[string]string{"bump.docs": "minor"},
			commits:       []*semrelv1.Commit{{RawMessage: "docs: publish API guide"}},
			wantBump:      semrelv1.BumpLevel_BUMP_LEVEL_MINOR,
			wantReason:    "maps docs to a minor bump",
		},
		{
			name:       "breaking change body marker produces major",
			commits:    []*semrelv1.Commit{{RawMessage: "fix(core): preserve ordering\n\nThis includes a BREAKING CHANGE for callers."}},
			wantBump:   semrelv1.BumpLevel_BUMP_LEVEL_MAJOR,
			wantReason: "breaking change",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			analyzer := NewAnalyzer(tt.analyzerConfig)

			result, err := analyzer.AnalyzeCommits(tt.commits, tt.runtimeConfig)
			require.NoError(t, err)
			require.Equal(t, tt.wantBump, result.Bump)
			require.Contains(t, result.Reason, tt.wantReason)
		})
	}
}

func TestAnalyzerAnalyzeCommitsInvalidOverride(t *testing.T) {
	t.Parallel()

	analyzer := NewAnalyzer()

	result, err := analyzer.AnalyzeCommits([]*semrelv1.Commit{{RawMessage: "fix: patch bug"}}, map[string]string{"bump.fix": "banana"})
	require.Nil(t, result)
	require.ErrorContains(t, err, "invalid bump override")
}

func TestAnalyzeCommitUnknownTypeDefaultsToNone(t *testing.T) {
	t.Parallel()

	result, matched := analyzeCommit(&semrelv1.Commit{RawMessage: "custom(scope): add extension"}, defaultBumpMapping)
	require.True(t, matched)
	require.Equal(t, semrelv1.BumpLevel_BUMP_LEVEL_NONE, result.Bump)
	require.Contains(t, result.Reason, "does not require a release")
}

func TestSplitCommitMessageNormalizesWindowsNewlines(t *testing.T) {
	t.Parallel()

	subject, body := splitCommitMessage("feat(parser): add support\r\n\r\nBREAKING-CHANGE: old format removed")
	require.Equal(t, "feat(parser): add support", subject)
	require.Equal(t, "BREAKING-CHANGE: old format removed", body)
	require.True(t, hasBreakingChange(body))
}

func TestHasBreakingChangeEmptyBody(t *testing.T) {
	t.Parallel()

	require.False(t, hasBreakingChange(""))
}
