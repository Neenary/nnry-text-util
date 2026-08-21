// Package flowapi provides typed Go structs for the Flow Launcher JSON-RPC API.
// Each struct represents a callable method on the IPublicAPI interface.
package flowapi

import "encoding/json"

// APIMethod is implemented by all Flow Launcher API call structs.
type APIMethod interface {
	// Method returns the Flow.Launcher.* method name.
	Method() string
	// Params returns the positional parameters for the JSON-RPC call.
	Params() []any
}

// JsonRPCAction is the action to execute when a user selects a result.
type JsonRPCAction struct {
	Method               string `json:"method"`
	Parameters           []any  `json:"parameters"`
	DontHideAfterAction  bool   `json:"dontHideAfterAction,omitempty"`
}

// Result is a single result item in a Flow Launcher query response.
type Result struct {
	Title       string         `json:"Title"`
	SubTitle    string         `json:"SubTitle,omitempty"`
	IcoPath     string         `json:"IcoPath,omitempty"`
	Score       int            `json:"Score,omitempty"`
	Action      *JsonRPCAction `json:"JsonRPCAction,omitempty"`
	ContextData any            `json:"ContextData,omitempty"`
}

// QueryResponse is the response sent back to Flow Launcher for a query.
type QueryResponse struct {
	Result []Result `json:"result"`
}

// ContextMenuResponse is the response sent back to Flow Launcher for a context menu request.
type ContextMenuResponse struct {
	Result []Result `json:"result"`
}

// NewResult creates a Result whose JsonRPCAction calls the given API method.
func NewResult(title, subtitle string, api APIMethod) *Result {
	return &Result{
		Title:    title,
		SubTitle: subtitle,
		Action: &JsonRPCAction{
			Method:     api.Method(),
			Parameters: api.Params(),
		},
	}
}

// NewResultWithIcoPath creates a Result with an icon path.
func NewResultWithIcoPath(title, subtitle, icoPath string, api APIMethod) *Result {
	r := NewResult(title, subtitle, api)
	r.IcoPath = icoPath
	return r
}

// MarshalJSON returns the JSON encoding of the result.
func (r *Result) MarshalJSON() ([]byte, error) {
	type Alias Result
	return json.Marshal((*Alias)(r))
}