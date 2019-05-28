package codec

import (
	"database/sql/driver"
	"reflect"
	"sync"

	"bitbucket.org/SianLoong/sqlike/reflext"
)

// DefaultMapper :
var (
	DefaultRegistry = buildDefaultRegistry()
)

func buildDefaultRegistry() *Registry {
	rg := NewRegistry()
	// DefaultDecoders{}.SetDecoders(rg)
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

// func (r *Registry) Encode(t reflect.Type) (interface{}, error) {

// }

// LookupEncoder :
func (r *Registry) LookupEncoder(t reflect.Type) (ValueEncoder, error) {
	var (
		enc  ValueEncoder
		isOk bool
	)
	v := reflext.Zero(t)
	if _, isOk := v.Interface().(driver.Valuer); isOk {
		return func(v reflect.Value) (interface{}, error) {
			return v.Interface().(driver.Valuer).Value()
		}, nil
	}

	enc, isOk = r.typeEncoders[t]
	if isOk {
		return enc, nil
	}

	enc, isOk = r.kindEncoders[t.Kind()]
	if isOk {
		return enc, nil
	}
	return nil, ErrNoEncoder{Type: t}
}

// LookupDecoder :
func (r *Registry) LookupDecoder(t reflect.Type) (ValueDecoder, error) {
	var (
		dec  ValueDecoder
		isOk bool
	)
	dec, isOk = r.typeDecoders[t]
	if isOk {
		return dec, nil
	}

	dec, isOk = r.kindDecoders[t.Kind()]
	if isOk {
		return dec, nil
	}
	return nil, ErrNoDecoder{Type: t}
}
