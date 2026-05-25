# analyzer-conventional

Conventional commits analyzer plugin for SemRel.

Determines release types from conventional commit messages and configurable type-to-bump rules.

## Documentation

- SemRel docs (planned): <https://github.com/SemRels/semrel/tree/main/docs/plugins/analyzer-conventional>
- Plugin template: <https://github.com/SemRels/plugin-template>
- Registry: <https://registry.semrel.io>

## Repository Layout

~~~text
cmd/plugin/              Plugin entry point
internal/plugin/         Business logic scaffold
internal/grpc/           gRPC transport scaffold
proto/v1                 Symlink to the SemRel protobuf contract
.github/workflows/       CI, release, and security automation
~~~

## Development

~~~bash
go build ./cmd/plugin
go test ./...
~~~

## Configuration Example

~~~yaml
plugins:
  - name: analyzer-conventional
    type: analyzer
    config:
      preset: conventionalcommits
      breaking_keywords:
        - BREAKING CHANGE
      type_map:
        feat: minor
        fix: patch
~~~

## Status

This repository is bootstrapped from SemRels/plugin-template and is ready for implementation.
