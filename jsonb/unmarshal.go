package jsonb

import (
	"bytes"
	"reflect"

	"github.com/si3nloong/sqlike/reflext"
	"golang.org/x/xerrors"
)

// Unmarshaller :
type Unmarshaller interface {
	UnmarshalJSONB([]byte) error
}

var skipByteMap = map[byte]bool{
	' ':  true,
	'\n': true,
	'\t': true,
	'\r': true,
}

var sequenceMap = map[byte][]byte{
	'{': []byte{'}', '"'},
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

	buf := bytes.NewBuffer(data)
	vv := reflext.Zero(t)
	vv = reflext.Indirect(vv)
	if err := decoder(buf, vv); err != nil {
		return err
	}
	reflext.Indirect(v).Set(vv)
	return nil
}
