package codec

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/paulmach/orb"
	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/spatial"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

type Codecer interface {
	SetTypeCoder(t reflect.Type, enc ValueEncoder, dec ValueDecoder)
	SetTypeEncoder(t reflect.Type, enc ValueEncoder)
	SetTypeDecoder(t reflect.Type, dec ValueDecoder)
	SetKindEncoder(k reflect.Kind, enc ValueEncoder)
	SetKindDecoder(k reflect.Kind, dec ValueDecoder)
	LookupEncoder(v reflect.Value) (ValueEncoder, error)
	LookupDecoder(t reflect.Type) (ValueDecoder, error)
}

// DefaultMapper :
var (
	DefaultRegistry = buildDefaultRegistry()
	sqlScanner      = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
)

func buildDefaultRegistry() Codecer {
	rg := NewRegistry()
	dec := DefaultDecoders{rg}
	enc := DefaultEncoders{rg}
	rg.SetTypeCoder(reflect.TypeOf([]byte{}), enc.EncodeByte, dec.DecodeByte)
	rg.SetTypeCoder(reflect.TypeOf(language.Tag{}), enc.EncodeStringer, dec.DecodeLanguage)
	rg.SetTypeCoder(reflect.TypeOf(currency.Unit{}), enc.EncodeStringer, dec.DecodeCurrency)
	rg.SetTypeCoder(reflect.TypeOf(time.Time{}), enc.EncodeTime, dec.DecodeTime)
	rg.SetTypeCoder(reflect.TypeOf(sql.RawBytes{}), enc.EncodeRawBytes, dec.DecodeRawBytes)
	rg.SetTypeCoder(reflect.TypeOf(json.RawMessage{}), enc.EncodeJSONRaw, dec.DecodeJSONRaw)
	rg.SetTypeCoder(reflect.TypeOf(orb.Point{}), enc.EncodeSpatial(spatial.Point), dec.DecodePoint)
	rg.SetTypeCoder(reflect.TypeOf(orb.LineString{}), enc.EncodeSpatial(spatial.LineString), dec.DecodeLineString)
	// rg.SetTypeCoder(reflect.TypeOf(orb.Polygon{}), enc.EncodeSpatial(spatial.Polygon), dec.DecodePolygon)
	// rg.SetTypeCoder(reflect.TypeOf(orb.MultiPoint{}), enc.EncodeSpatial(spatial.MultiPoint), dec.DecodeMultiPoint)
	// rg.SetTypeCoder(reflect.TypeOf(orb.MultiLineString{}), enc.EncodeSpatial(spatial.MultiLineString), dec.DecodeMultiLineString)
	// rg.SetTypeCoder(reflect.TypeOf(orb.MultiPolygon{}), enc.EncodeSpatial(spatial.MultiPolygon), dec.DecodeMultiPolygon)
	rg.SetKindCoder(reflect.String, enc.EncodeString, dec.DecodeString)
	rg.SetKindCoder(reflect.Bool, enc.EncodeBool, dec.DecodeBool)
	rg.SetKindCoder(reflect.Int, enc.EncodeInt, dec.DecodeInt)
	rg.SetKindCoder(reflect.Int8, enc.EncodeInt, dec.DecodeInt)
	rg.SetKindCoder(reflect.Int16, enc.EncodeInt, dec.DecodeInt)
	rg.SetKindCoder(reflect.Int32, enc.EncodeInt, dec.DecodeInt)
	rg.SetKindCoder(reflect.Int64, enc.EncodeInt, dec.DecodeInt)
	rg.SetKindCoder(reflect.Uint, enc.EncodeUint, dec.DecodeUint)
	rg.SetKindCoder(reflect.Uint8, enc.EncodeUint, dec.DecodeUint)
	rg.SetKindCoder(reflect.Uint16, enc.EncodeUint, dec.DecodeUint)
	rg.SetKindCoder(reflect.Uint32, enc.EncodeUint, dec.DecodeUint)
	rg.SetKindCoder(reflect.Uint64, enc.EncodeUint, dec.DecodeUint)
	rg.SetKindCoder(reflect.Float32, enc.EncodeFloat, dec.DecodeFloat)
	rg.SetKindCoder(reflect.Float64, enc.EncodeFloat, dec.DecodeFloat)
	rg.SetKindCoder(reflect.Ptr, enc.EncodePtr, dec.DecodePtr)
	rg.SetKindCoder(reflect.Struct, enc.EncodeStruct, dec.DecodeStruct)
	rg.SetKindCoder(reflect.Array, enc.EncodeArray, dec.DecodeArray)
	rg.SetKindCoder(reflect.Slice, enc.EncodeArray, dec.DecodeArray)
	rg.SetKindCoder(reflect.Map, enc.EncodeMap, dec.DecodeMap)
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

var _ Codecer = (*Registry)(nil)

// NewRegistry creates a new empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		typeEncoders: make(map[reflect.Type]ValueEncoder),
		typeDecoders: make(map[reflect.Type]ValueDecoder),
		kindEncoders: make(map[reflect.Kind]ValueEncoder),
		kindDecoders: make(map[reflect.Kind]ValueDecoder),
	}
}

// SetTypeCoder :
func (r *Registry) SetTypeCoder(t reflect.Type, enc ValueEncoder, dec ValueDecoder) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.typeEncoders[t] = enc
	r.typeDecoders[t] = dec
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

// SetKindCoder :
func (r *Registry) SetKindCoder(k reflect.Kind, enc ValueEncoder, dec ValueDecoder) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.kindEncoders[k] = enc
	r.kindDecoders[k] = dec
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

	// if !v.IsValid() || reflext.IsNull(v) {
	// 	return NilEncoder, nil
	// }
	if !v.IsValid() {
		return NilEncoder, nil
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

	ptrType := t
	if t.Kind() != reflect.Ptr {
		ptrType = reflect.PtrTo(t)
	}

	if ptrType.Implements(sqlScanner) {
		return func(it interface{}, v reflect.Value) error {
			if v.Kind() != reflect.Ptr {
				return v.Addr().Interface().(sql.Scanner).Scan(it)
			}
			return reflext.Init(v).Interface().(sql.Scanner).Scan(it)
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
	if !v.IsValid() || reflext.IsNull(v) {
		return nil, nil
	}
	x, ok := v.Interface().(driver.Valuer)
	if !ok {
		return nil, errors.New("codec: invalid type for assertion")
	}
	return x.Value()
}

// NilEncoder :
func NilEncoder(_ *reflext.StructField, _ reflect.Value) (interface{}, error) {
	return nil, nil
}
