# analyzer-conventional

`analyzer-conventional` is a SemRels commit analyzer plugin that inspects Conventional Commit messages from the release context and returns the highest required semantic version bump.

## Bump rules

- `BREAKING CHANGE:` footer or `!` after the type/scope (for example `feat!:` or `feat(api)!:`) → `MAJOR`
- `feat:` or `feat(scope):` → `MINOR`
- `fix:`, `perf:`, `revert:` → `PATCH`
- other conventional types such as `docs:`, `chore:`, `test:`, `ci:` and non-conventional commits → no bump

## Usage

Build the plugin binary:

~~~bash
go build -o analyzer-conventional.exe ./cmd/plugin
~~~

SemRels loads the binary through `hashicorp/go-plugin` using the shared handshake from `github.com/SemRels/semrel-api`.

## Development

~~~bash
go mod tidy
go build ./...
CGO_ENABLED=0 go test ./...
~~~

## Repository layout

~~~text
cmd/plugin/              go-plugin entry point
internal/plugin/         Conventional Commits analyzer implementation
internal/grpc/           placeholder package; transport is handled by go-plugin
proto/                   SemRels protobuf assets (currently unused by this binary)
~~~
