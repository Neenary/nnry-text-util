# AGENTS.md — nnry-text-util

A CLI text processing utility built with Go and [Cobra](https://github.com/spf13/cobra).

## Quick Start

```bash
go build -o nnry-text-util.exe .    # build binary
go run main.go                       # run directly
go vet ./...                         # lint
go test ./...                        # run tests (none exist yet)
```

## Project Status

**Early scaffold.** `main.go`, `cmd/root.go`, and `cmd/secret.go` exist. The `secret` subcommand is implemented. More text-processing commands can be added.

## Code Organization

```
main.go          — entry point, calls cmd.Execute()
cmd/root.go      — root cobra.Command definition
cmd/secret.go    — `secret` subcommand for generating random strings
```

- **Entry point**: `main.go` → `cmd.Execute()`
- **Commands**: Each subcommand lives in its own file under `cmd/` (e.g., `cmd/encode.go`, `cmd/decode.go`). Follow the pattern in `cmd/root.go` — use `cobra.Command` structs, register in `init()` with `rootCmd.AddCommand()`.
- **Shared logic**: Extract reusable text-processing functions into `pkg/` packages (e.g., `pkg/encoding/`, `pkg/transform/`). Keep `cmd/` thin — just wiring flags to function calls.

## Stack

- **Go 1.26.5** (on Windows)
- **Cobra v1.10.2** — CLI framework
- **Module**: `github.com/qk/nnry-text-util`

## Conventions

- **Package naming**: `cmd` for commands, `pkg/<name>` for library code.
- **Command files**: `cmd/<name>.go` — one file per subcommand.
- **No `_test.go` files exist yet** — add them alongside the code they test.
- **Error handling**: Cobra commands return errors; the root `Execute()` calls `os.Exit(1)` on error. Follow this pattern.
- **Flag style**: Uses `pflag` (Cobra's default). Prefer long flags (`--name`) with short single-char aliases for common flags.

## Gotchas

- **Windows binary**: The `.gitignore` excludes `nnry-text-util.exe` — the binary name is platform-specific.
- **No `init()` in `main.go`**: All command registration happens in `cmd/` package `init()` functions.
- **Toggle flag**: The default scaffold includes a `--toggle` / `-t` flag on the root command. Remove it when adding real commands.
- **No CI/CD config**: None exists. If adding GitHub Actions, the build command is `go build ./...` and test is `go test ./...`.
- **No Makefile**: Build via `go build`. If the project grows, consider adding a `Makefile` (or keep using raw `go` commands — the project is small).