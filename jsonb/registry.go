package jsonb

import (
	"github.com/si3nloong/sqlike/core/codec"
)

var registry = buildRegistry()

func buildRegistry() *codec.Registry {
	rg := codec.NewRegistry()
	ValueEncoder{}.SetEncoders(rg)
	ValueDecoder{}.SetDecoders(rg)
	return rg
}
