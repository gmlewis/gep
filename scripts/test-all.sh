#!/bin/bash -e
# -*- compile-command: "./test-all.sh"; -*-

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
cd ${SCRIPT_DIR}/..

go mod tidy
go test -race ./...
go vet ./...

go install honnef.co/go/tools/cmd/staticcheck@2025.1
staticcheck ./...

echo "Done."
