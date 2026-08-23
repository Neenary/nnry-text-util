package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/neenary/nnry-text-util/pkg/flowapi"
	"github.com/neenary/nnry-text-util/pkg/jsonrpc"
	"github.com/spf13/cobra"
)

// flowDispatcher is the JSON-RPC dispatcher for Flow Launcher.
var flowDispatcher = jsonrpc.New()

// toolDef describes a tool that can be invoked via Flow Launcher.
type toolDef struct {
	Name        string
	Aliases     []string
	Description string
	Example     string
	Handler     func(ctx context.Context, args []string) (string, error)
}

var tools = []toolDef{
	{
		Name:        "secret",
		Aliases:     []string{"pw", "password", "token"},
		Description: "Generate a cryptographically secure random string",
		Example:     "secret 64 --alpha --no-upper",
		Handler:     handleToolSecret,
	},
}

var flowCmd = &cobra.Command{
	Use:   "flow",
	Short: "Flow Launcher JSON-RPC plugin mode",
	Long: `Run as a Flow Launcher executable_v1 plugin.

This command is invoked automatically when the binary receives a JSON-RPC request
as the first argument. It can also be called directly for testing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("usage: nnry-text-util flow '<json-rpc-request>'")
		}
		result, err := flowDispatcher.Dispatch(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		fmt.Print(result)
		return nil
	},
}

func init() {
	flowDispatcher.Handle("query", handleFlowQuery)
	flowDispatcher.Handle("context_menu", handleFlowContextMenu)

	// Register execute_* handlers for each tool
	for _, t := range tools {
		tool := t // capture
		flowDispatcher.Handle("execute_"+tool.Name, func(ctx context.Context, req jsonrpc.Request) (any, *jsonrpc.Error) {
			params, _ := jsonrpc.ParseParametersAll[string](req)
			result, err := tool.Handler(ctx, params)
			if err != nil {
				return nil, jsonrpc.NewError(-1, err.Error())
			}
			// Return a typed APIMethod — dispatcher detects Method()/Params() interface
			// and marshals it as {"method":"Flow.Launcher.Xxx","parameters":[...]}
			return flowapi.CopyToClipboard{
				Text:                     result,
				DirectCopy:               false,
				ShowDefaultNotification:  true,
			}, nil
		})
	}

	rootCmd.AddCommand(flowCmd)
}

// handleFlowQuery processes a query request from Flow Launcher.
func handleFlowQuery(ctx context.Context, req jsonrpc.Request) (any, *jsonrpc.Error) {
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

	var results []flowapiResult

	if len(parts) == 0 {
		// No query: show all available tools
		results = listAllTools()
	} else {
		// Try to match a tool
		toolName := parts[0]
		toolArgs := parts[1:]

		tool := findTool(toolName)
		if tool == nil {
			// Show matching tools
			results = filterTools(toolName)
			if len(results) == 0 {
				results = listAllTools()
			}
		} else if len(toolArgs) == 0 {
			// Show tool info
			results = append(results, flowapiResult{
				Title:    fmt.Sprintf("%s — %s", tool.Name, tool.Description),
				SubTitle: fmt.Sprintf("Example: %s", tool.Example),
				Action: &flowapiAction{
					Method:     fmt.Sprintf("execute_%s", tool.Name),
					Parameters: []any{tool.Example},
				},
				ContextData: tool.Name,
			})
		} else {
			// Execute the tool with args (preview)
			result, err := tool.Handler(ctx, toolArgs)
			if err != nil {
				results = append(results, flowapiResult{
					Title:    fmt.Sprintf("Error: %s", err.Error()),
					SubTitle: tool.Description,
					ContextData: tool.Name,
				})
			} else {
				params := make([]any, len(toolArgs))
				for i, a := range toolArgs {
					params[i] = a
				}
				results = append(results, flowapiResult{
					Title:    fmt.Sprintf("%s (%s)", tool.Name, strings.Join(toolArgs, " ")),
					SubTitle: result,
					Action: &flowapiAction{
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

// handleFlowContextMenu processes a context_menu request from Flow Launcher.
func handleFlowContextMenu(ctx context.Context, req jsonrpc.Request) (any, *jsonrpc.Error) {
	params, err := jsonrpc.ParseParametersAll[string](req)
	if err != nil {
		return map[string]any{"result": []any{}}, nil
	}

	contextData := ""
	if len(params) > 0 {
		contextData = params[0]
	}

	results := []flowapiResult{
		{
			Title:    "Copy to clipboard",
			SubTitle: "Copy the result to clipboard",
			Action: &flowapiAction{
				Method:     flowapi.CopyToClipboard{}.Method(),
				Parameters: []any{contextData, false, true},
			},
		},
		{
			Title:    "Generate new",
			SubTitle: "Generate a new result",
			Action: &flowapiAction{
				Method: fmt.Sprintf("execute_%s", contextData),
				Parameters: []any{},
			},
		},
	}

	return map[string]any{"result": results}, nil
}

// handleToolSecret generates a secret string.
func handleToolSecret(ctx context.Context, args []string) (string, error) {
	opts := SecretOptions{
		Length: 32,
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--alpha":
			opts.Alpha = true
		case arg == "--numeric":
			opts.Numeric = true
		case arg == "--alphanumeric" || arg == "-a":
			opts.Alphanumeric = true
		case arg == "--no-symbols":
			opts.NoSymbols = true
		case arg == "--no-digits":
			opts.NoDigits = true
		case arg == "--no-lower":
			opts.NoLower = true
		case arg == "--no-upper":
			opts.NoUpper = true
		case arg == "--length" || arg == "-l":
			if i+1 < len(args) {
				i++
				n, err := strconv.Atoi(args[i])
				if err != nil {
					return "", fmt.Errorf("invalid length '%s': must be a number", args[i])
				}
				opts.Length = n
			}
		default:
			// Positional: first number is length
			n, err := strconv.Atoi(arg)
			if err == nil {
				opts.Length = n
			}
		}
	}

	if opts.Length < 1 {
		return "", fmt.Errorf("length must be at least 1")
	}
	if opts.Length > 10000 {
		return "", fmt.Errorf("length must be at most 10000")
	}

	charset := BuildCharset(opts)
	secret, err := generateSecret(opts.Length, charset)
	if err != nil {
		return "", fmt.Errorf("failed to generate secret: %w", err)
	}
	return secret, nil
}

// flowapiResult is a simplified result struct for Flow Launcher.
type flowapiResult struct {
	Title       string        `json:"Title"`
	SubTitle    string        `json:"SubTitle,omitempty"`
	IcoPath     string        `json:"IcoPath,omitempty"`
	Score       int           `json:"Score,omitempty"`
	Action      *flowapiAction `json:"JsonRPCAction,omitempty"`
	ContextData any           `json:"ContextData,omitempty"`
}

type flowapiAction struct {
	Method               string `json:"method"`
	Parameters           []any  `json:"parameters"`
	DontHideAfterAction  bool   `json:"dontHideAfterAction,omitempty"`
}

// findTool finds a tool by name or alias.
func findTool(name string) *toolDef {
	for _, t := range tools {
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
func listAllTools() []flowapiResult {
	var results []flowapiResult
	for _, t := range tools {
		results = append(results, flowapiResult{
			Title:    t.Name,
			SubTitle: t.Description,
			Action: &flowapiAction{
				Method:     fmt.Sprintf("execute_%s", t.Name),
				Parameters: []any{},
			},
			ContextData: t.Name,
		})
	}
	return results
}

// filterTools returns tools matching the given prefix.
func filterTools(prefix string) []flowapiResult {
	var results []flowapiResult
	for _, t := range tools {
		if strings.HasPrefix(strings.ToLower(t.Name), strings.ToLower(prefix)) {
			results = append(results, flowapiResult{
				Title:    t.Name,
				SubTitle: t.Description,
				Action: &flowapiAction{
					Method:     fmt.Sprintf("execute_%s", t.Name),
					Parameters: []any{},
				},
				ContextData: t.Name,
			})
		}
	}
	return results
}