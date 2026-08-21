package jsonrpc

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Request represents a JSON-RPC request.
type Request struct {
	JSONRPC    string          `json:"jsonrpc,omitempty"`
	Method     string          `json:"method"`
	Parameters json.RawMessage `json:"parameters,omitempty"`
	ID         any             `json:"id,omitempty"`
}

// Parse validates and parses a raw JSON-RPC string.
// It returns detailed error messages for invalid input.
func Parse(raw string) (Request, error) {
	var req Request

	// 1. Valid JSON
	if err := json.Unmarshal([]byte(raw), &req); err != nil {
		return req, fmt.Errorf("Error: not a valid JSON string")
	}

	// 2. Has Method
	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal([]byte(raw), &rawMap); err != nil {
		return req, fmt.Errorf("Error: not a valid JSON string")
	}

	if _, ok := rawMap["method"]; !ok {
		return req, fmt.Errorf("Error: missing 'method' field")
	}

	// 3. Method not empty
	if strings.TrimSpace(req.Method) == "" {
		return req, fmt.Errorf("Error: 'method' cannot be empty")
	}

	// 4. Parameters is an array (if present)
	if paramsRaw, ok := rawMap["parameters"]; ok {
		var params []any
		if err := json.Unmarshal(paramsRaw, &params); err != nil {
			return req, fmt.Errorf("Error: 'parameters' must be a JSON array")
		}
	}

	return req, nil
}