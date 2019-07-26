package jsonb

import (
	"fmt"
	"reflect"

	"github.com/si3nloong/sqlike/reflext"
	"golang.org/x/xerrors"
)

// Unmarshaller :
type Unmarshaller interface {
	UnmarshalJSONB([]byte) error
}

// Unmarshal :
func Unmarshal(data []byte, dst interface{}) error {
	v := reflext.ValueOf(dst)
	if !v.IsValid() {
		return xerrors.New("invalid value for jsonb.Unmarshal")
	}

	t := v.Type()
	if !reflext.IsKind(t, reflect.Ptr) {
		return xerrors.New("unaddressable destination")
	}

	v = v.Elem()
	if !v.CanSet() {
		return xerrors.New("unaddressable destination")
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

// unmarshallerDecoder
func unmarshallerDecoder() ValueDecoder {
	return func(r *Reader, v reflect.Value) error {
		it := v.Interface()
		return it.(Unmarshaller).UnmarshalJSONB(r.Bytes())
	}
}
