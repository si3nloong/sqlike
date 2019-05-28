package core

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/si3nloong/sqlike/core/codec"
	"golang.org/x/xerrors"
)

const timeFormat = "2006-01-02 15:04:05.999999"

// DefaultDecoders :
type DefaultDecoders struct {
	registry *codec.Registry
}

// SetDecoders :
func (dec DefaultDecoders) SetDecoders(rg *codec.Registry) {
	rg.SetTypeDecoder(reflect.TypeOf([]byte{}), dec.DecodeByte)
	rg.SetTypeDecoder(reflect.TypeOf(time.Time{}), dec.DecodeTime)
	rg.SetTypeDecoder(reflect.TypeOf(json.RawMessage{}), dec.DecodeJSONRaw)
	rg.SetKindDecoder(reflect.String, dec.DecodeString)
	rg.SetKindDecoder(reflect.Bool, dec.DecodeBool)
	rg.SetKindDecoder(reflect.Int, dec.DecodeInt)
	rg.SetKindDecoder(reflect.Int8, dec.DecodeInt)
	rg.SetKindDecoder(reflect.Int16, dec.DecodeInt)
	rg.SetKindDecoder(reflect.Int32, dec.DecodeInt)
	rg.SetKindDecoder(reflect.Int64, dec.DecodeInt)
	rg.SetKindDecoder(reflect.Uint, dec.DecodeUint)
	rg.SetKindDecoder(reflect.Uint8, dec.DecodeUint)
	rg.SetKindDecoder(reflect.Uint16, dec.DecodeUint)
	rg.SetKindDecoder(reflect.Uint32, dec.DecodeUint)
	rg.SetKindDecoder(reflect.Uint64, dec.DecodeUint)
	rg.SetKindDecoder(reflect.Float32, dec.DecodeFloat)
	rg.SetKindDecoder(reflect.Float64, dec.DecodeFloat)
	rg.SetKindDecoder(reflect.Ptr, dec.DecodePtr)
	rg.SetKindDecoder(reflect.Struct, dec.DecodeStruct)
	rg.SetKindDecoder(reflect.Array, dec.DecodeArray)
	rg.SetKindDecoder(reflect.Slice, dec.DecodeArray)
	// rg.SetKindDecoder(reflect.Map, dec.DecodeMap)
	dec.registry = rg
}

// DecodeByte :
func (dec DefaultDecoders) DecodeByte(r codec.ValueReader, v reflect.Value) error {
	b := r.Bytes()
	x := make([]byte, len(b), len(b))
	if _, err := base64.StdEncoding.Decode(x, b); err != nil {
		return err
	}
	v.SetBytes(x)
	return nil
}

// DecodeJSONRaw :
func (dec DefaultDecoders) DecodeJSONRaw(r codec.ValueReader, v reflect.Value) error {
	v.SetBytes(r.Bytes())
	return nil
}

// DecodeTime :
func (dec DefaultDecoders) DecodeTime(r codec.ValueReader, v reflect.Value) error {
	x, err := time.Parse(timeFormat, r.String())
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(x))
	return nil
}

// DecodeString :
func (dec DefaultDecoders) DecodeString(r codec.ValueReader, v reflect.Value) error {
	v.SetString(r.String())
	return nil
}

// DecodeBool :
func (dec DefaultDecoders) DecodeBool(r codec.ValueReader, v reflect.Value) error {
	x, err := strconv.ParseBool(r.String())
	if err != nil {
		return err
	}
	v.SetBool(x)
	return nil
}

// DecodeInt :
func (dec DefaultDecoders) DecodeInt(r codec.ValueReader, v reflect.Value) error {
	x, err := strconv.ParseInt(r.String(), 10, 64)
	if err != nil {
		return err
	}
	if v.OverflowInt(x) {
		return xerrors.New("integer overflow")
	}
	v.SetInt(x)
	return nil
}

// DecodeUint :
func (dec DefaultDecoders) DecodeUint(r codec.ValueReader, v reflect.Value) error {
	x, err := strconv.ParseUint(r.String(), 10, 64)
	if err != nil {
		return err
	}
	if v.OverflowUint(x) {
		return xerrors.New("unsigned integer overflow")
	}
	v.SetUint(x)
	return nil
}

// DecodeFloat :
func (dec DefaultDecoders) DecodeFloat(r codec.ValueReader, v reflect.Value) error {
	x, err := strconv.ParseFloat(r.String(), 64)
	if err != nil {
		return err
	}
	if v.OverflowFloat(x) {
		return xerrors.New("float overflow")
	}
	v.SetFloat(x)
	return nil
}

// DecodePtr :
func (dec *DefaultDecoders) DecodePtr(r codec.ValueReader, v reflect.Value) error {
	t := v.Type()
	b := r.Bytes()
	if b == nil {
		v.Set(reflect.Zero(t))
		return nil
	}
	t = t.Elem()
	decoder, err := dec.registry.LookupDecoder(t)
	if err != nil {
		return err
	}
	return decoder(r, v.Elem())
}

// DecodeStruct :
func (dec DefaultDecoders) DecodeStruct(r codec.ValueReader, v reflect.Value) error {
	// return jsonb.Unmarshal(r.Bytes(), v)
	return nil
}

// DecodeArray :
func (dec DefaultDecoders) DecodeArray(r codec.ValueReader, v reflect.Value) error {
	// return jsonb.Unmarshal(r.Bytes(), v)
	return nil
}
