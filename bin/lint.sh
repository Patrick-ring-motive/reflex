#!/bin/bash

LINTOUTFILE="gl-code-quality-report.json"
GOLANGCI_LINT_RUNTIME_ARGS="--tests=false --max-issues-per-linter 0 --max-same-issues 0 --timeout 3m"
GOLANGCI_LINT_ENABLED_LINTERS="--enable gosec,revive"
GOLANGCI_LINT_OUTPUT="--output.code-climate.path gl-code-quality-report.json --output.text.path stdout"

echo "go linting..."
go tool golangci-lint run --fix --color always ${GOLANGCI_LINT_RUNTIME_ARGS} ${GOLANGCI_LINT_ENABLED_LINTERS} ${GOLANGCI_LINT_OUTPUT} ./...
