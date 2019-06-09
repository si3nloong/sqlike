package jsonb

import (
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
	t := v.Type()
	if !reflext.IsKind(t, reflect.Ptr) {
		return xerrors.New("unaddressable destination")
	}

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
		return ErrDecode{
			callback: "Unmarshal",
		}
	}
	reflext.Indirect(v).Set(vv)
	return nil
}

// UnmarshalValue :
func UnmarshalValue(data []byte, v reflect.Value) error {
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
		return ErrDecode{
			callback: "Unmarshal",
		}
	}
	reflext.Indirect(v).Set(vv)
	return nil
}
