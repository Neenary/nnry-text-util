package jsonrpc

import "encoding/json"

// Error represents a JSON-RPC error object.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Standard JSON-RPC error codes.
var (
	ErrParse          = &Error{Code: -32700, Message: "Parse error"}
	ErrInvalidRequest = &Error{Code: -32600, Message: "Invalid Request"}
	ErrMethodNotFound = &Error{Code: -32601, Message: "Method not found"}
	ErrInvalidParams  = &Error{Code: -32602, Message: "Invalid params"}
	ErrInternal       = &Error{Code: -32603, Message: "Internal error"}
)

// NewError creates a custom error.
func NewError(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

// Error implements the error interface.
func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}