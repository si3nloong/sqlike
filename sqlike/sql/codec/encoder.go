package codec

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"time"

	"github.com/si3nloong/sqlike/jsonb"
)

// DefaultEncoders :
type DefaultEncoders struct {
	registry *Registry
}

// SetEncoders :
func (enc DefaultEncoders) SetEncoders(rg *Registry) {
	rg.SetTypeEncoder(reflect.TypeOf([]byte{}), enc.EncodeByte)
	rg.SetTypeEncoder(reflect.TypeOf(time.Time{}), enc.EncodeTime)
	rg.SetTypeEncoder(reflect.TypeOf(json.RawMessage{}), enc.EncodeJSONRaw)
	rg.SetKindEncoder(reflect.String, enc.EncodeString)
	rg.SetKindEncoder(reflect.Bool, enc.EncodeBool)
	rg.SetKindEncoder(reflect.Int, enc.EncodeInt)
	rg.SetKindEncoder(reflect.Int8, enc.EncodeInt)
	rg.SetKindEncoder(reflect.Int16, enc.EncodeInt)
	rg.SetKindEncoder(reflect.Int32, enc.EncodeInt)
	rg.SetKindEncoder(reflect.Int64, enc.EncodeInt)
	rg.SetKindEncoder(reflect.Uint, enc.EncodeUint)
	rg.SetKindEncoder(reflect.Uint8, enc.EncodeUint)
	rg.SetKindEncoder(reflect.Uint16, enc.EncodeUint)
	rg.SetKindEncoder(reflect.Uint32, enc.EncodeUint)
	rg.SetKindEncoder(reflect.Uint64, enc.EncodeUint)
	rg.SetKindEncoder(reflect.Float32, enc.EncodeFloat)
	rg.SetKindEncoder(reflect.Float64, enc.EncodeFloat)
	rg.SetKindEncoder(reflect.Ptr, enc.EncodePtr)
	rg.SetKindEncoder(reflect.Struct, enc.EncodeStruct)
	rg.SetKindEncoder(reflect.Array, enc.EncodeArray)
	rg.SetKindEncoder(reflect.Slice, enc.EncodeArray)
	rg.SetKindEncoder(reflect.Map, enc.EncodeMap)
	enc.registry = rg
}

// EncodeByte :
func (enc DefaultEncoders) EncodeByte(v reflect.Value) (interface{}, error) {
	buf := new(bytes.Buffer)
	b64 := base64.NewEncoder(base64.StdEncoding, buf)
	defer b64.Close()
	b64.Write(v.Bytes())
	b := buf.Bytes()
	if b == nil {
		return make([]byte, 0, 0), nil
	}
	return b, nil
}

// EncodeJSONRaw :
func (enc DefaultEncoders) EncodeJSONRaw(v reflect.Value) (interface{}, error) {
	if v.IsNil() {
		return []byte(`null`), nil
	}
	buf := new(bytes.Buffer)
	if err := json.Compact(buf, v.Bytes()); err != nil {
		return nil, err
	}
	if buf.Len() == 0 {
		return []byte(`{}`), nil
	}
	return json.RawMessage(buf.Bytes()), nil
}

// EncodeTime :
func (enc DefaultEncoders) EncodeTime(v reflect.Value) (interface{}, error) {
	x, isOk := v.Interface().(time.Time)
	if !isOk {

	}
	return x, nil
}

// EncodeString :
func (enc DefaultEncoders) EncodeString(v reflect.Value) (interface{}, error) {
	return v.String(), nil
}

// EncodeBool :
func (enc DefaultEncoders) EncodeBool(v reflect.Value) (interface{}, error) {
	return v.Bool(), nil
}

// EncodeInt :
func (enc DefaultEncoders) EncodeInt(v reflect.Value) (interface{}, error) {
	return v.Int(), nil
}

// EncodeUint :
func (enc DefaultEncoders) EncodeUint(v reflect.Value) (interface{}, error) {
	return v.Uint(), nil
}

// EncodeFloat :
func (enc DefaultEncoders) EncodeFloat(v reflect.Value) (interface{}, error) {
	return v.Float(), nil
}

// EncodePtr :
func (enc *DefaultEncoders) EncodePtr(v reflect.Value) (interface{}, error) {
	if v.IsNil() {
		return nil, nil
	}
	v = v.Elem()
	encoder, err := enc.registry.LookupEncoder(v.Type())
	if err != nil {
		return nil, err
	}
	return encoder(v)
}

// EncodeStruct :
func (enc DefaultEncoders) EncodeStruct(v reflect.Value) (interface{}, error) {
	return jsonb.Marshal(v)
}

// EncodeArray :
func (enc DefaultEncoders) EncodeArray(v reflect.Value) (interface{}, error) {
	return jsonb.Marshal(v)
}

// EncodeMap :
func (enc DefaultEncoders) EncodeMap(v reflect.Value) (interface{}, error) {
	return jsonb.Marshal(v)
}
