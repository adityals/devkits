// Package gots converts Go type declarations to TypeScript.
package gots

import (
	"fmt"
	"go/parser"
	"go/token"

	"github.com/gzuidhof/tygo/tygo"
)

type Converter struct {
	cfg tygo.PackageConfig
}

func New() *Converter {
	return &Converter{
		cfg: tygo.PackageConfig{
			// time.Time marshals to an RFC 3339 string.
			TypeMappings: map[string]string{
				"time.Time": "string",
			},
		},
	}
}

func (c *Converter) Convert(src string) (string, error) {
	out, err := tygo.ConvertGoToTypescript(stripPackageClause(src), c.cfg)
	if err != nil {
		return "", fmt.Errorf("convert go to ts: %w", err)
	}

	return out, nil
}

// stripPackageClause drops a leading `package x` clause, because tygo prepends its own.
func stripPackageClause(src string) string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.PackageClauseOnly)
	if err != nil {
		return src
	}

	return src[fset.Position(f.Name.End()).Offset:]
}
