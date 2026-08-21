package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nnry-text-util",
	Short: "A CLI tool for text processing utilities",
	Long: `nnry-text-util is a command-line tool for various text processing
operations such as encoding, decoding, formatting, and transforming text.`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// JSON-RPC mode: detect if the first arg is a valid JSON-RPC request
		if len(args) == 1 && looksLikeJSON(args[0]) {
			// First check if it looks like a valid request with a non-empty Method
			var req struct {
				Method string `json:"method"`
			}
			if err := json.Unmarshal([]byte(args[0]), &req); err == nil && req.Method != "" {
				result, err := flowDispatcher.Dispatch(cmd.Context(), args[0])
				if err != nil {
					fmt.Fprintln(cmd.ErrOrStderr(), err.Error())
					return nil
				}
				fmt.Print(result)
				return nil
			}
			// Valid JSON but not a valid JSON-RPC request — show parse errors
			result, err := flowDispatcher.Dispatch(cmd.Context(), args[0])
			if err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), err.Error())
				return nil
			}
			fmt.Print(result)
			return nil
		}
		return cmd.Help()
	},
}

func looksLikeJSON(s string) bool {
	s = strings.TrimSpace(s)
	return strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}