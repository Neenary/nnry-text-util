package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var timestampCmd = &cobra.Command{
	Use:   "timestamp",
	Short: "Generate an ISO timestamp string valid for Windows filenames",
	Long: `Generate an ISO 8601-like timestamp with colons replaced by hyphens
so the result is usable as a Windows filename.

Default format: YYYY-MM-DDTHH-MM-SS (local time).
Use --utc for UTC, --date for date only, --time for time only,
or --compact for a condensed format without separators.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		utc, _ := cmd.Flags().GetBool("utc")
		dateOnly, _ := cmd.Flags().GetBool("date")
		timeOnly, _ := cmd.Flags().GetBool("time")
		compact, _ := cmd.Flags().GetBool("compact")
		separator, _ := cmd.Flags().GetString("separator")

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
			// Replace colons with the chosen separator (only for full timestamp)
			result = strings.ReplaceAll(result, ":", separator)
		} else if timeOnly && !compact {
			result = strings.ReplaceAll(result, ":", separator)
		}

		fmt.Println(result)
		return nil
	},
}

func init() {
	timestampCmd.Flags().Bool("utc", false, "Use UTC instead of local time")
	timestampCmd.Flags().BoolP("date", "d", false, "Date only (YYYY-MM-DD)")
	timestampCmd.Flags().BoolP("time", "t", false, "Time only (HH-MM-SS)")
	timestampCmd.Flags().BoolP("compact", "c", false, "Compact format (YYYYMMDDTHHMMSS)")
	timestampCmd.Flags().StringP("separator", "s", "-", "Separator to replace colons")

	rootCmd.AddCommand(timestampCmd)
}