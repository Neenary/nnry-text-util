package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
)

// HandlerFunc processes a JSON-RPC request and returns a result or error.
type HandlerFunc func(ctx context.Context, req Request) (any, *Error)

// Dispatcher routes JSON-RPC requests to registered handlers.
type Dispatcher struct {
	handlers map[string]HandlerFunc
}

// New creates a new Dispatcher.
func New() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string]HandlerFunc),
	}
}

// Handle registers a handler for a method.
func (d *Dispatcher) Handle(method string, fn HandlerFunc) {
	d.handlers[method] = fn
}

// Dispatch parses and routes a raw JSON-RPC request string.
// Returns the JSON response string.
func (d *Dispatcher) Dispatch(ctx context.Context, raw string) (string, error) {
	req, err := Parse(raw)
	if err != nil {
		return "", err
	}

	// Look up handler
	fn, ok := d.handlers[req.Method]
	if !ok {
		return "", fmt.Errorf("Error: unknown method '%s'", req.Method)
	}

	result, rpcErr := fn(ctx, req)
	if rpcErr != nil {
		resp := NewErrorResponse(req.ID, rpcErr)
		b, _ := json.Marshal(resp)
		return string(b), nil
	}

	// If result is a string, return it directly (raw output)
	if s, ok := result.(string); ok {
		return s, nil
	}

	// If result implements Method()/Params() (e.g. api.APIMethod),
	// marshal as JSON-RPC request model: {"method":"...","parameters":[...]}
	if rm, ok := result.(interface{ Method() string; Params() []any }); ok {
		b, err := json.Marshal(map[string]any{
			"method":     rm.Method(),
			"parameters": rm.Params(),
		})
		if err != nil {
			return "", fmt.Errorf("Error: failed to marshal response: %w", err)
		}
		return string(b), nil
	}

	// If result is a map or slice, marshal it directly (v1 JSON-RPC format)
	// The caller already structured the response (e.g., {"result": [...]})
	b, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("Error: failed to marshal response: %w", err)
	}
	return string(b), nil
}

// ParseParameters unmarshals the Parameters field into the given target.
func ParseParameters[T any](req Request) (T, error) {
	var zero T
	if req.Parameters == nil {
		return zero, nil
	}
	var params []T
	if err := json.Unmarshal(req.Parameters, &params); err != nil {
		return zero, fmt.Errorf("invalid parameters: %w", err)
	}
	if len(params) == 0 {
		return zero, nil
	}
	return params[0], nil
}

// ParseParametersAll unmarshals all Parameters into a slice.
func ParseParametersAll[T any](req Request) ([]T, error) {
	if req.Parameters == nil {
		return nil, nil
	}
	var params []T
	if err := json.Unmarshal(req.Parameters, &params); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}
	return params, nil
}