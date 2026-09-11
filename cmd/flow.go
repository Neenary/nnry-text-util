package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	flowPlugin.Define(plugin.ToolDef{
		Name:        "timestamp",
		Aliases:     []string{"ts", "time", "now", "date"},
		Description: "Generate an ISO timestamp string valid for Windows filenames",
		Example:     "timestamp --utc --compact",
		Handler:     handleToolTimestamp,
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

// handleToolTimestamp generates a Windows-safe timestamp string.
func handleToolTimestamp(ctx context.Context, args []string) (string, error) {
	utc := false
	dateOnly := false
	timeOnly := false
	compact := false
	separator := "-"

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--utc":
			utc = true
		case arg == "--date" || arg == "-d":
			dateOnly = true
		case arg == "--time" || arg == "-t":
			timeOnly = true
		case arg == "--compact" || arg == "-c":
			compact = true
		case arg == "--separator" || arg == "-s":
			if i+1 < len(args) {
				i++
				separator = args[i]
			}
		}
	}

	now := time.Now()
	if utc {
		now = now.UTC()
	}

	var result string
	switch {
	case compact:
		if dateOnly {
			result = now.Format("20060102")
		} else if timeOnly {
			result = now.Format("150405")
		} else {
			result = now.Format("20060102T150405")
		}
	case dateOnly:
		result = now.Format("2006-01-02")
	case timeOnly:
		result = now.Format("15:04:05")
	default:
		result = now.Format("2006-01-02T15:04:05")
	}

	if !compact && !dateOnly && !timeOnly {
		result = strings.ReplaceAll(result, ":", separator)
	} else if timeOnly && !compact {
		result = strings.ReplaceAll(result, ":", separator)
	}

	return result, nil
}