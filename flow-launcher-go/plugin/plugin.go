// Package plugin provides a framework for building Flow Launcher plugins in Go.
// It handles JSON-RPC dispatch, tool registration, query/context_menu handling,
// and cobra integration.
package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/neenary/jsonrpc"
	"github.com/spf13/cobra"
)

// ToolDef describes a tool that can be invoked via Flow Launcher.
type ToolDef struct {
	Name        string
	Aliases     []string
	Description string
	Example     string
	Handler     func(ctx context.Context, args []string) (string, error)
}

// JsonRPCAction is the action to execute when a user selects a result.
type JsonRPCAction struct {
	Method               string `json:"method"`
	Parameters           []any  `json:"parameters"`
	DontHideAfterAction  bool   `json:"dontHideAfterAction,omitempty"`
}

// Result is a single result item in a Flow Launcher query response.
type Result struct {
	Title       string         `json:"Title"`
	SubTitle    string         `json:"SubTitle,omitempty"`
	IcoPath     string         `json:"IcoPath,omitempty"`
	Score       int            `json:"Score,omitempty"`
	Action      *JsonRPCAction `json:"JsonRPCAction,omitempty"`
	ContextData any            `json:"ContextData,omitempty"`
}

// NewResult creates a Result whose JsonRPCAction calls the given method with params.
func NewResult(title, subtitle, method string, params []any) *Result {
	return &Result{
		Title:    title,
		SubTitle: subtitle,
		Action: &JsonRPCAction{
			Method:     method,
			Parameters: params,
		},
	}
}

// NewResultWithIcoPath creates a Result with an icon path.
func NewResultWithIcoPath(title, subtitle, icoPath, method string, params []any) *Result {
	r := NewResult(title, subtitle, method, params)
	r.IcoPath = icoPath
	return r
}

// Plugin is the Flow Launcher plugin runtime.
type Plugin struct {
	dispatcher *jsonrpc.Dispatcher
	tools      []ToolDef
}

// New creates a new Plugin.
func New() *Plugin {
	p := &Plugin{
		dispatcher: jsonrpc.New(),
	}
	p.registerBuiltinHandlers()
	return p
}

// Define registers a tool with the plugin.
func (p *Plugin) Define(tool ToolDef) {
	p.tools = append(p.tools, tool)
	p.dispatcher.Handle("execute_"+tool.Name, p.makeExecuteHandler(tool))
}

// Attach wires the plugin into a cobra.Command.
// It sets Args to cobra.ArbitraryArgs and wraps RunE to detect JSON-RPC requests.
func (p *Plugin) Attach(cmd *cobra.Command) {
	cmd.Args = cobra.ArbitraryArgs
	existingRunE := cmd.RunE
	cmd.RunE = func(c *cobra.Command, args []string) error {
		if len(args) == 1 && looksLikeJSON(args[0]) {
			var req struct {
				Method string `json:"method"`
			}
			if err := json.Unmarshal([]byte(args[0]), &req); err == nil && req.Method != "" {
				result, err := p.dispatcher.Dispatch(c.Context(), args[0])
				if err != nil {
					fmt.Fprintln(c.ErrOrStderr(), err.Error())
					return nil
				}
				fmt.Print(result)
				return nil
			}
			// Valid JSON but not a valid JSON-RPC request — show parse errors
			result, err := p.dispatcher.Dispatch(c.Context(), args[0])
			if err != nil {
				fmt.Fprintln(c.ErrOrStderr(), err.Error())
				return nil
			}
			fmt.Print(result)
			return nil
		}
		if existingRunE != nil {
			return existingRunE(c, args)
		}
		return cmd.Help()
	}
}

// Dispatcher returns the underlying JSON-RPC dispatcher for advanced use.
func (p *Plugin) Dispatcher() *jsonrpc.Dispatcher {
	return p.dispatcher
}

// Tools returns the registered tool definitions.
func (p *Plugin) Tools() []ToolDef {
	return p.tools
}

func (p *Plugin) registerBuiltinHandlers() {
	p.dispatcher.Handle("query", p.handleQuery)
	p.dispatcher.Handle("context_menu", p.handleContextMenu)
}

func (p *Plugin) makeExecuteHandler(tool ToolDef) jsonrpc.HandlerFunc {
	return func(ctx context.Context, req jsonrpc.Request) (any, *jsonrpc.Error) {
		params, _ := jsonrpc.ParseParametersAll[string](req)
		result, err := tool.Handler(ctx, params)
		if err != nil {
			return nil, jsonrpc.NewError(-1, err.Error())
		}
		// Return a string so dispatcher returns it as raw output
		return result, nil
	}
}

// handleQuery processes a query request from Flow Launcher.
func (p *Plugin) handleQuery(ctx context.Context, req jsonrpc.Request) (any, *jsonrpc.Error) {
	params, err := jsonrpc.ParseParametersAll[string](req)
	if err != nil {
		params = []string{""}
	}

	query := ""
	if len(params) > 0 {
		query = params[0]
	}

	query = strings.TrimSpace(query)
	parts := strings.Fields(query)

	var results []Result

	if len(parts) == 0 {
		results = p.listAllTools()
	} else {
		toolName := parts[0]
		toolArgs := parts[1:]

		tool := p.findTool(toolName)
		if tool == nil {
			results = p.filterTools(toolName)
			if len(results) == 0 {
				results = p.listAllTools()
			}
		} else if len(toolArgs) == 0 {
			results = append(results, Result{
				Title:    fmt.Sprintf("%s — %s", tool.Name, tool.Description),
				SubTitle: fmt.Sprintf("Example: %s", tool.Example),
				Action: &JsonRPCAction{
					Method:     fmt.Sprintf("execute_%s", tool.Name),
					Parameters: []any{tool.Example},
				},
				ContextData: tool.Name,
			})
		} else {
			result, err := tool.Handler(ctx, toolArgs)
			if err != nil {
				results = append(results, Result{
					Title:       fmt.Sprintf("Error: %s", err.Error()),
					SubTitle:    tool.Description,
					ContextData: tool.Name,
				})
			} else {
				params := make([]any, len(toolArgs))
				for i, a := range toolArgs {
					params[i] = a
				}
				results = append(results, Result{
					Title:    fmt.Sprintf("%s (%s)", tool.Name, strings.Join(toolArgs, " ")),
					SubTitle: result,
					Action: &JsonRPCAction{
						Method:     fmt.Sprintf("execute_%s", tool.Name),
						Parameters: params,
					},
					ContextData: tool.Name,
				})
			}
		}
	}

	return map[string]any{"result": results}, nil
}

// handleContextMenu processes a context_menu request from Flow Launcher.
func (p *Plugin) handleContextMenu(ctx context.Context, req jsonrpc.Request) (any, *jsonrpc.Error) {
	params, err := jsonrpc.ParseParametersAll[string](req)
	if err != nil {
		return map[string]any{"result": []any{}}, nil
	}

	contextData := ""
	if len(params) > 0 {
		contextData = params[0]
	}

	results := []Result{
		{
			Title:    "Copy to clipboard",
			SubTitle: "Copy the result to clipboard",
			Action: &JsonRPCAction{
				Method:     "Flow.Launcher.CopyToClipboard",
				Parameters: []any{contextData, false, true},
			},
		},
		{
			Title:    "Generate new",
			SubTitle: "Generate a new result",
			Action: &JsonRPCAction{
				Method:     fmt.Sprintf("execute_%s", contextData),
				Parameters: []any{},
			},
		},
	}

	return map[string]any{"result": results}, nil
}

// findTool finds a tool by name or alias.
func (p *Plugin) findTool(name string) *ToolDef {
	for _, t := range p.tools {
		if strings.EqualFold(t.Name, name) {
			return &t
		}
		for _, a := range t.Aliases {
			if strings.EqualFold(a, name) {
				return &t
			}
		}
	}
	return nil
}

// listAllTools returns results for all available tools.
func (p *Plugin) listAllTools() []Result {
	var results []Result
	for _, t := range p.tools {
		results = append(results, Result{
			Title:    t.Name,
			SubTitle: t.Description,
			Action: &JsonRPCAction{
				Method:     fmt.Sprintf("execute_%s", t.Name),
				Parameters: []any{},
			},
			ContextData: t.Name,
		})
	}
	return results
}

// filterTools returns tools matching the given prefix.
func (p *Plugin) filterTools(prefix string) []Result {
	var results []Result
	for _, t := range p.tools {
		if strings.HasPrefix(strings.ToLower(t.Name), strings.ToLower(prefix)) {
			results = append(results, Result{
				Title:    t.Name,
				SubTitle: t.Description,
				Action: &JsonRPCAction{
					Method:     fmt.Sprintf("execute_%s", t.Name),
					Parameters: []any{},
				},
				ContextData: t.Name,
			})
		}
	}
	return results
}

func looksLikeJSON(s string) bool {
	s = strings.TrimSpace(s)
	return strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")
}