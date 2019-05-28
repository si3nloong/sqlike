package codec

import (
	"database/sql"
	"encoding/base64"
	"reflect"
)

// DefaultDecoders :
type DefaultDecoders struct {
	registry *Registry
}

// SetDecoders :
func (dec DefaultDecoders) SetDecoders(rg *Registry) {
	rg.SetTypeDecoder(reflect.TypeOf([]byte{}), dec.DecodeByte)
	// rg.SetTypeDecoder(reflect.TypeOf(time.Time{}), dec.DecodeTime)
	// rg.SetTypeDecoder(reflect.TypeOf(json.RawMessage{}), dec.DecodeJSONRaw)
	// rg.SetKindDecoder(reflect.String, dec.DecodeString)
	// rg.SetKindDecoder(reflect.Bool, dec.DecodeBool)
	// rg.SetKindDecoder(reflect.Int, dec.DecodeInt)
	// rg.SetKindDecoder(reflect.Int8, dec.DecodeInt)
	// rg.SetKindDecoder(reflect.Int16, dec.DecodeInt)
	// rg.SetKindDecoder(reflect.Int32, dec.DecodeInt)
	// rg.SetKindDecoder(reflect.Int64, dec.DecodeInt)
	// rg.SetKindDecoder(reflect.Uint, dec.DecodeUint)
	// rg.SetKindDecoder(reflect.Uint8, dec.DecodeUint)
	// rg.SetKindDecoder(reflect.Uint16, dec.DecodeUint)
	// rg.SetKindDecoder(reflect.Uint32, dec.DecodeUint)
	// rg.SetKindDecoder(reflect.Uint64, dec.DecodeUint)
	// rg.SetKindDecoder(reflect.Float32, dec.DecodeFloat)
	// rg.SetKindDecoder(reflect.Float64, dec.DecodeFloat)
	// rg.SetKindDecoder(reflect.Ptr, dec.DecodePtr)
	// rg.SetKindDecoder(reflect.Struct, dec.DecodeStruct)
	// rg.SetKindDecoder(reflect.Array, dec.DecodeArray)
	// rg.SetKindDecoder(reflect.Slice, dec.DecodeArray)
	// rg.SetKindDecoder(reflect.Map, dec.DecodeMap)
	dec.registry = rg
}

// DecodeByte :
func (dec DefaultDecoders) DecodeByte(it interface{}, v reflect.Value) error {
	b := it.(sql.RawBytes)
	x := make([]byte, len(b), len(b))
	if _, err := base64.StdEncoding.Decode(x, b); err != nil {
		return err
	}
	v.SetBytes(x)
	return nil
}

// // DecodeJSONRaw :
// func (dec DefaultDecoders) DecodeJSONRaw(b []byte, v reflect.Value) error {
// 	v.SetBytes(b)
// 	return nil
// }

// // DecodeTime :
// func (dec DefaultDecoders) DecodeTime(b []byte, v reflect.Value) error {
// 	// format := time.RFC3339
// 	x, err := time.Parse(timeFormat, b2s(b))
// 	if err != nil {
// 		return err
// 	}
// 	v.Set(reflect.ValueOf(x))
// 	return nil
// }

// // DecodeString :
// func (dec DefaultDecoders) DecodeString(b []byte, v reflect.Value) error {
// 	v.SetString(b2s(b))
// 	return nil
// }

// // DecodeBool :
// func (dec DefaultDecoders) DecodeBool(b []byte, v reflect.Value) error {
// 	x, err := strconv.ParseBool(b2s(b))
// 	if err != nil {
// 		return err
// 	}
// 	v.SetBool(x)
// 	return nil
// }

// // DecodeInt :
// func (dec DefaultDecoders) DecodeInt(b []byte, v reflect.Value) error {
// 	x, err := strconv.ParseInt(b2s(b), 10, 64)
// 	if err != nil {
// 		return err
// 	}
// 	if v.OverflowInt(x) {
// 		return xerrors.New("integer overflow")
// 	}
// 	v.SetInt(x)
// 	return nil
// }

// // DecodeUint :
// func (dec DefaultDecoders) DecodeUint(b []byte, v reflect.Value) error {
// 	x, err := strconv.ParseUint(b2s(b), 10, 64)
// 	if err != nil {
// 		return err
// 	}
// 	if v.OverflowUint(x) {
// 		return xerrors.New("unsigned integer overflow")
// 	}
// 	v.SetUint(x)
// 	return nil
// }

// // DecodeFloat :
// func (dec DefaultDecoders) DecodeFloat(b []byte, v reflect.Value) error {
// 	x, err := strconv.ParseFloat(b2s(b), 64)
// 	if err != nil {
// 		return err
// 	}
// 	if v.OverflowFloat(x) {
// 		return xerrors.New("float overflow")
// 	}
// 	v.SetFloat(x)
// 	return nil
// }

// DecodePtr :
func (dec *DefaultDecoders) DecodePtr(b []byte, v reflect.Value) error {
	t := v.Type()
	if b == nil {
		v.Set(reflect.Zero(t))
		return nil
	}
	t = t.Elem()
	decoder, err := dec.registry.LookupDecoder(t)
	if err != nil {
		return err
	}
	return decoder(b, v.Elem())
}

// DecodeStruct :
func (dec DefaultDecoders) DecodeStruct(b []byte, v reflect.Value) error {
	// v.SetBytes(b)
	return nil
}

// DecodeArray :
func (dec DefaultDecoders) DecodeArray(b []byte, v reflect.Value) error {
	// v.SetBytes(b)
	return nil
}
