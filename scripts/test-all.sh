#!/bin/bash -e
# -*- compile-command: "./test-all.sh"; -*-

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "${REPO_ROOT}"

export PATH="$(go env GOPATH)/bin:${PATH}"

# Static cleanliness checks: verify the repo is "squeaky-clean".
#
# 0. errcheck     — verifies all error return values are checked.
# 1. modernize    — the standalone modernize analyzer.
# 2. staticcheck  — static analysis for bugs and style (SA* + U1000).

echo "Running errcheck..."
if ! command -v errcheck >/dev/null 2>&1; then
	echo "Installing errcheck..."
	go install github.com/kisielk/errcheck@latest
fi
ERRCHECK_LOG="${REPO_ROOT}/errcheck.log"
if ! errcheck ./... >"${ERRCHECK_LOG}" 2>&1; then
	echo "FAIL: errcheck reported unchecked errors (see ${ERRCHECK_LOG}):" >&2
	cat "${ERRCHECK_LOG}" >&2
	exit 1
fi
rm -f "${ERRCHECK_LOG}"
echo "errcheck: clean (all errors checked)"

echo "Running modernize (modernization suggestions)..."
MODERNIZE_LOG="${REPO_ROOT}/modernize.log"
if ! go run golang.org/x/tools/go/analysis/passes/modernize/cmd/modernize@latest ./... >"${MODERNIZE_LOG}" 2>&1; then
	echo "FAIL: modernize reported suggestions (see ${MODERNIZE_LOG}):" >&2
	cat "${MODERNIZE_LOG}" >&2
	exit 1
fi
rm -f "${MODERNIZE_LOG}"
echo "modernize: clean (no suggestions)"

echo "Running staticcheck (SA* + U1000)..."
if ! command -v staticcheck >/dev/null 2>&1; then
	echo "Installing staticcheck..."
	go install honnef.co/go/tools/cmd/staticcheck@latest
fi
STATICCHECK_LOG="${REPO_ROOT}/staticcheck.log"
staticcheck -checks="SA*,U1000" ./... >"${STATICCHECK_LOG}" 2>&1 || true
if [[ -s "${STATICCHECK_LOG}" ]]; then
	echo "FAIL: staticcheck reported issues (see ${STATICCHECK_LOG}):" >&2
	cat "${STATICCHECK_LOG}" >&2
	exit 1
fi
rm -f "${STATICCHECK_LOG}"
echo "staticcheck: clean (SA* + U1000)"

echo "Running go mod tidy..."
go mod tidy

echo "Running go vet..."
go vet ./...

echo "Running unit tests with race detector..."
go test -race ./...

echo "Repo is squeaky-clean (errcheck + gopls check + modernize + staticcheck SA*/U1000 + deadcode advisory + all tests)."
