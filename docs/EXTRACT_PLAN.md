# Plan: Extract Flow Launcher Plugin into Standalone Go Library

## New module

Create a new Go module `github.com/qk/flow-launcher-go` (or similar name) in a separate repository.

## Package structure

```
flow-launcher-go/
├── go.mod
├── plugin/
│   ├── plugin.go        # Plugin core: Dispatch, Handler, ToolDef, NewResult, Result types
│   ├── plugin_test.go
│   ├── cobra.go         # Cobra integration: attach to root command, JSON detection
│   └── cobra_test.go
│
├── api/
│   ├── api.go           # 47 typed IPublicAPI method structs (same as current flowapi)
│   └── api_test.go
│
├── jsonrpc/
│   ├── request.go       # Parse, validate
│   ├── response.go      # Response types
│   ├── error.go         # Error codes
│   ├── dispatcher.go    # Handler registry, dispatch, duck-type Method()/Params()
│   └── jsonrpc_test.go
│
├── manifest/
│   ├── manifest.go      # Plugin manifest struct, New() with defaults, Validate(), Marshal()
│   └── manifest_test.go
│
├── deploy/
│   └── deploy.go        # Build + copy to Flow's plugin dir (platform-aware paths)
│
├── cmd/
│   └── flow-cli/
│       └── main.go      # Optional: CLI tool to scaffold/generate/test plugins
│
└── examples/
    └── secret/
        ├── main.go       # Example plugin using the library
        ├── plugin.json
        └── deploy.sh
```

## What each package provides

### `plugin` — main entry point for plugin authors

```go
// Define a tool
plugin.Define(plugin.Tool{
    Name:        "secret",
    Aliases:     []string{"pw", "password"},
    Description: "Generate a random string",
    Example:     "secret 64 --alpha",
    Handler: func(ctx context.Context, args []string) (string, error) {
        // ... generate secret
        return result, nil
    },
})

// Attach to a cobra.Command (or use standalone)
err := plugin.Attach(rootCmd)  // sets ArbitraryArgs, RunE with JSON detection
```

### `api` — typed Flow API method structs

```go
api.CopyToClipboard{Text: "hello", DirectCopy: false, ShowDefaultNotification: true}
api.ShowMsg{Title: "Done", SubTitle: "Generated", IconPath: "icon.png"}
api.ShellRun{Cmd: "notepad.exe", Filename: "cmd.exe"}
```

All implement `api.Method` interface: `Method() string` + `Params() []any`.

### `manifest` — programmatic plugin.json

```go
m := manifest.New(manifest.Config{
    ID:             "uuid",
    ActionKeyword:  "nt",
    Name:           "My Plugin",
    Description:    "Does stuff",
    Author:         "me",
    Version:        "1.0.0",
    Language:       "executable",
    IcoPath:        "Images\\app.png",
    ExeFileName:    "my-plugin.exe",
})
json := m.Marshal() // returns pretty-printed JSON
```

### `deploy` — one-call deploy

```go
deploy.ToFlow("path/to/plugin-dir") // builds binary, copies files
```

## Dependency graph

```
nnry-text-util
  └── github.com/qk/flow-launcher-go
        ├── plugin (ToolDef, Dispatch, Attach, NewResult)
        ├── api    (APIMethod structs)
        ├── jsonrpc (Parse, Dispatcher, error codes)
        ├── manifest (plugin.json builder)
        └── deploy  (build + deploy logic)
```

## Migration phases

### Phase 1: Create new module
- `go mod init github.com/qk/flow-launcher-go`
- Copy `pkg/jsonrpc/` → `jsonrpc/` (no changes needed)

### Phase 2: Copy `pkg/flowapi/` → `api/`
- Rename package from `flowapi` to `api`
- Keep `APIMethod` interface, `Result`, `JsonRPCAction`, `NewResult`
- Keep all 47 method structs

### Phase 3: Create `plugin/`
- Extract tool registration, dispatch, `flowDispatcher` from `cmd/flow.go`
- Extract `NewResult`, `Result`, `JsonRPCAction` types (moved from `api/` to `plugin/` to avoid circular deps)
- Move `api.APIMethod` interface to `plugin` package (or keep in `api`)
- `plugin.Attach()` handles:
  - `rootCmd.Args = cobra.ArbitraryArgs`
  - `rootCmd.RunE` with JSON detection
  - Registering `query` and `context_menu` handlers
  - Registering `execute_*` handlers for each tool

### Phase 4: Create `manifest/`
- `Config` struct with all plugin.json fields
- `New()` with defaults (UUID generation, version)
- `Validate()` — check required fields
- `Marshal()` → pretty JSON
- `Unmarshal()` from existing `plugin.json`

### Phase 5: Create `deploy/`
- `ToFlow(pluginDir string)` — builds binary, copies to `%APPDATA%\FlowLauncher\Plugins\`
- Platform-aware: Windows `%APPDATA%`, Linux/macOS `~/.config/FlowLauncher/Plugins/`
- `Build()` — runs `go build` in the plugin directory

### Phase 6: Update `nnry-text-util`
- Update `go.mod` to depend on `github.com/qk/flow-launcher-go`
- `cmd/flow.go` becomes:
  ```go
  package cmd

  import "github.com/qk/flow-launcher-go/plugin"

  func init() {
      plugin.Define(plugin.Tool{
          Name:    "secret",
          Aliases: []string{"pw", "password"},
          Handler: handleToolSecret,
      })
      plugin.Attach(rootCmd)
  }
  ```
- Delete `pkg/jsonrpc/`, `pkg/flowapi/`, `cmd/flow.go` (replaced by library import)
- `cmd/secret.go` stays — `SecretOptions`, `BuildCharset`, `generateSecret` are app-specific

### Phase 7: Polish
- Add tests for all packages
- Add CI (GitHub Actions: `go test ./...`)
- Add `examples/secret/` — full working plugin
- Add godoc comments
- Publish module

## What stays in `nnry-text-util`

| File | Reason |
|------|--------|
| `cmd/secret.go` | App-specific logic: `SecretOptions`, `BuildCharset`, `generateSecret` |
| `cmd/flow.go` | Thin wrapper: calls `plugin.Define()` + `plugin.Attach()` |
| `flow-plugin/plugin.json` | Static manifest (or generated via `manifest` package) |
| `deploy-flow.sh` | Calls `deploy.ToFlow()` or remains a shell script |
| `cmd/root.go` | Remains as-is, but `plugin.Attach()` handles the JSON detection |

## Key design decisions

- **Duck-type interface** in `jsonrpc.Dispatcher` (`interface{ Method() string; Params() []any }`) instead of importing `api` — avoids circular deps.
- **`plugin` package** owns `Result`, `NewResult`, `JsonRPCAction` — these are plugin concepts, not API concepts.
- **`api` package** only has method structs — pure data, no Flow-specific types.
- **`manifest` as a separate package** — useful without importing the full plugin runtime.
- **`deploy` as a separate package** — can be used standalone in a deploy script or CI.