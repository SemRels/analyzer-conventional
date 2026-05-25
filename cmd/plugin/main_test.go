// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The analyzer-conventional Authors

package main

import (
	"bytes"
	"encoding/json"
	"testing"

	plugin "github.com/SemRels/analyzer-conventional/internal/plugin"
	"github.com/stretchr/testify/require"
)

func TestRunWritesAnalysisResult(t *testing.T) {
	t.Parallel()

	getenv := func(key string) string {
		if key == "SEMREL_COMMITS" {
			return `["fix: resolve panic","feat: add search"]`
		}
		return ""
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run(&stdout, &stderr, getenv)

	require.Equal(t, 0, code)
	require.Empty(t, stderr.String())

	var result plugin.AnalysisResult
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &result))
	require.Equal(t, plugin.BumpMinor, result.Bump)
}

func TestRunRejectsInvalidCommitJSON(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run(&stdout, &stderr, func(string) string { return `[` })

	require.Equal(t, 1, code)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), "invalid SEMREL_COMMITS JSON")
}

func TestRunAllowsMissingCommitEnv(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run(&stdout, &stderr, func(string) string { return "" })

	require.Equal(t, 0, code)
	require.Empty(t, stderr.String())
	require.Contains(t, stdout.String(), "\"bump\":\"none\"")
}
