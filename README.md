<!-- SPDX-License-Identifier: Apache-2.0 -->
<!-- SPDX-FileCopyrightText: 2026 The go-semrel Authors -->

# analyzer-conventional

`analyzer-conventional` is a SemRel commit-analyzer plugin that reads conventional commit messages and returns the semantic version bump required for a release.

## What the plugin does

- Parses conventional commit subjects from commit `raw_message` values
- Detects breaking changes from `!`, `BREAKING CHANGE:`, and `BREAKING-CHANGE:` markers
- Chooses the highest required bump across all commits
- Supports per-type bump overrides through the plugin config map

## Supported conventional commit formats

The analyzer recognizes these subject formats:

- `type(scope): description`
- `type!: description`
- `type(scope)!: description`

Examples:

~~~text
feat: add release notes rendering
fix(parser): preserve footer whitespace
refactor(core)!: remove legacy adapter
~~~

Breaking changes can also be declared in the body or footer:

~~~text
feat(api): change payload format

BREAKING CHANGE: payload consumers must update their schema
~~~

## Default bump mapping

| Commit type | Default bump |
| --- | --- |
| `feat` | `minor` |
| `fix` | `patch` |
| `build` | `none` |
| `chore` | `none` |
| `ci` | `none` |
| `docs` | `none` |
| `perf` | `none` |
| `refactor` | `none` |
| `style` | `none` |
| `test` | `none` |
| any breaking change | `major` |
| unknown / non-conventional | `none` |

## Configurable overrides

Override the default mapping by setting `bump.<type>` keys in the plugin config block.

Example `.semrel.yaml`:

~~~yaml
plugins:
  - name: analyzer-conventional
    type: analyzer
    config:
      bump.chore: patch
      bump.docs: minor
      bump.refactor: patch
~~~

Accepted override values are `none`, `patch`, `minor`, and `major`.

## Error scenarios

The plugin returns an error when:

- the gRPC request is missing `ctx`
- a `bump.<type>` override uses an unsupported value

Non-conventional or empty commit messages are handled gracefully and simply contribute `none`.

## Build, run, and test

~~~bash
go mod tidy
go build ./...
go test -v -cover ./...
make build
make test
~~~

Run the plugin locally:

~~~bash
go run cmd/plugin/main.go
~~~

The server binds to `127.0.0.1:0`, logs the chosen listening address to stderr, and serves the SemRel `CommitAnalyzerPlugin` gRPC API.
