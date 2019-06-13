package jsonb

import (
	"reflect"
	"sync"
)

// ValueDecoder :
type ValueDecoder func(*Reader, reflect.Value) error

// ValueEncoder :
type ValueEncoder func(*Writer, reflect.Value) error

// Registry :
type Registry struct {
	mutex        sync.Mutex
	typeEncoders map[reflect.Type]ValueEncoder
	typeDecoders map[reflect.Type]ValueDecoder
	kindEncoders map[reflect.Kind]ValueEncoder
	kindDecoders map[reflect.Kind]ValueDecoder
}

var registry = buildRegistry()

func buildRegistry() *Registry {
	rg := NewRegistry()
	Decoder{}.SetDecoders(rg)
	Encoder{}.SetEncoders(rg)
	return rg
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
func (r *Registry) LookupEncoder(t reflect.Type) (ValueEncoder, error) {
	var (
		enc  ValueEncoder
		isOk bool
	)

	it := reflect.TypeOf((*Marshaller)(nil)).Elem()
	if t.Implements(it) {
		return marshallerEncoder(), nil
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

	it := reflect.TypeOf((*Unmarshaller)(nil)).Elem()
	if t.Implements(it) {
		return unmarshallerDecoder(), nil
	}

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
