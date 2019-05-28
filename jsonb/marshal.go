package jsonb

import (
	"bytes"

	"github.com/si3nloong/sqlike/reflext"
)

var jsonNull = []byte(`null`)

// Marshaller :
type Marshaller interface {
	MarshalJSONB() ([]byte, error)
}

// Marshal :
func Marshal(src interface{}) (b []byte, err error) {
	v := reflext.ValueOf(src)
	if src == nil || reflext.IsNull(v) {
		b = jsonNull
		return
	}

	t := reflext.Deref(v.Type())
	encoder, err := registry.LookupEncoder(t)
	if err != nil {
		return nil, err
	}

	v = reflext.Indirect(v)
	w := new(bytes.Buffer)
	if err := encoder(w, v); err != nil {
		return nil, err
	}

	b = w.Bytes()
	return
}
