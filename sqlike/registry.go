package sqlike

import "github.com/si3nloong/sqlike/core/codec"

// DefaultMapper :
var (
	DefaultRegistry = buildDefaultRegistry()
)

func buildDefaultRegistry() *codec.Registry {
	rg := codec.NewRegistry()
	DefaultDecoders{}.SetDecoders(rg)
	// DefaultEncoders{}.SetEncoders(rg)
	return rg
}
