#!/usr/bin/env bash
set -euo pipefail

output="${1:-nnry-text-util}"
os="${GOOS:-}"
arch="${GOARCH:-}"

if [[ "$os" == "windows" ]] || [[ -z "$os" && "$(go env GOOS)" == "windows" ]]; then
	output="${output}.exe"
fi

go_build=(go build -o "$output")
if [[ -n "$os" ]]; then
	go_build+=( -ldflags="-s -w" )
fi

GOOS="$os" GOARCH="$arch" "${go_build[@]}" .