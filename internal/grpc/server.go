// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The go-semrel Authors

package grpc

import (
	"context"
	"fmt"
	"log"

	semrelv1 "github.com/SemRels/analyzer-conventional/api/gen/v1"
	"github.com/SemRels/analyzer-conventional/internal/plugin"
)

// CommitAnalyzerServer implements semrelv1.CommitAnalyzerPluginServer.
type CommitAnalyzerServer struct {
	semrelv1.UnimplementedCommitAnalyzerPluginServer
	analyzer *plugin.Analyzer
}

func NewCommitAnalyzerServer() *CommitAnalyzerServer {
	return &CommitAnalyzerServer{analyzer: plugin.NewAnalyzer()}
}

func (s *CommitAnalyzerServer) AnalyzeCommits(ctx context.Context, req *semrelv1.AnalyzeCommitsRequest) (*semrelv1.AnalyzeCommitsResponse, error) {
	_ = ctx
	if req.GetCtx() == nil {
		return nil, fmt.Errorf("missing release context")
	}

	result, err := s.analyzer.AnalyzeCommits(req.GetCtx().GetCommits(), req.GetCtx().GetConfig())
	if err != nil {
		return nil, err
	}

	log.Printf("analyzer result: bump=%s reason=%s", result.Bump, result.Reason)

	return &semrelv1.AnalyzeCommitsResponse{
		Bump:   result.Bump,
		Reason: result.Reason,
	}, nil
}
