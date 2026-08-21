# Journey Notes: Flow Launcher Plugin (`flow` command)

## Overview

Added a `flow` command to `nnry-text-util` that turns it into a Flow Launcher `executable` v1 plugin, plus a `pkg/flowapi/` for typed IPublicAPI method schemas and a `pkg/jsonrpc/` for generic JSON-RPC v2 parsing.

---

## What Went Right

### 1. Generic JSON-RPC parser (`pkg/jsonrpc/`)
- Kept `Parse()` as a standalone function with no global state — tests are easy, no init ordering issues.
- Detailed error messages (`missing 'method' field`, `'parameters' must be a JSON array`) make debugging obvious.
- Duck-typed `Method()/Params()` interface detection in `Dispatch()` — any struct with those methods auto-marshals to `{"method":"...","parameters":[...]}`. This is what lets `flowapi.CopyToClipboard{...}` work directly as a return value.

### 2. Typed API structs (`pkg/flowapi/`)
- One struct per IPublicAPI method, each implementing `APIMethod` interface.
- `NewResult(title, subtitle, api)` makes it impossible to mistype a method name in query results.
- The interface is the same shape as the duck-type check in the dispatcher — no circular imports needed.

### 3. Extracted `BuildCharset` + `SecretOptions`
- Shared between CLI (`cmd/secret.go`) and Flow handler (`cmd/flow.go`).
- Both paths support the same flags: `--alpha`, `--numeric`, `--alphanumeric`, `--no-symbols`, `--no-digits`, `--no-lower`, `--no-upper`, `--length`.

### 4. Root command detection
- `rootCmd.Args = cobra.ArbitraryArgs` + first arg JSON detection.
- Standalone CLI (`secret --length 16`) works unchanged.
- JSON-RPC (`{"method":"query"...}`) auto-routes to the flow dispatcher.

### 5. `deploy-flow.sh` is separate from `build.sh`
- `build.sh` stays focused on the binary.
- Deployment is a separate concern — easy to CI/CD later.

---

## What Went Wrong & How It Was Fixed

### Problem 1: Case sensitivity — `Method` vs `method`

**Issue:** Flow Launcher sends JSON-RPC with lowercase keys (`"method"`, `"parameters"`), but the Go struct tags said `json:"Method"` (uppercase). The v1 `JsonRPCPlugin` uses `PropertyNameCaseInsensitive = true` on its deserializer, but the plugin's own `Parse()` function was strict.

**Error:**
```
Error: missing 'Method' field
```

**Fix:** Changed all JSON tags to lowercase (`json:"method"`, `json:"parameters"`) and updated the `rawMap` key checks.

**Lesson:** Always match the case that Flow actually sends. For v1 (`ExecutablePlugin`), the request is serialized via `JsonSerializer.Serialize` with `PropertyNamingPolicy = CamelCase`. For v2 (`JsonRPCPluginV2`), `StreamJsonRpc` also uses camelCase. Lowercase is always correct.

---

### Problem 2: Execute handler returned raw string instead of JSON-RPC request model

**Issue:** The `execute_secret` handler returned the raw secret string (e.g., `(mwxA1$HY...`). Flow reads stdout and tries to `JsonSerializer.DeserializeAsync<JsonRPCRequestModel>(...)` on it. The `(` character is not valid JSON, so it threw:

```
System.Text.Json.JsonException: ')' is an invalid start of a value.
```

**Root cause:** The v1 `JsonRPCPlugin.ExecuteResultAsync()` method:
1. Deserializes stdout as `JsonRPCRequestModel` (expecting `{method, parameters}`).
2. If the method starts with `"Flow.Launcher."`, it calls `ExecuteFlowLauncherAPI()`.
3. Otherwise, it sends the request model back to the plugin as a new JSON-RPC call.

So the plugin must return a **JSON-RPC request model**, not a raw string.

**Fix:** Execute handler now returns `flowapi.CopyToClipboard{Text: result}`, which the dispatcher's duck-type check auto-marshals to `{"method":"Flow.Launcher.CopyToClipboard","parameters":[...]}`.

**Lesson:** `execute_*` handlers must return a JSON-RPC request model (method + params), not raw output. The model is either a `Flow.Launcher.*` API call or a custom method that Flow will re-invoke.

---

### Problem 3: Flow v1 query sends `settings` field

**Issue:** Flow's `JsonRPCPlugin.QueryAsync()` sends:
```json
{"Id":0,"Method":"query","Parameters":["secret"],"Settings":{}}
```

The `Settings` field is a `Dictionary<string, object>` passed as the second element of the `Parameters` array. But our `Parse()` function only checks for `parameters` — it doesn't validate the `settings` field. This wasn't actually a bug (the extra field is silently ignored), but it's worth noting for future v2 support where `settings` matters.

**Fix:** None needed for v1 — Flow's `PropertyNameCaseInsensitive` deserializer handles the response regardless. The `settings` field is just ignored during parsing.

---

### Problem 4: `flow` subcommand interferes with root detection

**Issue:** Both `rootCmd` (via `ArbitraryArgs` + JSON detection) and `flowCmd` (a subcommand) can handle JSON-RPC requests. If someone runs `nnry-text-util flow '{"method":"query"}'`, it works via the subcommand. If someone runs `nnry-text-util '{"method":"query"}'` directly, it works via root detection. But having two paths is confusing.

**Fix:** Kept both — the `flow` subcommand is useful for manual testing (`nnry-text-util flow '{"method":"query"}'`), while root detection handles the Flow Launcher invocation automatically.

---

### Problem 5: `string` vs `[]any` type mismatch in `flowapiAction.Parameters`

**Issue:** `handleToolSecret` returns tool args as `[]string`, but `flowapiAction.Parameters` is `[]any`. Go doesn't allow implicit `[]string` to `[]any` conversion.

**Fix:** Explicit loop: `for i, a := range toolArgs { params[i] = a }`.

**Lesson:** Always convert slices explicitly. `[]T` is not `[]any` in Go.

---

## Architecture Decisions

### Why `executable` (v1) over `executable_v2` (v2)

| Factor | v1 | v2 |
|--------|----|----|
| Process model | One-shot per query | Persistent process |
| Input | `os.Args[1]` (JSON string) | stdin pipe (newline-delimited JSON-RPC) |
| Output | stdout (single JSON), exit | stdout pipe, stay alive |
| Protocol | Simple JSON | Full JSON-RPC 2.0 with `jsonrpc`/`id` |
| Complexity | Trivial | Must implement StreamJsonRPC protocol |

v1 is a perfect fit for a Go CLI tool — Go starts instantly, no need to amortize startup cost. v2 can be added later as a second mode if needed.

### Why `pkg/jsonrpc/` is separate from `pkg/flowapi/`

- `jsonrpc` is protocol-level (parsing, dispatching, errors) — no Flow dependency.
- `flowapi` is Flow-specific (IPublicAPI method structs, result types).
- Other apps can use `jsonrpc.Dispatcher` with their own handlers.

### Why duck-type interface instead of importing `flowapi` in `jsonrpc`

`jsonrpc` is a lower-level package. Importing `flowapi` from `jsonrpc` would create a circular dependency if `flowapi` ever needed `jsonrpc`. The duck-type check (`interface{ Method() string; Params() []any }`) achieves the same result without any import.

---

## Future Improvements

### v2 support
- Add `runV2()` mode: read stdin line-by-line, parse JSON-RPC 2.0, dispatch, write JSON-RPC 2.0 response, loop.
- Detect mode: if stdin is a pipe → v2, if `os.Args[1]` is JSON → v1.
- `plugin.json` changes to `"Language": "executable_v2"`.

### Settings support
- v1 passes `settings` as part of the request. Currently ignored.
- Could parse settings from `params[1]` when present and feed them to tool handlers.

### More tools
- Follow the `toolDef` pattern: add `Name`, `Aliases`, `Description`, `Example`, `Handler`.
- Handlers return `(string, error)` — the string is the raw result, and the execute wrapper wraps it in `CopyToClipboard`.

### Plugin icon
- `flow-plugin/Images/app.png` is missing. Generate a 256x256 PNG placeholder.
- Without it, Flow shows a default icon.