package codec

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
	"sync"
	"time"

	"cloud.google.com/go/civil"
	"github.com/paulmach/orb"
	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/spatial"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

// Codecer :
type Codecer interface {
	RegisterTypeCodec(t reflect.Type, enc ValueEncoder, dec ValueDecoder)
	RegisterTypeEncoder(t reflect.Type, enc ValueEncoder)
	RegisterTypeDecoder(t reflect.Type, dec ValueDecoder)
	RegisterKindEncoder(k reflect.Kind, enc ValueEncoder)
	RegisterKindDecoder(k reflect.Kind, dec ValueDecoder)
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
	rg.RegisterTypeCodec(reflect.TypeOf([]byte{}), enc.EncodeByte, dec.DecodeByte)
	rg.RegisterTypeCodec(reflect.TypeOf(language.Tag{}), enc.EncodeStringer, dec.DecodeLanguage)
	rg.RegisterTypeCodec(reflect.TypeOf(currency.Unit{}), enc.EncodeStringer, dec.DecodeCurrency)
	rg.RegisterTypeCodec(reflect.TypeOf(time.Time{}), enc.EncodeTime, dec.DecodeTime)
	rg.RegisterTypeCodec(reflect.TypeOf(civil.Date{}), enc.EncodeStringer, dec.DecodeCivilDate)
	rg.RegisterTypeCodec(reflect.TypeOf(sql.RawBytes{}), enc.EncodeRawBytes, dec.DecodeRawBytes)
	rg.RegisterTypeCodec(reflect.TypeOf(json.RawMessage{}), enc.EncodeJSONRaw, dec.DecodeJSONRaw)
	rg.RegisterTypeCodec(reflect.TypeOf(orb.Point{}), enc.EncodeSpatial(spatial.Point), dec.DecodePoint)
	rg.RegisterTypeCodec(reflect.TypeOf(orb.LineString{}), enc.EncodeSpatial(spatial.LineString), dec.DecodeLineString)
	// rg.RegisterTypeCodec(reflect.TypeOf(orb.Polygon{}), enc.EncodeSpatial(spatial.Polygon), dec.DecodePolygon)
	// rg.RegisterTypeCodec(reflect.TypeOf(orb.MultiPoint{}), enc.EncodeSpatial(spatial.MultiPoint), dec.DecodeMultiPoint)
	// rg.RegisterTypeCodec(reflect.TypeOf(orb.MultiLineString{}), enc.EncodeSpatial(spatial.MultiLineString), dec.DecodeMultiLineString)
	// rg.RegisterTypeCodec(reflect.TypeOf(orb.MultiPolygon{}), enc.EncodeSpatial(spatial.MultiPolygon), dec.DecodeMultiPolygon)
	rg.RegisterKindCodec(reflect.String, enc.EncodeString, dec.DecodeString)
	rg.RegisterKindCodec(reflect.Bool, enc.EncodeBool, dec.DecodeBool)
	rg.RegisterKindCodec(reflect.Int, enc.EncodeInt, dec.DecodeInt)
	rg.RegisterKindCodec(reflect.Int8, enc.EncodeInt, dec.DecodeInt)
	rg.RegisterKindCodec(reflect.Int16, enc.EncodeInt, dec.DecodeInt)
	rg.RegisterKindCodec(reflect.Int32, enc.EncodeInt, dec.DecodeInt)
	rg.RegisterKindCodec(reflect.Int64, enc.EncodeInt, dec.DecodeInt)
	rg.RegisterKindCodec(reflect.Uint, enc.EncodeUint, dec.DecodeUint)
	rg.RegisterKindCodec(reflect.Uint8, enc.EncodeUint, dec.DecodeUint)
	rg.RegisterKindCodec(reflect.Uint16, enc.EncodeUint, dec.DecodeUint)
	rg.RegisterKindCodec(reflect.Uint32, enc.EncodeUint, dec.DecodeUint)
	rg.RegisterKindCodec(reflect.Uint64, enc.EncodeUint, dec.DecodeUint)
	rg.RegisterKindCodec(reflect.Float32, enc.EncodeFloat, dec.DecodeFloat)
	rg.RegisterKindCodec(reflect.Float64, enc.EncodeFloat, dec.DecodeFloat)
	rg.RegisterKindCodec(reflect.Ptr, enc.EncodePtr, dec.DecodePtr)
	rg.RegisterKindCodec(reflect.Struct, enc.EncodeStruct, dec.DecodeStruct)
	rg.RegisterKindCodec(reflect.Array, enc.EncodeArray, dec.DecodeArray)
	rg.RegisterKindCodec(reflect.Slice, enc.EncodeArray, dec.DecodeArray)
	rg.RegisterKindCodec(reflect.Map, enc.EncodeMap, dec.DecodeMap)
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

// RegisterTypeCodec :
func (r *Registry) RegisterTypeCodec(t reflect.Type, enc ValueEncoder, dec ValueDecoder) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.typeEncoders[t] = enc
	r.typeDecoders[t] = dec
}

// RegisterTypeEncoder :
func (r *Registry) RegisterTypeEncoder(t reflect.Type, enc ValueEncoder) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.typeEncoders[t] = enc
}

// RegisterTypeDecoder :
func (r *Registry) RegisterTypeDecoder(t reflect.Type, dec ValueDecoder) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.typeDecoders[t] = dec
}

// RegisterKindCodec :
func (r *Registry) RegisterKindCodec(k reflect.Kind, enc ValueEncoder, dec ValueDecoder) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.kindEncoders[k] = enc
	r.kindDecoders[k] = dec
}

// RegisterKindEncoder :
func (r *Registry) RegisterKindEncoder(k reflect.Kind, enc ValueEncoder) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.kindEncoders[k] = enc
}

// RegisterKindDecoder :
func (r *Registry) RegisterKindDecoder(k reflect.Kind, dec ValueDecoder) {
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

func encodeValue(_ reflext.StructFielder, v reflect.Value) (interface{}, error) {
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
func NilEncoder(_ reflext.StructFielder, _ reflect.Value) (interface{}, error) {
	return nil, nil
}

func sqlScannerDecoder(it interface{}, v reflect.Value) error {
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
