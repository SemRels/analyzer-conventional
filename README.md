# analyzer-conventional

Conventional commit analyzer plugin for Semantic Release.

Analyzes conventional commits to determine the next Semantic Release version.

## Documentation

- Docs (coming soon): <https://github.com/SemRels/semrel/tree/main/docs/plugins/analyzer-conventional>
- Template source: <https://github.com/SemRels/plugin-template>

## Repository Layout

`	ext
cmd/plugin/              Plugin entry point
internal/plugin/         Business logic scaffold
internal/grpc/           gRPC transport scaffold
proto/v1                 Symlink to the SemRel protobuf contract
.github/workflows/       CI, release, and security automation
`

## Development

`ash
go build ./cmd/plugin
go test ./...
`

## Configuration Example

`yaml
plugins:
  - name: analyzer-conventional
    type: analyzer
    config:
      preset: conventionalcommits
      release_rules:
        - type: feat
          release: minor
`

## Status

This repository is bootstrapped from SemRels/plugin-template and is ready for implementation.