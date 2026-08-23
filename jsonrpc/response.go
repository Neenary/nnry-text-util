package jsonrpc

import "encoding/json"

// Response represents a JSON-RPC response.
type Response struct {
	JSONRPC string `json:"jsonrpc,omitempty"`
	Result  any    `json:"result,omitempty"`
	Error   *Error `json:"error,omitempty"`
	ID      any    `json:"id"`
}

// NewResponse creates a success response.
func NewResponse(id any, result any) Response {
	return Response{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	}
}

// NewErrorResponse creates an error response.
func NewErrorResponse(id any, err *Error) Response {
	return Response{
		JSONRPC: "2.0",
		Error:   err,
		ID:      id,
	}
}

// MarshalJSON returns the JSON encoding of the response.
func (r Response) MarshalJSON() ([]byte, error) {
	// Use a type alias to avoid infinite recursion
	type Alias Response
	return json.Marshal(Alias(r))
}