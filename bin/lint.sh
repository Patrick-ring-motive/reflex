#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
repo_dir="$(cd -- "$script_dir/.." >/dev/null 2>&1 && pwd)"
cd "$repo_dir"

unformatted="$(find . -name '*.go' -not -path './vendor/*' -exec gofmt -l {} +)"
if [[ -n "$unformatted" ]]; then
	printf 'Run gofmt on:\n%s\n' "$unformatted" >&2
	exit 1
fi

go vet ./...
