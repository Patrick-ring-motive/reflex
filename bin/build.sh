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

coverage_file="$(mktemp "${TMPDIR:-/tmp}/reflex-coverage.XXXXXX")"
trap 'rm -f "$coverage_file"' EXIT

go test -race -covermode=atomic -coverprofile="$coverage_file" ./...
coverage="$(go tool cover -func="$coverage_file" | awk '/^total:/ { print $3 }')"
printf 'Statement coverage: %s\n' "$coverage"
if [[ "$coverage" != "100.0%" ]]; then
	printf 'Statement coverage must be 100.0%%.\n' >&2
	exit 1
fi

printf 'Build checks passed.\n'