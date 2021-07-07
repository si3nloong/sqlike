package codec

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"reflect"
	"sync"

	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/x/reflext"
)

var (
	sqlScanner = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
)

// Registry :
type Registry struct {
	mutex        sync.Mutex
	typeEncoders map[reflect.Type]db.ValueEncoder
	typeDecoders map[reflect.Type]db.ValueDecoder
	kindEncoders map[reflect.Kind]db.ValueEncoder
	kindDecoders map[reflect.Kind]db.ValueDecoder
}

var _ db.Codecer = (*Registry)(nil)

// NewRegistry : creates a new empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		typeEncoders: make(map[reflect.Type]db.ValueEncoder),
		typeDecoders: make(map[reflect.Type]db.ValueDecoder),
		kindEncoders: make(map[reflect.Kind]db.ValueEncoder),
		kindDecoders: make(map[reflect.Kind]db.ValueDecoder),
	}
}

// RegisterTypeCodec :
func (r *Registry) RegisterTypeCodec(t reflect.Type, enc db.ValueEncoder, dec db.ValueDecoder) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.typeEncoders[t] = enc
	r.typeDecoders[t] = dec
}

// RegisterTypeEncoder :
func (r *Registry) RegisterTypeEncoder(t reflect.Type, enc db.ValueEncoder) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.typeEncoders[t] = enc
}

// RegisterTypeDecoder :
func (r *Registry) RegisterTypeDecoder(t reflect.Type, dec db.ValueDecoder) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.typeDecoders[t] = dec
}

// RegisterKindCodec :
func (r *Registry) RegisterKindCodec(k reflect.Kind, enc db.ValueEncoder, dec db.ValueDecoder) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.kindEncoders[k] = enc
	r.kindDecoders[k] = dec
}

// RegisterKindEncoder :
func (r *Registry) RegisterKindEncoder(k reflect.Kind, enc db.ValueEncoder) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.kindEncoders[k] = enc
}

// RegisterKindDecoder :
func (r *Registry) RegisterKindDecoder(k reflect.Kind, dec db.ValueDecoder) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.kindDecoders[k] = dec
}

// LookupEncoder :
func (r *Registry) LookupEncoder(v reflect.Value) (db.ValueEncoder, error) {
	var (
		enc db.ValueEncoder
		ok  bool
	)

	if !v.IsValid() {
		return NilEncoder, nil
	}

	if _, ok := v.Interface().(driver.Valuer); ok {
		return encodeDriverValue, nil
	}

	t := v.Type()
	enc, ok = r.typeEncoders[t]
	if ok {
		return enc, nil
	}

	enc, ok = r.kindEncoders[t.Kind()]
	if ok {
		return enc, nil
	}
	return nil, ErrNoEncoder{Type: t}
}

// LookupDecoder :
func (r *Registry) LookupDecoder(t reflect.Type) (db.ValueDecoder, error) {
	var (
		dec db.ValueDecoder
		ok  bool
	)

	ptrType := t
	if t.Kind() != reflect.Ptr {
		ptrType = reflect.PtrTo(t)
	}

	if ptrType.Implements(sqlScanner) {
		return sqlScannerDecoder, nil
	}

	dec, ok = r.typeDecoders[t]
	if ok {
		return dec, nil
	}

	dec, ok = r.kindDecoders[t.Kind()]
	if ok {
		return dec, nil
	}
	return nil, ErrNoDecoder{Type: t}
}

// NilEncoder :
func NilEncoder(_ context.Context, _ reflect.Value) (interface{}, error) {
	return nil, nil
}

func encodeDriverValue(_ context.Context, v reflect.Value) (interface{}, error) {
	if !v.IsValid() || reflext.IsNull(v) {
		return nil, nil
	}
	x, ok := v.Interface().(driver.Valuer)
	if !ok {
		return nil, errors.New("codec: invalid type for assertion")
	}
	return x.Value()
}

func sqlScannerDecoder(_ context.Context, it interface{}, v reflect.Value) error {
	if it == nil {
		// Avoid from sql.scanner when the value is nil
		v.Set(reflect.Zero(v.Type()))
		return nil
	}

	if v.Kind() != reflect.Ptr {
		return v.Addr().Interface().(sql.Scanner).Scan(it)
	}

	return reflext.Init(v).Interface().(sql.Scanner).Scan(it)
}
