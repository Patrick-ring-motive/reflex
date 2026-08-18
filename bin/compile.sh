#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
repo_dir="$(cd -- "$script_dir/.." >/dev/null 2>&1 && pwd)"
cd "$repo_dir"

go build ./...
go test -run '^$' ./...
