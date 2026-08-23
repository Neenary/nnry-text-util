module github.com/neenary/nnry-text-util

go 1.26.5

require (
	github.com/neenary/flow-launcher-go v0.0.0
	github.com/neenary/jsonrpc v0.0.0
	github.com/spf13/cobra v1.10.2
	github.com/spf13/pflag v1.0.9
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
)

replace (
	github.com/neenary/flow-launcher-go => ./flow-launcher-go
	github.com/neenary/jsonrpc => ./jsonrpc
)