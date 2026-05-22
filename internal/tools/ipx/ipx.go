// Package ipx is the Adapter that exposes pkg/ipx as Operations.
package ipx

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/ipx"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{Tool: "ip", Op: "calc",
			Description: "IPv4/IPv6 subnet math from CIDR notation.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["cidr"],
  "properties":{"cidr":{"type":"string"},
                "maxList":{"type":"integer","minimum":0,"maximum":1024,"default":0,
                           "description":"enumerate up to N host addresses"}}
}`),
			Handler: handleCalc,
		},
	}
}

type calcArgs struct {
	CIDR    string `json:"cidr"`
	MaxList int    `json:"maxList,omitempty"`
}

func handleCalc(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a calcArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.CIDRCalc(a.CIDR, enginepkg.Options{MaxList: a.MaxList})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
