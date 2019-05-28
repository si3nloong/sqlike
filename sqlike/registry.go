package sqlike

import "bitbucket.org/SianLoong/sqlike/core/codec"

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
