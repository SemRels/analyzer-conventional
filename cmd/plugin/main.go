// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The analyzer-conventional Authors

package main

import (
	analyzerplugin "github.com/SemRels/analyzer-conventional/internal/plugin"
	semrelapi "github.com/SemRels/semrel-api/plugin"
	goplugin "github.com/hashicorp/go-plugin"
)

func main() {
	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: semrelapi.HandshakeConfig,
		Plugins: map[string]goplugin.Plugin{
			"analyzer": &semrelapi.CommitAnalyzerGRPCPlugin{
				Impl: analyzerplugin.New(),
			},
		},
		GRPCServer: goplugin.DefaultGRPCServer,
	})
}
