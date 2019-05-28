package jsonb

import (
	"bitbucket.org/SianLoong/sqlike/core/codec"
)

var registry = buildRegistry()

func buildRegistry() *codec.Registry {
	rg := codec.NewRegistry()
	ValueEncoder{}.SetEncoders(rg)
	ValueDecoder{}.SetDecoders(rg)
	return rg
}
