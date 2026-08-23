# AGENTS.md — nnry-text-util

A CLI text processing utility built with Go and [Cobra](https://github.com/spf13/cobra).

## Quick Start

```bash
go build -o nnry-text-util.exe .    # build binary
go run main.go                       # run directly
go vet ./...                         # lint
go test ./...                        # run tests
```

## Project Status

Active. `secret` subcommand is implemented. The Flow Launcher plugin system has been extracted into a standalone Go module at `flow-launcher-go/`.

## Code Organization

```
main.go                  — entry point, calls cmd.Execute()
cmd/
  root.go                — root cobra.Command definition
  secret.go              — `secret` subcommand, SecretOptions, BuildCharset, generateSecret
  flow.go                — thin wrapper: registers tools with plugin, attaches to rootCmd
jsonrpc/                 — 📁 standalone Go module (github.com/neenary/jsonrpc)
  request.go             — JSON-RPC request parsing/validation
  response.go            — JSON-RPC response types
  error.go               — JSON-RPC error codes
  dispatcher.go          — handler registry, dispatch, generic parameter parsing
flow-launcher-go/        — 📁 standalone Go module (github.com/neenary/flow-launcher-go)
  api/                   — typed Flow Launcher IPublicAPI method structs (47 methods)
  plugin/                — plugin framework: ToolDef, Plugin, query/context_menu dispatch
  manifest/              — plugin.json builder/validator
  deploy/                — build + deploy to Flow Launcher plugin dir
flow-plugin/
  plugin.json            — static manifest
  Images/                — plugin icons
```

## Stack

- **Go 1.26.5** (on Windows)
- **Cobra v1.10.2** — CLI framework
- **Module**: `github.com/neenary/nnry-text-util`
- **Submodules** (local replace in go.mod):
  - `github.com/neenary/jsonrpc` — JSON-RPC 2.0 library (zero deps)
  - `github.com/neenary/flow-launcher-go` — Flow Launcher plugin framework (depends on jsonrpc)

## Conventions

- **Package naming**: `cmd` for commands, `flow-launcher-go/<pkg>` for library code.
- **Command files**: `cmd/<name>.go` — one file per subcommand.
- **Error handling**: Cobra commands return errors; the root `Execute()` calls `os.Exit(1)` on error. Follow this pattern.
- **Flag style**: Uses `pflag` (Cobra's default). Prefer long flags (`--name`) with short single-char aliases for common flags.

## Gotchas

- **Windows binary**: The `.gitignore` excludes `nnry-text-util.exe` — the binary name is platform-specific.
- **No `init()` in `main.go`**: All command registration happens in `cmd/` package `init()` functions.
- **No CI/CD config**: None exists. If adding GitHub Actions, the build command is `go build ./...` and test is `go test ./...`.
- **No Makefile**: Build via `go build`. If the project grows, consider adding a `Makefile` (or keep using raw `go` commands).
- **Submodules**: Both `jsonrpc/` and `flow-launcher-go/` are standalone Go modules. They can be published to GitHub independently. The root `go.mod` uses `replace` directives to point at local copies. `flow-launcher-go` also has a `replace` for `jsonrpc`.