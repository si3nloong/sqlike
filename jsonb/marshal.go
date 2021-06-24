package jsonb

import (
	"encoding"
	"encoding/json"
	"reflect"

	"github.com/si3nloong/sqlike/v2/x/reflext"
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
		x := v.Interface().(Marshaler)
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
		x := v.Interface().(json.Marshaler)
		b, err := x.MarshalJSON()
		if err != nil {
			return err
		}
		w.Write(b)
		return nil
	}
}

func textMarshalerEncoder() ValueEncoder {
	return func(w *Writer, v reflect.Value) error {
		x := v.Interface().(encoding.TextMarshaler)
		b, err := x.MarshalText()
		if err != nil {
			return err
		}
		length := len(b)
		w.WriteByte('"')
		for i := 0; i < length; i++ {
			char := b[0]
			b = b[1:]
			if x, ok := escapeCharMap[char]; ok {
				w.Write(x)
				continue
			}
			w.WriteByte(char)
		}
		w.WriteByte('"')
		return nil
	}
}
