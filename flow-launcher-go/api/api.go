// Package api provides typed Go structs for the Flow Launcher JSON-RPC API.
// Each struct represents a callable method on the IPublicAPI interface.
package api

// Method is implemented by all Flow Launcher API call structs.
type Method interface {
	// Method returns the Flow.Launcher.* method name.
	Method() string
	// Params returns the positional parameters for the JSON-RPC call.
	Params() []any
}