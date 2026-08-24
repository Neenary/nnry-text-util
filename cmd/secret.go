package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/spf13/cobra"
)

const (
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits    = "0123456789"
	symbols   = "!@#$%^&*()-_=+[]{}|;:,.<>?/~"
)

// SecretOptions defines the character set options for secret generation.
type SecretOptions struct {
	Length       int
	Alpha        bool // only letters (no digits, no symbols)
	Numeric      bool // only digits
	Alphanumeric bool // letters + digits, no symbols
	NoSymbols    bool // exclude symbols
	NoDigits     bool // exclude digits
	NoLower      bool
	NoUpper      bool
}

// BuildCharset builds the character set string based on the options.
func BuildCharset(opts SecretOptions) string {
	var charset strings.Builder

	if opts.Alpha {
		if !opts.NoLower {
			charset.WriteString(lowercase)
		}
		if !opts.NoUpper {
			charset.WriteString(uppercase)
		}
		return charset.String()
	}

	if opts.Numeric {
		return digits
	}

	if !opts.NoLower {
		charset.WriteString(lowercase)
	}
	if !opts.NoUpper {
		charset.WriteString(uppercase)
	}
	if !opts.NoDigits {
		charset.WriteString(digits)
	}
	if !opts.Alphanumeric && !opts.NoSymbols {
		charset.WriteString(symbols)
	}

	if charset.Len() == 0 {
		charset.WriteString(lowercase + uppercase + digits + symbols)
	}

	return charset.String()
}

func generateSecret(length int, charset string) (string, error) {
	if charset == "" {
		return "", fmt.Errorf("charset is empty")
	}
	var result strings.Builder
	result.Grow(length)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result.WriteByte(charset[n.Int64()])
	}

	return result.String(), nil
}

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Generate a cryptographically secure random secret string",
	Long: `Generate a cryptographically secure random secret string with customizable character sets.

By default, generates a 32-character secret using letters, digits, and symbols.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		length, _ := cmd.Flags().GetInt("length")
		alpha, _ := cmd.Flags().GetBool("alpha")
		numeric, _ := cmd.Flags().GetBool("numeric")
		alphanumeric, _ := cmd.Flags().GetBool("alphanumeric")
		noDigits, _ := cmd.Flags().GetBool("no-digits")
		noLower, _ := cmd.Flags().GetBool("no-lower")
		noUpper, _ := cmd.Flags().GetBool("no-upper")
		noSymbols, _ := cmd.Flags().GetBool("no-symbols")

		opts := SecretOptions{
			Length:       length,
			Alpha:        alpha,
			Numeric:      numeric,
			Alphanumeric: alphanumeric,
			NoSymbols:    noSymbols,
			NoDigits:     noDigits,
			NoLower:      noLower,
			NoUpper:      noUpper,
		}

		charset := BuildCharset(opts)

		result, err := generateSecret(length, charset)
		if err != nil {
			return err
		}

		fmt.Println(result)
			return nil
	},
}

func init() {
	secretCmd.Flags().IntP("length", "l", 32, "Length of the secret string")
	secretCmd.Flags().Bool("alpha", false, "Only letters (no digits, no symbols)")
	secretCmd.Flags().Bool("numeric", false, "Only digits")
	secretCmd.Flags().BoolP("alphanumeric", "a", false, "Only letters and digits (no symbols)")
	secretCmd.Flags().Bool("no-digits", false, "Exclude digits")
	secretCmd.Flags().Bool("no-lower", false, "Exclude lowercase letters")
	secretCmd.Flags().Bool("no-upper", false, "Exclude uppercase letters")
	secretCmd.Flags().Bool("no-symbols", false, "Exclude symbols (default: symbols included)")

	rootCmd.AddCommand(secretCmd)
}