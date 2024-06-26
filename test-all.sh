#!/bin/bash -ex

# go generate ./...
# Find any gofmt errors:
diff -u <(echo -n) <(gofmt -d -s .) || \
    (echo "Please check go formatting!" && exit 1)

# Test the Go code:
go mod tidy
go test -race ./...
go vet ./...
# golangci-lint run - buggy with Go 1.18

# Perform staticcheck:
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...

echo "Done."
