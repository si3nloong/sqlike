package codec

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"golang.org/x/xerrors"
)

// DefaultDecoders :
type DefaultDecoders struct {
	registry *Registry
}

// SetDecoders :
func (dec DefaultDecoders) SetDecoders(rg *Registry) {
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
func (dec DefaultDecoders) DecodeByte(it interface{}, v reflect.Value) error {
	var (
		x   []byte
		err error
	)
	switch vi := it.(type) {
	case string:
		x, err = base64.StdEncoding.DecodeString(vi)
		if err != nil {
			return err
		}
	case []byte:
		x, err = base64.StdEncoding.DecodeString(b2s(vi))
		if err != nil {
			return err
		}
	case nil:
		x = make([]byte, 0, 0)
	}
	v.SetBytes(x)
	return nil
}

// DecodeJSONRaw :
func (dec DefaultDecoders) DecodeJSONRaw(it interface{}, v reflect.Value) error {
	b := new(bytes.Buffer)
	switch vi := it.(type) {
	case string:
		if err := json.Compact(b, []byte(vi)); err != nil {
			return err
		}
	case []byte:
		if err := json.Compact(b, vi); err != nil {
			return err
		}
	case nil:
	}
	v.SetBytes(b.Bytes())
	return nil
}

// DecodeTime :
func (dec DefaultDecoders) DecodeTime(it interface{}, v reflect.Value) error {
	var (
		x   time.Time
		err error
	)
	switch vi := it.(type) {
	case time.Time:
		x = vi
	case string:
		x, err = time.Parse(time.RFC3339, vi)
		if err != nil {
			return err
		}
	case []byte:
		x, err = time.Parse(time.RFC3339, b2s(vi))
		if err != nil {
			return err
		}
	case nil:
	}
	// convert back to UTC
	v.Set(reflect.ValueOf(x.UTC()))
	return nil
}

// DecodeString :
func (dec DefaultDecoders) DecodeString(it interface{}, v reflect.Value) error {
	var x string
	switch vi := it.(type) {
	case string:
		x = vi
	case []byte:
		x = string(vi)
	case nil:
	}
	v.SetString(x)
	return nil
}

// DecodeBool :
func (dec DefaultDecoders) DecodeBool(it interface{}, v reflect.Value) error {
	var (
		x   bool
		err error
	)
	switch vi := it.(type) {
	case []byte:
		x, err = strconv.ParseBool(b2s(vi))
		if err != nil {
			return err
		}
	case string:
		x, err = strconv.ParseBool(vi)
		if err != nil {
			return err
		}
	case bool:
		x = vi
	case nil:
	}
	v.SetBool(x)
	return nil
}

// DecodeInt :
func (dec DefaultDecoders) DecodeInt(it interface{}, v reflect.Value) error {
	var (
		x   int64
		err error
	)
	switch vi := it.(type) {
	case []byte:
		x, err = strconv.ParseInt(b2s(vi), 10, 64)
		if err != nil {
			return err
		}
	case string:
		x, err = strconv.ParseInt(vi, 10, 64)
		if err != nil {
			return err
		}
	case int64:
		x = vi
	case nil:
	}
	if v.OverflowInt(x) {
		return xerrors.New("integer overflow")
	}
	v.SetInt(x)
	return nil
}

// DecodeUint :
func (dec DefaultDecoders) DecodeUint(it interface{}, v reflect.Value) error {
	var (
		x   uint64
		err error
	)
	switch vi := it.(type) {
	case []byte:
		x, err = strconv.ParseUint(b2s(vi), 10, 64)
		if err != nil {
			return err
		}
	case string:
		x, err = strconv.ParseUint(vi, 10, 64)
		if err != nil {
			return err
		}
	case uint64:
		x = vi
	case nil:
	}
	if v.OverflowUint(x) {
		return xerrors.New("unsigned integer overflow")
	}
	v.SetUint(x)
	return nil
}

// DecodeFloat :
func (dec DefaultDecoders) DecodeFloat(it interface{}, v reflect.Value) error {
	var (
		x   float64
		err error
	)
	switch vi := it.(type) {
	case []byte:
		x, err = strconv.ParseFloat(b2s(vi), 64)
		if err != nil {
			return err
		}
	case string:
		x, err = strconv.ParseFloat(vi, 64)
		if err != nil {
			return err
		}
	case float64:
		x = vi
	case nil:

	}
	if v.OverflowFloat(x) {
		return xerrors.New("float overflow")
	}
	v.SetFloat(x)
	return nil
}

// DecodePtr :
func (dec *DefaultDecoders) DecodePtr(it interface{}, v reflect.Value) error {
	t := v.Type()
	if it == nil {
		v.Set(reflect.Zero(t))
		return nil
	}
	t = t.Elem()
	decoder, err := dec.registry.LookupDecoder(t)
	if err != nil {
		return err
	}
	return decoder(it, v.Elem())
}

// DecodeStruct :
func (dec DefaultDecoders) DecodeStruct(it interface{}, v reflect.Value) error {
	switch vi := it.(type) {
	case string:
		return json.Unmarshal([]byte(vi), v.Interface())
	case []byte:
		return json.Unmarshal(vi, v.Interface())
	}
	return nil
}

// DecodeArray :
func (dec DefaultDecoders) DecodeArray(it interface{}, v reflect.Value) error {
	switch vi := it.(type) {
	case string:
		return json.Unmarshal([]byte(vi), v.Interface())
	case []byte:
		return json.Unmarshal(vi, v.Interface())
	}
	return nil
}
