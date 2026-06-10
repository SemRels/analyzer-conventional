// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The analyzer-conventional Authors

package plugin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnalyzerAnalyze(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		commits    []string
		wantBump   BumpLevel
		wantReason string
	}{
		{
			name:       "major from breaking footer",
			commits:    []string{"feat: add config\n\nBREAKING CHANGE: config schema changed"},
			wantBump:   BumpMajor,
			wantReason: "breaking commit(s)",
		},
		{
			name:       "major from bang header",
			commits:    []string{"feat!: remove deprecated endpoint"},
			wantBump:   BumpMajor,
			wantReason: "breaking commit(s)",
		},
		{
			name:       "major from bang with scope",
			commits:    []string{"feat(api)!: rename /users to /accounts"},
			wantBump:   BumpMajor,
			wantReason: "breaking commit(s)",
		},
		{
			name:       "major from fix bang",
			commits:    []string{"fix!: drop support for Node.js 14"},
			wantBump:   BumpMajor,
			wantReason: "breaking commit(s)",
		},
		{
			name:       "minor from feat",
			commits:    []string{"feat(api): add pagination"},
			wantBump:   BumpMinor,
			wantReason: "feat commit(s)",
		},
		{
			name:       "patch from fix",
			commits:    []string{"fix: resolve panic"},
			wantBump:   BumpPatch,
			wantReason: "fix commit(s)",
		},
		{
			name:       "none for non conventional commits",
			commits:    []string{"docs: update usage", "merge branch 'main'"},
			wantBump:   BumpNone,
			wantReason: "no conventional commits",
		},
		{
			name:       "highest bump wins",
			commits:    []string{"fix: resolve panic", "feat(ui): add filters"},
			wantBump:   BumpMinor,
			wantReason: "feat commit(s)",
		},
		{
			name:       "mixed patch types summarized",
			commits:    []string{"fix: resolve panic", "perf: avoid extra allocations"},
			wantBump:   BumpPatch,
			wantReason: "patch-level commit(s)",
		},
	}

	analyzer := New()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := analyzer.Analyze(tt.commits)
			require.Equal(t, tt.wantBump, result.Bump)
			require.Contains(t, result.Reason, tt.wantReason)
		})
	}
}

func TestParseCommit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		commit       string
		wantType     string
		wantBreaking bool
	}{
		{
			name:         "empty",
			commit:       "   ",
			wantBreaking: false,
		},
		{
			name:         "non conventional with breaking footer",
			commit:       "update docs\n\nBREAKING CHANGE: rewrite config",
			wantBreaking: true,
		},
		{
			name:         "conventional fix",
			commit:       "fix(scope): patch issue",
			wantType:     "fix",
			wantBreaking: false,
		},
		{
			name:         "bang without scope triggers breaking",
			commit:       "feat!: remove deprecated endpoint",
			wantType:     "feat",
			wantBreaking: true,
		},
		{
			name:         "bang with scope triggers breaking",
			commit:       "feat(api)!: rename /users to /accounts",
			wantType:     "feat",
			wantBreaking: true,
		},
		{
			name:         "fix bang without scope triggers breaking",
			commit:       "fix!: drop support for Node.js 14",
			wantType:     "fix",
			wantBreaking: true,
		},
		{
			name:         "fix bang with scope triggers breaking",
			commit:       "fix(auth)!: remove legacy token format",
			wantType:     "fix",
			wantBreaking: true,
		},
		{
			name:         "bang and breaking footer both present",
			commit:       "feat!: major overhaul\n\nBREAKING CHANGE: all endpoints renamed",
			wantType:     "feat",
			wantBreaking: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotType, gotBreaking := parseCommit(tt.commit)
			require.Equal(t, tt.wantType, gotType)
			require.Equal(t, tt.wantBreaking, gotBreaking)
		})
	}
}

func TestSummarizePatchTypes(t *testing.T) {
	t.Parallel()

	require.Equal(t, "fix", summarizePatchTypes(map[string]int{"fix": 1}))
	require.Equal(t, "patch-level", summarizePatchTypes(map[string]int{"fix": 1, "perf": 1}))
}
