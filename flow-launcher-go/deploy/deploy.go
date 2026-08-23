// Package deploy provides helpers to build and deploy Flow Launcher plugins.
package deploy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Config defines the deployment configuration.
type Config struct {
	// BinaryName is the name of the output binary (e.g., "my-plugin.exe").
	BinaryName string
	// SourceDir is the directory containing main.go.
	SourceDir string
	// PluginDir is the target Flow Launcher plugin directory.
	// If empty, the default Flow Launcher plugins directory is used.
	PluginDir string
	// ExtraFiles are additional files to copy (e.g., plugin.json, Images/).
	ExtraFiles []string
}

// Build compiles the plugin binary.
func Build(cfg Config) error {
	output := filepath.Join(cfg.SourceDir, cfg.BinaryName)
	cmd := exec.Command("go", "build", "-o", output, ".")
	cmd.Dir = cfg.SourceDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ToFlow deploys the plugin to the Flow Launcher plugins directory.
// It builds the binary and copies all necessary files.
func ToFlow(cfg Config) error {
	// Determine target directory
	pluginDir := cfg.PluginDir
	if pluginDir == "" {
		var err error
		pluginDir, err = defaultFlowPluginDir()
		if err != nil {
			return fmt.Errorf("deploy: cannot determine Flow Launcher plugin dir: %w", err)
		}
	}

	// Build the binary
	if err := Build(cfg); err != nil {
		return fmt.Errorf("deploy: build failed: %w", err)
	}

	// Create target directory
	targetDir := filepath.Join(pluginDir, filepath.Base(cfg.SourceDir))
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("deploy: cannot create target dir: %w", err)
	}

	// Copy binary
	src := filepath.Join(cfg.SourceDir, cfg.BinaryName)
	dst := filepath.Join(targetDir, cfg.BinaryName)
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("deploy: cannot read binary: %w", err)
	}
	if err := os.WriteFile(dst, data, 0755); err != nil {
		return fmt.Errorf("deploy: cannot write binary: %w", err)
	}

	// Copy extra files
	for _, f := range cfg.ExtraFiles {
		src := filepath.Join(cfg.SourceDir, f)
		dst := filepath.Join(targetDir, f)
		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("deploy: cannot copy %s: %w", f, err)
		}
	}

	fmt.Printf("Deployed to %s\n", targetDir)
	return nil
}

func defaultFlowPluginDir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA not set")
		}
		return filepath.Join(appData, "FlowLauncher", "Plugins"), nil
	default:
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".config", "FlowLauncher", "Plugins"), nil
	}
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}