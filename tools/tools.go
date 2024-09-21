//go:build tools
// +build tools

//go:generate oapi-codegen --config=../oapi.yaml ../api.yaml
//go:generate sqlc generate -f ../sqlc.yaml

package tools

import (
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
)
