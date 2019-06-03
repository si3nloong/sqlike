package jsonb

import (
	"log"
	"reflect"

	"github.com/si3nloong/sqlike/reflext"
	"golang.org/x/xerrors"
)

// Unmarshaller :
type Unmarshaller interface {
	UnmarshalJSONB([]byte) error
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

	log.Println("Data :", string(data))
	r := NewReader(data)
	vv := reflext.Zero(t)

	vv = reflext.Indirect(vv)
	if err := decoder(r, vv); err != nil {
		return err
	}

	c := r.nextToken()
	log.Println("leftover byte :", c, string(c))
	// if c != 0 {
	// 	return xerrors.New("invalid json format, extra characters")
	// }
	reflext.Indirect(v).Set(vv)
	return nil
}
