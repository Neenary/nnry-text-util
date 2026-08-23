package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nnry-text-util",
	Short: "A CLI tool for text processing utilities",
	Long: `nnry-text-util is a command-line tool for various text processing
operations such as encoding, decoding, formatting, and transforming text.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}