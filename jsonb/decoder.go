package jsonb

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/si3nloong/sqlike/core"
	"github.com/si3nloong/sqlike/core/codec"
	"golang.org/x/xerrors"
)

// Decoder :
type Decoder struct {
	registry *Registry
}

// SetDecoders :
func (dec Decoder) SetDecoders(rg *Registry) {
	// rg.SetTypeDecoder(reflect.TypeOf([]byte{}), dec.DecodeByte)
	rg.SetTypeDecoder(reflect.TypeOf(json.RawMessage{}), dec.DecodeJSONRaw)
	rg.SetKindDecoder(reflect.String, dec.DecodeString)
	rg.SetKindDecoder(reflect.Bool, dec.DecodeBool)
	rg.SetKindDecoder(reflect.Int, dec.DecodeInt)
	rg.SetKindDecoder(reflect.Int8, dec.DecodeInt)
	rg.SetKindDecoder(reflect.Int16, dec.DecodeInt)
	rg.SetKindDecoder(reflect.Int32, dec.DecodeInt)
	rg.SetKindDecoder(reflect.Int64, dec.DecodeInt)
	// rg.SetKindDecoder(reflect.Uint, dec.DecodeUint)
	// rg.SetKindDecoder(reflect.Uint8, dec.DecodeUint)
	// rg.SetKindDecoder(reflect.Uint16, dec.DecodeUint)
	// rg.SetKindDecoder(reflect.Uint32, dec.DecodeUint)
	// rg.SetKindDecoder(reflect.Uint64, dec.DecodeUint)
	rg.SetKindDecoder(reflect.Float32, dec.DecodeFloat)
	rg.SetKindDecoder(reflect.Float64, dec.DecodeFloat)
	rg.SetKindDecoder(reflect.Struct, dec.DecodeStruct)
	rg.SetKindDecoder(reflect.Slice, dec.DecodeSlice)
	rg.SetKindDecoder(reflect.Interface, dec.DecodeInterface)
	// rg.SetKindDecoder(reflect.Array, dec.DecodeArray)
	dec.registry = rg
}

// DecodeByte :
func (dec Decoder) DecodeByte(r codec.ValueReader, v reflect.Value) error {
	b := r.Bytes()
	if b == nil {
		return nil
	}
	length := len(b)
	if b[0] != '"' && b[length-1] != '"' {
		return xerrors.New("invalid byte format")
	}
	var err error
	b, err = base64.StdEncoding.DecodeString(r.String()[1 : length-1])
	if err != nil {
		return err
	}
	v.SetBytes(b)
	return nil
}

// DecodeTime :
func (dec Decoder) DecodeTime(r *Reader, v reflect.Value) error {
	x, err := time.Parse(time.RFC3339Nano, string(r.Bytes()))
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(x))
	return nil
}

// DecodeJSONRaw :
func (dec Decoder) DecodeJSONRaw(r *Reader, v reflect.Value) error {
	v.SetBytes(r.Bytes())
	return nil
}

// DecodeString :
func (dec Decoder) DecodeString(r *Reader, v reflect.Value) error {
	v.SetString(r.ReadString())
	return nil
}

// DecodeBool :
func (dec Decoder) DecodeBool(r *Reader, v reflect.Value) error {
	x, err := r.ReadBoolean()
	if err != nil {
		return err
	}
	v.SetBool(x)
	return nil
}

// DecodeInt :
func (dec Decoder) DecodeInt(r *Reader, v reflect.Value) error {
	x, err := strconv.ParseInt(string(r.Bytes()), 10, 64)
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
func (dec Decoder) DecodeUint(r codec.ValueReader, v reflect.Value) error {
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
func (dec Decoder) DecodeFloat(r *Reader, v reflect.Value) error {
	x, err := strconv.ParseFloat(string(r.Bytes()), 64)
	if err != nil {
		return err
	}
	if v.OverflowFloat(x) {
		return xerrors.New("float overflow")
	}
	v.SetFloat(x)
	return nil
}

// DecodeStruct :
func (dec *Decoder) DecodeStruct(r *Reader, v reflect.Value) error {
	mapper := core.DefaultMapper
	return r.ReadFlattenObject(func(it *Reader, k string) error {
		vv, exists := mapper.LookUpFieldByName(v, k)
		if !exists {
			return nil
		}
		decoder, err := dec.registry.LookupDecoder(vv.Type())
		if err != nil {
			return err
		}
		return decoder(it, vv)
	})
}

// DecodeSlice :
func (dec Decoder) DecodeSlice(r *Reader, v reflect.Value) error {
	// tkn := r.nextToken()
	// if tkn.typ != jsonArray {
	// 	return xerrors.New("expected array")
	// }

	// reflect.MakeSlice(v.Type().Elem(), 0, 0)
	// // decoder := dec.registry.LookupDecoder(v.Type().Elem())
	// log.Println(v.Type().Elem())
	// log.Println(tkn.typ.String())
	// log.Println(r)
	return nil
}

// DecodeInterface :
func (dec Decoder) DecodeInterface(r *Reader, v reflect.Value) error {
	x, err := r.ReadValue()
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(x))
	return nil
}
