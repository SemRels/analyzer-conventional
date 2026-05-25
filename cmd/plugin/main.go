// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The go-semrel Authors

package main

import (
	"log"
	"net"

	semrelv1 "github.com/SemRels/analyzer-conventional/api/gen/v1"
	grpcserver "github.com/SemRels/analyzer-conventional/internal/grpc"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	semrelv1.RegisterCommitAnalyzerPluginServer(srv, grpcserver.NewCommitAnalyzerServer())

	log.Printf("analyzer-conventional gRPC server listening on %s", lis.Addr())
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
