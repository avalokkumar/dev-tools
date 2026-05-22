// Package datax is the Adapter that exposes pkg/datax as Operations.
package datax

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/datax"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{Tool: "data", Op: "json_to_csv",
			Description: "Convert a JSON array of objects to CSV.",
			InputSchema: schemaInputDelim,
			Handler:     handleJSONToCSV,
		},
		{Tool: "data", Op: "csv_to_json",
			Description: "Convert CSV to a JSON array of objects.",
			InputSchema: schemaCSVToJSON,
			Handler:     handleCSVToJSON,
		},
		{Tool: "data", Op: "json_to_xml",
			Description: "Convert JSON to XML.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},"root":{"type":"string","default":"root"},"indent":{"type":"integer","default":2}}
}`),
			Handler: handleJSONToXML,
		},
		{Tool: "data", Op: "xml_to_json",
			Description: "Convert XML to JSON.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},"indent":{"type":"integer","default":2}}
}`),
			Handler: handleXMLToJSON,
		},
		{Tool: "data", Op: "flatten",
			Description: "Flatten nested JSON into a single-level dotted-path map.",
			InputSchema: schemaSep,
			Handler:     handleFlatten,
		},
		{Tool: "data", Op: "unflatten",
			Description: "Inverse of flatten: restore nested JSON from a dotted-path map.",
			InputSchema: schemaSep,
			Handler:     handleUnflatten,
		},
		{Tool: "data", Op: "key_rename",
			Description: "Rename JSON keys per exact-match rules (recursive).",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input","rules"],
  "properties":{"input":{"type":"string"},
                "rules":{"type":"array","items":{"type":"object","required":["from","to"],
                         "properties":{"from":{"type":"string"},"to":{"type":"string"}}}}}
}`),
			Handler: handleKeyRename,
		},
	}
}

var (
	schemaInputDelim = json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},"delimiter":{"type":"string","maxLength":1}}
}`)
	schemaCSVToJSON = json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},"delimiter":{"type":"string","maxLength":1},"hasHeader":{"type":"boolean","default":true}}
}`)
	schemaSep = json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},"sep":{"type":"string","default":"."}}
}`)
)

type inDelim struct {
	Input     string `json:"input"`
	Delimiter string `json:"delimiter,omitempty"`
}

func handleJSONToCSV(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a inDelim
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.JSONToCSV([]byte(a.Input), enginepkg.JSONToCSVOptions{Delimiter: a.Delimiter})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type csvIn struct {
	Input     string `json:"input"`
	Delimiter string `json:"delimiter,omitempty"`
	HasHeader *bool  `json:"hasHeader,omitempty"`
}

func handleCSVToJSON(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a csvIn
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.CSVToJSON([]byte(a.Input), enginepkg.CSVToJSONOptions{Delimiter: a.Delimiter, HasHeader: a.HasHeader})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type xmlOut struct {
	Input  string `json:"input"`
	Root   string `json:"root,omitempty"`
	Indent int    `json:"indent,omitempty"`
}

func handleJSONToXML(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a xmlOut
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.JSONToXML([]byte(a.Input), enginepkg.JSONToXMLOptions{Root: a.Root, Indent: a.Indent})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type xmlIn struct {
	Input  string `json:"input"`
	Indent int    `json:"indent,omitempty"`
}

func handleXMLToJSON(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a xmlIn
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.XMLToJSON([]byte(a.Input), enginepkg.XMLToJSONOptions{Indent: a.Indent})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type sepIn struct {
	Input string `json:"input"`
	Sep   string `json:"sep,omitempty"`
}

func handleFlatten(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a sepIn
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.Flatten([]byte(a.Input), enginepkg.FlattenOptions{Sep: a.Sep})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

func handleUnflatten(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a sepIn
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.Unflatten([]byte(a.Input), enginepkg.FlattenOptions{Sep: a.Sep})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type renameArgs struct {
	Input string                  `json:"input"`
	Rules []enginepkg.KeyRenameRule `json:"rules"`
}

func handleKeyRename(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a renameArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.KeyRename([]byte(a.Input), a.Rules)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
