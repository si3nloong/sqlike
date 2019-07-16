package codec

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/si3nloong/sqlike/reflext"

	"github.com/si3nloong/sqlike/jsonb"
	"golang.org/x/xerrors"
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
func (enc DefaultEncoders) EncodeByte(_ *reflext.StructField, v reflect.Value) (interface{}, error) {
	b := v.Bytes()
	if b == nil {
		return make([]byte, 0, 0), nil
	}
	x := base64.StdEncoding.EncodeToString(b)
	return []byte(x), nil
}

// EncodeJSONRaw :
func (enc DefaultEncoders) EncodeJSONRaw(_ *reflext.StructField, v reflect.Value) (interface{}, error) {
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
func (enc DefaultEncoders) EncodeTime(_ *reflext.StructField, v reflect.Value) (interface{}, error) {
	x, isOk := v.Interface().(time.Time)
	if !isOk {
		return nil, xerrors.New("invalid data type")
	}
	// convert to UTC before storing into DB
	return x.UTC(), nil
}

// EncodeString :
func (enc DefaultEncoders) EncodeString(sf *reflext.StructField, v reflect.Value) (interface{}, error) {
	str := v.String()
	if str == "" && sf != nil {
		if val, isOk := sf.Tag.LookUp("enum"); isOk {
			enums := strings.Split(val, "|")
			if len(enums) > 0 {
				return enums[0], nil
			}
		}
	}
	return str, nil
}

// EncodeBool :
func (enc DefaultEncoders) EncodeBool(_ *reflext.StructField, v reflect.Value) (interface{}, error) {
	return v.Bool(), nil
}

// EncodeInt :
func (enc DefaultEncoders) EncodeInt(_ *reflext.StructField, v reflect.Value) (interface{}, error) {
	return v.Int(), nil
}

// EncodeUint :
func (enc DefaultEncoders) EncodeUint(_ *reflext.StructField, v reflect.Value) (interface{}, error) {
	return v.Uint(), nil
}

// EncodeFloat :
func (enc DefaultEncoders) EncodeFloat(_ *reflext.StructField, v reflect.Value) (interface{}, error) {
	return v.Float(), nil
}

// EncodePtr :
func (enc *DefaultEncoders) EncodePtr(sf *reflext.StructField, v reflect.Value) (interface{}, error) {
	if v.IsNil() {
		return nil, nil
	}
	v = v.Elem()
	encoder, err := enc.registry.LookupEncoder(v.Type())
	if err != nil {
		return nil, err
	}
	return encoder(sf, v)
}

// EncodeStruct :
func (enc DefaultEncoders) EncodeStruct(_ *reflext.StructField, v reflect.Value) (interface{}, error) {
	return jsonb.Marshal(v)
}

// EncodeArray :
func (enc DefaultEncoders) EncodeArray(_ *reflext.StructField, v reflect.Value) (interface{}, error) {
	return jsonb.Marshal(v)
}

// EncodeMap :
func (enc DefaultEncoders) EncodeMap(_ *reflext.StructField, v reflect.Value) (interface{}, error) {
	return jsonb.Marshal(v)
}
