package jsonb

import (
	"bytes"
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"

	"errors"

	"github.com/si3nloong/sqlike/x/reflext"
)

// Unmarshaler :
type Unmarshaler interface {
	UnmarshalJSONB([]byte) error
}

// Unmarshal :
func Unmarshal(data []byte, dst interface{}) error {
	v := reflext.ValueOf(dst)
	if !v.IsValid() {
		return errors.New("invalid value for Unmarshal")
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
		if r.IsNull() {
			v.Set(reflect.New(v.Type()).Elem())
			return r.skipNull()
		}

		r.pos = r.len
		if v.Kind() != reflect.Ptr {
			return v.Addr().Interface().(Unmarshaler).UnmarshalJSONB(r.Bytes())
		}
		return reflext.Init(v).Interface().(Unmarshaler).UnmarshalJSONB(r.Bytes())
	}
}

// jsonUnmarshalerDecoder
func jsonUnmarshalerDecoder() ValueDecoder {
	return func(r *Reader, v reflect.Value) error {
		if r.IsNull() {
			v.Set(reflect.New(v.Type()).Elem())
			return r.skipNull()
		}

		r.pos = r.len
		if v.Kind() != reflect.Ptr {
			return v.Addr().Interface().(json.Unmarshaler).UnmarshalJSON(r.Bytes())
		}
		return reflext.Init(v).Interface().(json.Unmarshaler).UnmarshalJSON(r.Bytes())
	}
}

// textUnmarshalerDecoder
func textUnmarshalerDecoder() ValueDecoder {
	return func(r *Reader, v reflect.Value) error {
		if r.IsNull() {
			v.Set(reflect.New(v.Type()).Elem())
			return r.skipNull()
		}

		r.pos = r.len
		b := r.Bytes()
		b = bytes.Trim(b, `"`)
		if v.Kind() != reflect.Ptr {
			return v.Addr().Interface().(encoding.TextUnmarshaler).UnmarshalText(b)
		}
		return reflext.Init(v).Interface().(encoding.TextUnmarshaler).UnmarshalText(b)
	}
}
