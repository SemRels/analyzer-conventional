# analyzer-conventional

Determines the semantic version bump from Conventional Commit messages.

This plugin is distributed as the standalone Go binary `semrel-plugin-analyzer-conventional`. Semrel executes the binary as a subprocess, provides plugin configuration through `SEMREL_PLUGIN_*` environment variables, provides release context through `SEMREL_*` environment variables, reads standard output, and treats exit code `0` as success and any non-zero exit code as failure. Install the binary in `~/.semrel/plugins/` or anywhere on your `$PATH`.

## Installation

### Binary

```bash
go install github.com/SemRels/analyzer-conventional/cmd/plugin@latest
```

### Docker

Pre-built, multi-platform images (linux/amd64, linux/arm64) are published to the GitHub Container Registry on every release:

```bash
docker pull ghcr.io/semrels/analyzer-conventional:latest
```

Images are signed with [cosign](https://github.com/sigstore/cosign) and include a full SBOM attestation. Verify the signature:

```bash
cosign verify ghcr.io/semrels/analyzer-conventional:latest \
  --certificate-identity-regexp 'https://github.com/SemRels/analyzer-conventional/.github/workflows/release.yml.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
```


## Configuration

```yaml
plugins:
  - name: analyzer-conventional
    path: ~/.semrel/plugins/semrel-plugin-analyzer-conventional
    env:
      SEMREL_PLUGIN_BREAKING_CHANGE_LABEL: "BREAKING CHANGE"
      SEMREL_PLUGIN_MINOR_TYPES: "feat"
      SEMREL_PLUGIN_PATCH_TYPES: "fix,perf,refactor"
```

## `SEMREL_PLUGIN_*` variables

| Name | Required | Description | Default |
| --- | --- | --- | --- |
| `SEMREL_PLUGIN_BREAKING_CHANGE_LABEL` | Optional | Label used to detect breaking changes in commit messages. | BREAKING CHANGE |
| `SEMREL_PLUGIN_MINOR_TYPES` | Optional | Comma-separated commit types that trigger a minor release. | feat |
| `SEMREL_PLUGIN_PATCH_TYPES` | Optional | Comma-separated commit types that trigger a patch release. | fix,perf,refactor |

## `SEMREL_*` release context used

| Variable | Description |
| --- | --- |
| `SEMREL_BUMP` | Calculated bump level such as major, minor, or patch. |

## Example behavior

The plugin reads commit history, classifies commits using Conventional Commit rules, and emits the resulting bump decision for semrel to use.

### Breaking changes

Two forms trigger a **major** bump:

1. The `!` shorthand after type or scope (Conventional Commits v1.0.0 §11):

   ```
   feat!: remove deprecated authentication endpoint
   fix!: drop support for Node.js 14
   feat(api)!: rename /users to /accounts
   ```

2. A `BREAKING CHANGE:` footer in the commit body:

   ```
   feat: overhaul config

   BREAKING CHANGE: config schema changed
   ```

Both forms can appear together in the same commit.

## License

Apache-2.0
