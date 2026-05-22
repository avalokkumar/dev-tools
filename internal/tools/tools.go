// Package tools wires built-in Engines into a single Registry consumed by
// every Surface (CLI, Web, MCP). Adapters live in subpackages.
//
// This package is the only place where Surfaces learn the catalog of
// Operations. Adding a new Engine = add one RegisterX call below.
package tools

import (
	"github.com/devforge/devforge/internal/mcpserver"
	"github.com/devforge/devforge/internal/tools/codefmt"
	"github.com/devforge/devforge/internal/tools/colorx"
	"github.com/devforge/devforge/internal/tools/cronx"
	"github.com/devforge/devforge/internal/tools/cryptox"
	"github.com/devforge/devforge/internal/tools/csvfmt"
	"github.com/devforge/devforge/internal/tools/datax"
	"github.com/devforge/devforge/internal/tools/devx"
	"github.com/devforge/devforge/internal/tools/enc"
	"github.com/devforge/devforge/internal/tools/faker"
	"github.com/devforge/devforge/internal/tools/gitx"
	"github.com/devforge/devforge/internal/tools/httpx"
	"github.com/devforge/devforge/internal/tools/idx"
	"github.com/devforge/devforge/internal/tools/ipx"
	"github.com/devforge/devforge/internal/tools/jsonfmt"
	"github.com/devforge/devforge/internal/tools/jwtx"
	toolsmathx "github.com/devforge/devforge/internal/tools/mathx"
	"github.com/devforge/devforge/internal/tools/mdx"
	"github.com/devforge/devforge/internal/tools/netx"
	"github.com/devforge/devforge/internal/tools/regextool"
	"github.com/devforge/devforge/internal/tools/smartdiff"
	"github.com/devforge/devforge/internal/tools/sqlfmt"
	"github.com/devforge/devforge/internal/tools/strx"
	"github.com/devforge/devforge/internal/tools/timex"
	"github.com/devforge/devforge/internal/tools/totpx"
	"github.com/devforge/devforge/internal/tools/tzconv"
	"github.com/devforge/devforge/internal/tools/uuid"
	"github.com/devforge/devforge/internal/tools/yamlfmt"
)

// Register adds every built-in Operation to reg.
// Returns the first registration error encountered.
func Register(reg *mcpserver.Registry) error {
	groups := [][]mcpserver.Operation{
		// Phase A-D originals
		uuid.Operations(),
		jsonfmt.Operations(),
		yamlfmt.Operations(),
		csvfmt.Operations(),
		smartdiff.Operations(),
		regextool.Operations(),
		cronx.Operations(),
		jwtx.Operations(),
		tzconv.Operations(),
		faker.Operations(),
		// Phase E1-E3
		enc.Operations(),
		strx.Operations(),
		timex.Operations(),
		sqlfmt.Operations(),
		mdx.Operations(),
		idx.Operations(),
		colorx.Operations(),
		// Phase E4
		datax.Operations(),
		// Phase E5
		gitx.Operations(),
		devx.Operations(),
		// Phase E6 + E7
		netx.Operations(),
		httpx.Operations(),
		// Phase E8
		cryptox.Operations(),
		totpx.Operations(),
		// Phase E9
		codefmt.Operations(),
		toolsmathx.Operations(),
		ipx.Operations(),
	}
	for _, ops := range groups {
		for _, op := range ops {
			if err := reg.Register(op); err != nil {
				return err
			}
		}
	}
	return nil
}
