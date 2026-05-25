// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The go-semrel Authors

package grpc

import (
	"context"
	"testing"

	semrelv1 "github.com/SemRels/analyzer-conventional/api/gen/v1"
	"github.com/stretchr/testify/require"
)

func TestCommitAnalyzerServerAnalyzeCommits(t *testing.T) {
	t.Parallel()

	server := NewCommitAnalyzerServer()

	resp, err := server.AnalyzeCommits(context.Background(), &semrelv1.AnalyzeCommitsRequest{
		Ctx: &semrelv1.ReleaseContext{
			Commits: []*semrelv1.Commit{{RawMessage: "feat(cli): add grpc startup"}},
		},
	})
	require.NoError(t, err)
	require.Equal(t, semrelv1.BumpLevel_BUMP_LEVEL_MINOR, resp.GetBump())
	require.Contains(t, resp.GetReason(), "minor bump")
}

func TestCommitAnalyzerServerAnalyzeCommitsMissingContext(t *testing.T) {
	t.Parallel()

	server := NewCommitAnalyzerServer()

	resp, err := server.AnalyzeCommits(context.Background(), &semrelv1.AnalyzeCommitsRequest{})
	require.Nil(t, resp)
	require.ErrorContains(t, err, "missing release context")
}
