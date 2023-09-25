//go:build tools

package tools

import (
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "go.uber.org/mock/mockgen"
	_ "golang.org/x/vuln/cmd/govulncheck"
)
