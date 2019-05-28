package sqlike

import (
	"log"
	"reflect"
	"bitbucket.org/SianLoong/sqlike/core"
	"bitbucket.org/SianLoong/sqlike/reflext"
	"golang.org/x/xerrors"
)

// ErrNoResult :
var ErrNoResult = xerrors.New("goloquent: no result in return")

// ErrUnaddressableEntity :
var ErrUnaddressableEntity = xerrors.New("unaddressable entity")

// Result :
type Result struct {
	err error
	csr *Cursor
}

// Decode :
func (r Result) Decode(dest interface{}) error {
	if r.err != nil {
		return r.err
	}

	v := reflext.Indirect(reflect.ValueOf(dest))
	if !v.IsValid() || !v.CanSet() {
		return ErrUnaddressableEntity
	}

	t := v.Type()
	cdc := core.DefaultMapper.CodecByType(t)
	log.Println(cdc)
	// vv := reflect.New(vi.Type())
	// v.Elem().Set(vv.Elem())
	return nil
}

// Error :
func (r *Result) Error() error {
	return r.err
}

// Close :
func (r *Result) Close() error {
	return r.csr.Close()
}
