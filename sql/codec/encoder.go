package codec

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/si3nloong/sqlike/reflext"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"

	"github.com/si3nloong/sqlike/jsonb"
)

// DefaultEncoders :
type DefaultEncoders struct {
	registry *Registry
}

// SetEncoders :
func (enc DefaultEncoders) SetEncoders(rg *Registry) {
	rg.SetTypeEncoder(reflect.TypeOf([]byte{}), enc.EncodeByte)
	rg.SetTypeEncoder(reflect.TypeOf(language.Tag{}), enc.EncodeStringer)
	rg.SetTypeEncoder(reflect.TypeOf(currency.Unit{}), enc.EncodeStringer)
	rg.SetTypeEncoder(reflect.TypeOf(time.Time{}), enc.EncodeTime)
	rg.SetTypeEncoder(reflect.TypeOf(sql.RawBytes{}), enc.EncodeRawBytes)
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

// EncodeRawBytes :
func (enc DefaultEncoders) EncodeRawBytes(_ *reflext.StructField, v reflect.Value) (interface{}, error) {
	return sql.RawBytes(v.Bytes()), nil
}

// EncodeJSONRaw :
func (enc DefaultEncoders) EncodeJSONRaw(_ *reflext.StructField, v reflect.Value) (interface{}, error) {
	if v.IsNil() {
		return []byte("null"), nil
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

// EncodeStringer :
func (enc DefaultEncoders) EncodeStringer(_ *reflext.StructField, v reflect.Value) (interface{}, error) {
	x := v.Interface().(fmt.Stringer)
	return x.String(), nil
}

// EncodeTime :
func (enc DefaultEncoders) EncodeTime(_ *reflext.StructField, v reflect.Value) (interface{}, error) {
	x := v.Interface().(time.Time)
	// if x.IsZero() {
	// 	x, _ = time.Parse(time.RFC3339, "1970-01-01T08:00:00Z")
	// 	return x, nil
	// }
	// convert to UTC before storing into DB
	return x.UTC(), nil
}

// EncodeString :
func (enc DefaultEncoders) EncodeString(sf *reflext.StructField, v reflect.Value) (interface{}, error) {
	str := v.String()
	if str == "" && sf != nil {
		if val, ok := sf.Tag.LookUp("enum"); ok {
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
	if !v.IsValid() || v.IsNil() {
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
	t := v.Type()
	k := t.Key()
	if k.Kind() != reflect.String {
		return nil, fmt.Errorf("sqlike/sql/codec: unsupported data type %q for map key, it must be string", k.Kind())
	}
	k = t.Elem()
	if !isBaseType(k) {
		return nil, fmt.Errorf("sqlike/sql/codec: unsupported data type %q for map value", k.Kind())
	}
	if v.IsNil() {
		return string("null"), nil
	}
	return jsonb.Marshal(v)
}

func isBaseType(t reflect.Type) bool {
	for {
		k := t.Kind()
		switch k {
		case reflect.String:
			return true
		case reflect.Bool:
			return true
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return true
		case reflect.Float32, reflect.Float64:
			return true
		case reflect.Ptr:
			t = t.Elem()
		default:
			return false
		}
	}
}
