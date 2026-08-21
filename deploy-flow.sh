#!/usr/bin/env bash
set -euo pipefail

# Build binary
go build -o flow-plugin/nnry-text-util.exe .

# Flow Launcher plugin directory
flow_plugins="${APPDATA}/FlowLauncher/Plugins"
target="${flow_plugins}/nnry-text-util"

mkdir -p "$target"

# Copy plugin files
cp -r flow-plugin/* "$target"

echo "Deployed to $target"