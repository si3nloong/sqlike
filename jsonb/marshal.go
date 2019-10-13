package jsonb

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/si3nloong/sqlike/reflext"
)

// Marshaler :
type Marshaler interface {
	MarshalJSONB() ([]byte, error)
}

// Marshal :
func Marshal(src interface{}) (b []byte, err error) {
	v := reflext.ValueOf(src)
	if src == nil || !v.IsValid() || reflext.IsNull(v) {
		b = []byte(null)
		return
	}

	encoder, err := registry.LookupEncoder(v)
	if err != nil {
		return nil, err
	}

	w := NewWriter()
	if err := encoder(w, v); err != nil {
		return nil, err
	}
	b = w.Bytes()
	return
}

// marshalerEncoder
func marshalerEncoder() ValueEncoder {
	return func(w *Writer, v reflect.Value) error {
		x, ok := v.Interface().(Marshaler)
		if !ok {
			return errors.New("codec: invalid type for assertion")
		}
		b, err := x.MarshalJSONB()
		if err != nil {
			return err
		}
		w.Write(b)
		return nil
	}
}

func jsonMarshalerEncoder() ValueEncoder {
	return func(w *Writer, v reflect.Value) error {
		x, ok := v.Interface().(json.Marshaler)
		if !ok {
			return errors.New("codec: invalid type for assertion")
		}
		b, err := x.MarshalJSON()
		if err != nil {
			return err
		}
		w.Write(b)
		return nil
	}
}
