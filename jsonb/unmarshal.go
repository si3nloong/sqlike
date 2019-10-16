package jsonb

import (
	"encoding/json"
	"fmt"
	"reflect"

	"errors"

	"github.com/si3nloong/sqlike/reflext"
)

// Unmarshaler :
type Unmarshaler interface {
	UnmarshalJSONB([]byte) error
}

// Unmarshal :
func Unmarshal(data []byte, dst interface{}) error {
	v := reflext.ValueOf(dst)
	if !v.IsValid() {
		return errors.New("invalid value for jsonb.Unmarshal")
	}

	t := v.Type()
	if !reflext.IsKind(t, reflect.Ptr) {
		return errors.New("unaddressable destination")
	}

	v = v.Elem()
	if !v.CanSet() {
		return errors.New("unaddressable destination")
	}

	t = t.Elem()
	decoder, err := registry.LookupDecoder(t)
	if err != nil {
		return err
	}

	r := NewReader(data)
	vv := reflext.Zero(t)
	if err := decoder(r, vv); err != nil {
		return err
	}

	c := r.nextToken()
	if c != 0 {
		return ErrInvalidJSON{
			callback: "Unmarshal",
			message:  fmt.Sprintf("invalid json string, extra char %q found", c),
		}
	}
	v.Set(vv)
	return nil
}

// UnmarshalValue :
func UnmarshalValue(data []byte, v reflect.Value) error {
	if data == nil {
		v.Set(reflect.Zero(v.Type()))
		return nil
	}

	t := v.Type()
	decoder, err := registry.LookupDecoder(reflext.Deref(t))
	if err != nil {
		return err
	}

	r := NewReader(data)
	vv := reflext.Zero(t)

	vv = reflext.Indirect(vv)
	if err := decoder(r, vv); err != nil {
		return err
	}

	c := r.nextToken()
	if c != 0 {
		return ErrInvalidJSON{
			callback: "Unmarshal",
		}
	}
	reflext.Indirect(v).Set(vv)
	return nil
}

// unmarshalerDecoder
func unmarshalerDecoder() ValueDecoder {
	return func(r *Reader, v reflect.Value) error {
		x, ok := v.Interface().(Unmarshaler)
		if !ok {
			return errors.New("codec: invalid type for assertion")
		}
		return x.UnmarshalJSONB(r.Bytes())
	}
}

// jsonUnmarshalerDecoder
func jsonUnmarshalerDecoder() ValueDecoder {
	return func(r *Reader, v reflect.Value) error {
		x, ok := v.Interface().(json.Unmarshaler)
		if !ok {
			return errors.New("codec: invalid type for assertion")
		}
		return x.UnmarshalJSON(r.Bytes())
	}
}
