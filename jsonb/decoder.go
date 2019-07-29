package jsonb

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/si3nloong/sqlike/core"
	"github.com/si3nloong/sqlike/reflext"
	"golang.org/x/xerrors"
)

// Decoder :
type Decoder struct {
	registry *Registry
}

// SetDecoders :
func (dec Decoder) SetDecoders(rg *Registry) {
	rg.SetTypeDecoder(reflect.TypeOf([]byte{}), dec.DecodeByte)
	rg.SetTypeDecoder(reflect.TypeOf(time.Time{}), dec.DecodeTime)
	rg.SetTypeDecoder(reflect.TypeOf(json.RawMessage{}), dec.DecodeJSONRaw)
	rg.SetTypeDecoder(reflect.TypeOf(json.Number("")), dec.DecodeJSONNumber)
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
	rg.SetKindDecoder(reflect.Slice, dec.DecodeSlice)
	rg.SetKindDecoder(reflect.Map, dec.DecodeMap)
	rg.SetKindDecoder(reflect.Interface, dec.DecodeInterface)
	dec.registry = rg
}

// DecodeByte :
func (dec Decoder) DecodeByte(r *Reader, v reflect.Value) error {
	x, err := r.ReadRawString()
	if err != nil {
		return err
	}
	var b []byte
	if x == null {
		v.SetBytes(b)
		return nil
	} else if x == "" {
		v.SetBytes(make([]byte, 0))
		return nil
	}
	b, err = base64.StdEncoding.DecodeString(x)
	if err != nil {
		return err
	}
	v.SetBytes(b)
	return nil
}

// DecodeTime :
func (dec Decoder) DecodeTime(r *Reader, v reflect.Value) error {
	b, err := r.ReadBytes()
	if err != nil {
		return err
	}
	str := string(b)
	if str == null || str == `""` {
		v.Set(reflect.ValueOf(time.Time{}))
		return nil
	}
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return errors.New("jsonb: invalid format of date")
	}
	str = string(b[1 : len(b)-1])
	var x time.Time
	switch {
	case regexp.MustCompile(`^\d{4}\-\d{2}\-\d{2}$`).MatchString(str):
		x, err = time.Parse(`2006-01-02`, str)
	case regexp.MustCompile(`^\d{4}\-\d{2}\-\d{2}\s\d{2}\:\d{2}:\d{2}$`).MatchString(str):
		x, err = time.Parse(`2006-01-02 15:04:05`, str)
	default:
		x, err = time.Parse(time.RFC3339Nano, str)
	}
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

// DecodeJSONNumber :
func (dec Decoder) DecodeJSONNumber(r *Reader, v reflect.Value) error {
	x, err := r.ReadNumber()
	if err != nil {
		return err
	}
	v.SetString(x.String())
	return nil
}

// DecodeString :
func (dec Decoder) DecodeString(r *Reader, v reflect.Value) error {
	x, err := r.ReadEscapeString()
	if err != nil {
		return err
	}
	v.SetString(x)
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
	num, err := r.ReadNumber()
	if err != nil {
		return err
	}
	x, err := strconv.ParseInt(num.String(), 10, 64)
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
func (dec Decoder) DecodeUint(r *Reader, v reflect.Value) error {
	num, err := r.ReadNumber()
	if err != nil {
		return err
	}
	x, err := strconv.ParseUint(num.String(), 10, 64)
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
	num, err := r.ReadNumber()
	if err != nil {
		return err
	}
	x, err := strconv.ParseFloat(num.String(), 64)
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
func (dec *Decoder) DecodePtr(r *Reader, v reflect.Value) error {
	t := v.Type()
	if r.IsNull() {
		v.Set(reflect.Zero(t))
		return r.skipNull()
	}

	t = t.Elem()
	decoder, err := dec.registry.LookupDecoder(t)
	if err != nil {
		return err
	}
	vv := reflect.New(t)
	if err := decoder(r, vv.Elem()); err != nil {
		return err
	}
	v.Set(vv)
	return nil
}

// DecodeStruct :
func (dec *Decoder) DecodeStruct(r *Reader, v reflect.Value) error {
	mapper := core.DefaultMapper
	if r.IsNull() {
		v.Set(reflect.Zero(v.Type()))
		return r.skipNull()
	}

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

// DecodeArray :
func (dec *Decoder) DecodeArray(r *Reader, v reflect.Value) error {
	t := v.Type()
	if r.IsNull() {
		v.Set(reflect.Zero(t))
		return r.skipNull()
	}

	i, length := 0, v.Len()
	t = t.Elem()
	if err := r.ReadArray(func(it *Reader) error {
		if i >= length {
			return errors.New("jsonb: invalid array length")
		}
		vv := v.Index(i)
		i++
		decoder, err := dec.registry.LookupDecoder(t)
		if err != nil {
			return err
		}
		return decoder(it, vv)
	}); err != nil {
		return err
	}
	return nil
}

// DecodeSlice :
func (dec *Decoder) DecodeSlice(r *Reader, v reflect.Value) error {
	t := v.Type()
	if r.IsNull() {
		v.Set(reflect.Zero(t))
		return r.skipNull()
	}
	v.Set(reflect.MakeSlice(t, 0, 0))
	t = t.Elem()
	return r.ReadArray(func(it *Reader) error {
		v.Set(reflect.Append(v, reflext.Zero(t)))
		vv := v.Index(v.Len() - 1)
		decoder, err := dec.registry.LookupDecoder(t)
		if err != nil {
			return err
		}
		return decoder(it, vv)
	})
}

// DecodeMap :
func (dec *Decoder) DecodeMap(r *Reader, v reflect.Value) error {
	if r.IsNull() {
		v.Set(reflect.Zero(v.Type()))
		return r.skipNull()
	}
	t := v.Type()
	if t.Key().Kind() != reflect.String {
		return fmt.Errorf("jsonb: unsupported data type of map key, %q", t.Key().Kind())
	}
	decoder, err := dec.registry.LookupDecoder(t.Elem())
	if err != nil {
		return err
	}
	x := reflect.MakeMap(t)
	if err := r.ReadObject(func(it *Reader, k string) error {
		vi := reflext.Zero(t.Elem())
		err = decoder(it, vi)
		if err != nil {
			return err
		}
		x.SetMapIndex(reflect.ValueOf(k), vi)
		return nil
	}); err != nil {
		return err
	}
	v.Set(x)
	return nil
}

// DecodeInterface :
func (dec Decoder) DecodeInterface(r *Reader, v reflect.Value) error {
	x, err := r.ReadValue()
	if err != nil {
		return err
	}
	if x != nil {
		v.Set(reflect.ValueOf(x))
	}
	return nil
}
