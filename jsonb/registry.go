package jsonb

import (
	"encoding"
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"github.com/si3nloong/sqlike/reflext"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

var (
	registry = buildDefaultRegistry()

	jsonbUnmarshaler = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
	jsonUnmarshaler  = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	textUnmarshaler  = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	textMarshaler    = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

// ValueDecoder :
type ValueDecoder func(*Reader, reflect.Value) error

// ValueEncoder :
type ValueEncoder func(*Writer, reflect.Value) error

// Registry :
type Registry struct {
	mutex        *sync.Mutex
	typeEncoders map[reflect.Type]ValueEncoder
	typeDecoders map[reflect.Type]ValueDecoder
	kindEncoders map[reflect.Kind]ValueEncoder
	kindDecoders map[reflect.Kind]ValueDecoder
}

func buildDefaultRegistry() *Registry {
	rg := NewRegistry()
	enc := DefaultEncoder{rg}
	dec := DefaultDecoder{rg}
	rg.SetTypeCoder(reflect.TypeOf([]byte{}), enc.EncodeByte, dec.DecodeByte)
	rg.SetTypeCoder(reflect.TypeOf(language.Tag{}), enc.EncodeStringer, dec.DecodeLanguage)
	rg.SetTypeCoder(reflect.TypeOf(currency.Unit{}), enc.EncodeStringer, dec.DecodeCurrency)
	rg.SetTypeCoder(reflect.TypeOf(time.Time{}), enc.EncodeTime, dec.DecodeTime)
	rg.SetTypeCoder(reflect.TypeOf(json.RawMessage{}), enc.EncodeJSONRaw, dec.DecodeJSONRaw)
	rg.SetTypeCoder(reflect.TypeOf(json.Number("")), enc.EncodeStringer, dec.DecodeJSONNumber)
	rg.SetKindCoder(reflect.String, enc.EncodeString, dec.DecodeString)
	rg.SetKindCoder(reflect.Bool, enc.EncodeBool, dec.DecodeBool)
	rg.SetKindCoder(reflect.Int, enc.EncodeInt, dec.DecodeInt(false))
	rg.SetKindCoder(reflect.Int8, enc.EncodeInt, dec.DecodeInt(false))
	rg.SetKindCoder(reflect.Int16, enc.EncodeInt, dec.DecodeInt(false))
	rg.SetKindCoder(reflect.Int32, enc.EncodeInt, dec.DecodeInt(false))
	rg.SetKindCoder(reflect.Int64, enc.EncodeInt, dec.DecodeInt(false))
	rg.SetKindCoder(reflect.Uint, enc.EncodeUint, dec.DecodeUint(false))
	rg.SetKindCoder(reflect.Uint8, enc.EncodeUint, dec.DecodeUint(false))
	rg.SetKindCoder(reflect.Uint16, enc.EncodeUint, dec.DecodeUint(false))
	rg.SetKindCoder(reflect.Uint32, enc.EncodeUint, dec.DecodeUint(false))
	rg.SetKindCoder(reflect.Uint64, enc.EncodeUint, dec.DecodeUint(false))
	rg.SetKindCoder(reflect.Float32, enc.EncodeFloat, dec.DecodeFloat)
	rg.SetKindCoder(reflect.Float64, enc.EncodeFloat, dec.DecodeFloat)
	rg.SetKindCoder(reflect.Ptr, enc.EncodePtr, dec.DecodePtr)
	rg.SetKindCoder(reflect.Struct, enc.EncodeStruct, dec.DecodeStruct)
	rg.SetKindCoder(reflect.Array, enc.EncodeArray, dec.DecodeArray)
	rg.SetKindCoder(reflect.Slice, enc.EncodeArray, dec.DecodeSlice)
	rg.SetKindCoder(reflect.Map, enc.EncodeMap, dec.DecodeMap)
	rg.SetKindCoder(reflect.Interface, enc.EncodeInterface, dec.DecodeInterface)
	return rg
}

// NewRegistry creates a new empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		mutex:        new(sync.Mutex),
		typeEncoders: make(map[reflect.Type]ValueEncoder),
		typeDecoders: make(map[reflect.Type]ValueDecoder),
		kindEncoders: make(map[reflect.Kind]ValueEncoder),
		kindDecoders: make(map[reflect.Kind]ValueDecoder),
	}
}

// SetTypeCoder :
func (r *Registry) SetTypeCoder(t reflect.Type, enc ValueEncoder, dec ValueDecoder) {
	if enc == nil {
		panic("missing encoder")
	}
	if dec == nil {
		panic("missing decoder")
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.typeEncoders[t] = enc
	r.typeDecoders[t] = dec
}

// SetKindCoder :
func (r *Registry) SetKindCoder(k reflect.Kind, enc ValueEncoder, dec ValueDecoder) {
	if enc == nil {
		panic("missing encoder")
	}
	if dec == nil {
		panic("missing decoder")
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.kindEncoders[k] = enc
	r.kindDecoders[k] = dec
}

// LookupEncoder :
func (r *Registry) LookupEncoder(v reflect.Value) (ValueEncoder, error) {
	var (
		enc ValueEncoder
		ok  bool
	)

	if !v.IsValid() || reflext.IsNull(v) {
		return func(w *Writer, v reflect.Value) error {
			w.WriteString(null)
			return nil
		}, nil
	}

	it := v.Interface()
	if _, ok := it.(Marshaler); ok {
		return marshalerEncoder(), nil
	}

	if _, ok := it.(json.Marshaler); ok {
		return jsonMarshalerEncoder(), nil
	}

	if _, ok := it.(encoding.TextMarshaler); ok {
		return textMarshalerEncoder(), nil
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

	ptrType := t
	if t.Kind() != reflect.Ptr {
		ptrType = reflect.PtrTo(t)
	}

	if ptrType.Implements(jsonbUnmarshaler) {
		return unmarshalerDecoder(), nil
	}

	dec, ok = r.typeDecoders[t]
	if ok {
		return dec, nil
	}

	if ptrType.Implements(jsonUnmarshaler) {
		return jsonUnmarshalerDecoder(), nil
	}

	if ptrType.Implements(textUnmarshaler) {
		return textUnmarshalerDecoder(), nil
	}

	dec, ok = r.kindDecoders[t.Kind()]
	if ok {
		return dec, nil
	}
	return nil, ErrNoDecoder{Type: t}
}
