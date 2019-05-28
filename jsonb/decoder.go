package jsonb

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"reflect"
	"strconv"

	"bitbucket.org/SianLoong/sqlike/core/codec"
	"bitbucket.org/SianLoong/sqlike/core"
	"bitbucket.org/SianLoong/sqlike/util"
	"golang.org/x/xerrors"
)

// ValueDecoder :
type ValueDecoder struct {
	registry *codec.Registry
}

// SetDecoders :
func (dec ValueDecoder) SetDecoders(rg *codec.Registry) {
	rg.SetTypeDecoder(reflect.TypeOf([]byte{}), dec.DecodeByte)
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
	rg.SetKindDecoder(reflect.Struct, dec.DecodeStruct)
	dec.registry = rg
}

// DecodeByte :
func (dec ValueDecoder) DecodeByte(r codec.ValueReader, v reflect.Value) error {
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

// DecodeJSONRaw :
func (dec ValueDecoder) DecodeJSONRaw(r codec.ValueReader, v reflect.Value) error {
	buf := new(bytes.Buffer)
	if err := json.Compact(buf, r.Bytes()); err != nil {
		return err
	}
	v.SetBytes(buf.Bytes())
	return nil
}

// DecodeString :
func (dec ValueDecoder) DecodeString(r codec.ValueReader, v reflect.Value) error {
	str := r.String()
	length := len(str)
	if str[0] != '"' && str[length-1] != '"' {
		return xerrors.New("invalid string format")
	}
	blr := util.AcquireString()
	defer util.ReleaseString(blr)
	unescapeString(blr, str[1:length-1])
	v.SetString(blr.String())
	return nil
}

// DecodeBool :
func (dec ValueDecoder) DecodeBool(r codec.ValueReader, v reflect.Value) error {
	// x, err := strconv.ParseBool(b2s(b))
	// if err != nil {
	// 	return err
	// }
	// v.SetBool(x)
	return nil
}

// DecodeInt :
func (dec ValueDecoder) DecodeInt(r codec.ValueReader, v reflect.Value) error {
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
func (dec ValueDecoder) DecodeUint(r codec.ValueReader, v reflect.Value) error {
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
func (dec ValueDecoder) DecodeFloat(r codec.ValueReader, v reflect.Value) error {
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

// DecodeStruct :
func (dec ValueDecoder) DecodeStruct(r codec.ValueReader, v reflect.Value) error {
	mapper := core.DefaultMapper
	cdc := mapper.CodecByType(v.Type())
	log.Println(cdc)
	return nil
}
