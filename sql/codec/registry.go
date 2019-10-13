package codec

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"reflect"
	"sync"

	"github.com/si3nloong/sqlike/reflext"
)

// DefaultMapper :
var (
	DefaultRegistry = buildDefaultRegistry()
)

func buildDefaultRegistry() *Registry {
	rg := NewRegistry()
	DefaultDecoders{}.SetDecoders(rg)
	DefaultEncoders{}.SetEncoders(rg)
	return rg
}

// Registry :
type Registry struct {
	mutex        sync.Mutex
	typeEncoders map[reflect.Type]ValueEncoder
	typeDecoders map[reflect.Type]ValueDecoder
	kindEncoders map[reflect.Kind]ValueEncoder
	kindDecoders map[reflect.Kind]ValueDecoder
}

// NewRegistry creates a new empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		typeEncoders: make(map[reflect.Type]ValueEncoder),
		typeDecoders: make(map[reflect.Type]ValueDecoder),
		kindEncoders: make(map[reflect.Kind]ValueEncoder),
		kindDecoders: make(map[reflect.Kind]ValueDecoder),
	}
}

// SetTypeEncoder :
func (r *Registry) SetTypeEncoder(t reflect.Type, enc ValueEncoder) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.typeEncoders[t] = enc
}

// SetTypeDecoder :
func (r *Registry) SetTypeDecoder(t reflect.Type, dec ValueDecoder) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.typeDecoders[t] = dec
}

// SetKindEncoder :
func (r *Registry) SetKindEncoder(k reflect.Kind, enc ValueEncoder) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.kindEncoders[k] = enc
}

// SetKindDecoder :
func (r *Registry) SetKindDecoder(k reflect.Kind, dec ValueDecoder) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.kindDecoders[k] = dec
}

// LookupEncoder :
func (r *Registry) LookupEncoder(v reflect.Value) (ValueEncoder, error) {
	var (
		enc ValueEncoder
		ok  bool
	)

	if !v.IsValid() || reflext.IsNull(v) {
		return func(_ *reflext.StructField, _ reflect.Value) (interface{}, error) {
			return nil, nil
		}, nil
	}

	if _, ok := v.Interface().(driver.Valuer); ok {
		return encodeValue, nil
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
func (r *Registry) LookupDecoder(t reflect.Type) (ValueDecoder, error) {
	var (
		dec ValueDecoder
		ok  bool
	)

	v := reflext.Zero(t)
	if _, ok := v.Addr().Interface().(sql.Scanner); ok {
		return func(it interface{}, v reflect.Value) error {
			return v.Addr().Interface().(sql.Scanner).Scan(it)
		}, nil
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

func encodeValue(_ *reflext.StructField, v reflect.Value) (interface{}, error) {
	x, ok := v.Interface().(driver.Valuer)
	if !ok {
		return nil, errors.New("codec: invalid type for assertion")
	}
	return x.Value()
}
