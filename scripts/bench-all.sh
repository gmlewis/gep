#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}/.."

# Keep the repo-wide benchmark sweep practical for local and CI use.
go test -run='^$' -bench=. -benchmem -count=1 -benchtime=1x ./...
