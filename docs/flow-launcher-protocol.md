# Flow Launcher JSON-RPC Protocol

## Plugin Lifecycle

Flow Launcher communicates with the plugin via stdin/stdout using JSON-RPC 2.0.

### 1. Query

User types in the search bar → Flow Launcher sends a `query` request:

```json
{"method": "query", "params": ["<search text>"], "id": 1}
```

Plugin responds with a list of `Result` objects:

```json
{"result": [{"Title": "...", "SubTitle": "...", "JsonRPCAction": {"method": "execute_<tool>", "parameters": [...]}, "ContextData": "..."}]}
```

### 2. Execute

User presses Enter on a result → Flow Launcher sends an `execute_<tool>` request:

```json
{"method": "execute_<tool>", "params": ["<arg1>", "<arg2>", ...], "id": 2}
```

Plugin responds with the **next action** as a `JsonRPCRequestModel`:

```json
{"method": "Flow.Launcher.CopyToClipboard", "parameters": ["<result>", true, true]}
```

Key properties of `JsonRPCRequestModel`:

| Field        | Type   | Description |
|-------------|--------|-------------|
| `method`     | string | Action to execute (built-in `Flow.Launcher.*` or custom) |
| `parameters` | array  | Parameters for the action |
| `id`         | int    | Request ID (optional) |

### 3. Context Menu

User right-clicks a result → Flow Launcher sends a `context_menu` request:

```json
{"method": "context_menu", "params": ["<context_data>"], "id": 3}
```

Plugin responds with context menu items (same `Result` format as query).

---

## Built-in Flow Launcher Methods

| Method | Parameters | Effect |
|--------|-----------|--------|
| `Flow.Launcher.CopyToClipboard` | `[text, showNotification, alsoCopyToAll]` | Copies text to clipboard |
| `Flow.Launcher.OpenUrl` | `[url]` | Opens URL in default browser |
| `Flow.Launcher.ShowMsg` | `[title, subtitle, icoPath]` | Shows a message box |
| `Flow.Launcher.ShellRun` | `[command]` | Runs a shell command |

---

## Response Format Rules

- If response is **empty** (no stdout output), Flow Launcher hides the panel.
- If response has a `method` starting with `Flow.Launcher.`, it's handled internally.
- If response has a custom `method`, Flow Launcher sends another JSON-RPC request to the plugin.

---

## Key gotchas

- **Never return a raw string** from an execute handler — Flow Launcher tries to parse it as JSON and crashes with `JsonException`.
- Always return a JSON object with `method` and `parameters` from execute handlers.
- To copy the result and close, use `Flow.Launcher.CopyToClipboard`.
- The `query` response is `{"result": [...]}`, NOT a full JSON-RPC envelope.