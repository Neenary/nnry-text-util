package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/neenary/flow-launcher-go/plugin"
)

func init() {
	flowPlugin.Define(plugin.ToolDef{
		Name:        "secret",
		Aliases:     []string{"pw", "password", "token"},
		Description: "Generate a cryptographically secure random string",
		Example:     "secret 64 --alpha --no-upper",
		Handler:     handleToolSecret,
	})
	flowPlugin.Attach(rootCmd)
}

var flowPlugin = plugin.New()

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