// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The analyzer-conventional Authors

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	plugin "github.com/SemRels/analyzer-conventional/internal/plugin"
)

func main() {
	os.Exit(run(os.Stdout, os.Stderr, os.Getenv))
}

func run(stdout, stderr io.Writer, getenv func(string) string) int {
	raw := getenv("SEMREL_COMMITS")

	var commits []string
	if raw != "" {
		if err := json.Unmarshal([]byte(raw), &commits); err != nil {
			_, _ = fmt.Fprintln(stderr, "analyzer-conventional: invalid SEMREL_COMMITS JSON:", err)
			return 1
		}
	}

	result := plugin.New().Analyze(commits)
	result.PluginSchemaVersion = plugin.PluginSchemaVersion

	if err := json.NewEncoder(stdout).Encode(result); err != nil {
		_, _ = fmt.Fprintln(stderr, "analyzer-conventional:", err)
		return 1
	}

	return 0
}
