package jsonb

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/si3nloong/sqlike/reflext"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

// Encoder :
type Encoder struct {
	registry *Registry
}

// SetEncoders :
func (enc Encoder) SetEncoders(rg *Registry) {
	rg.SetTypeEncoder(reflect.TypeOf([]byte{}), enc.EncodeByte)
	rg.SetTypeEncoder(reflect.TypeOf(language.Tag{}), enc.EncodeStringer)
	rg.SetTypeEncoder(reflect.TypeOf(currency.Unit{}), enc.EncodeStringer)
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
	rg.SetKindEncoder(reflect.Interface, enc.EncodeInterface)
	enc.registry = rg
}

// EncodeByte :
func (enc Encoder) EncodeByte(w *Writer, v reflect.Value) error {
	if v.IsNil() {
		w.WriteString(null)
		return nil
	}
	w.WriteRune('"')
	w.WriteString(base64.StdEncoding.EncodeToString(v.Bytes()))
	w.WriteRune('"')
	return nil
}

// EncodeStringer :
func (enc Encoder) EncodeStringer(w *Writer, v reflect.Value) error {
	x := v.Interface().(fmt.Stringer)
	w.WriteString(strconv.Quote(x.String()))
	return nil
}

// EncodeJSONRaw :
func (enc Encoder) EncodeJSONRaw(w *Writer, v reflect.Value) error {
	if v.IsNil() {
		w.WriteString(null)
		return nil
	}
	buf := new(bytes.Buffer)
	if err := json.Compact(buf, v.Bytes()); err != nil {
		return err
	}
	if buf.Len() == 0 {
		w.Write([]byte(`{}`))
		return nil
	}
	w.Write(buf.Bytes())
	return nil
}

// EncodeTime :
func (enc Encoder) EncodeTime(w *Writer, v reflect.Value) error {
	var temp [20]byte
	x := v.Interface().(time.Time)
	w.Write(x.UTC().AppendFormat(temp[:0], `"`+time.RFC3339Nano+`"`))
	return nil
}

// EncodeString :
func (enc Encoder) EncodeString(w *Writer, v reflect.Value) error {
	w.WriteRune('"')
	escapeString(w, v.String())
	w.WriteRune('"')
	return nil
}

// EncodeBool :
func (enc Encoder) EncodeBool(w *Writer, v reflect.Value) error {
	var temp [4]byte
	w.Write(strconv.AppendBool(temp[:0], v.Bool()))
	return nil
}

// EncodeInt :
func (enc Encoder) EncodeInt(w *Writer, v reflect.Value) error {
	var temp [8]byte
	w.Write(strconv.AppendInt(temp[:0], v.Int(), 10))
	return nil
}

// EncodeUint :
func (enc Encoder) EncodeUint(w *Writer, v reflect.Value) error {
	var temp [10]byte
	w.Write(strconv.AppendUint(temp[:0], v.Uint(), 10))
	return nil
}

// EncodeFloat :
func (enc Encoder) EncodeFloat(w *Writer, v reflect.Value) error {
	f64 := v.Float()
	if f64 <= 0 {
		w.WriteRune('0')
		return nil
	}
	w.WriteString(strconv.FormatFloat(f64, 'E', -1, 64))
	return nil
}

// EncodePtr :
func (enc *Encoder) EncodePtr(w *Writer, v reflect.Value) error {
	if v.IsNil() {
		w.WriteString(null)
		return nil
	}

	v = v.Elem()
	encoder, err := enc.registry.LookupEncoder(v)
	if err != nil {
		return err
	}
	return encoder(w, v)
}

// EncodeStruct :
func (enc *Encoder) EncodeStruct(w *Writer, v reflect.Value) error {
	w.WriteRune('{')
	mapper := reflext.DefaultMapper
	cdc := mapper.CodecByType(v.Type())
	for i, sf := range cdc.Properties {
		if i > 0 {
			w.WriteRune(',')
		}
		w.WriteString(strconv.Quote(sf.Path))
		w.WriteRune(':')
		fv := mapper.FieldByIndexesReadOnly(v, sf.Index)
		encoder, err := enc.registry.LookupEncoder(fv)
		if err != nil {
			return err
		}
		if err := encoder(w, fv); err != nil {
			return err
		}
	}
	w.WriteRune('}')
	return nil
}

// EncodeArray :
func (enc *Encoder) EncodeArray(w *Writer, v reflect.Value) error {
	if v.Kind() == reflect.Slice && v.IsNil() {
		w.WriteString(null)
		return nil
	}
	w.WriteRune('[')
	length := v.Len()
	for i := 0; i < length; i++ {
		if i > 0 {
			w.WriteRune(',')
		}

		fv := v.Index(i)
		encoder, err := enc.registry.LookupEncoder(fv)
		if err != nil {
			return err
		}
		if err := encoder(w, fv); err != nil {
			return err
		}
	}
	w.WriteRune(']')
	return nil
}

// EncodeInterface :
func (enc *Encoder) EncodeInterface(w *Writer, v reflect.Value) error {
	it := v.Interface()
	if it == nil {
		w.WriteString(null)
		return nil
	}
	encoder, err := enc.registry.LookupEncoder(v)
	if err != nil {
		return err
	}
	return encoder(w, reflect.ValueOf(it))
}

// EncodeMap :
func (enc *Encoder) EncodeMap(w *Writer, v reflect.Value) error {
	t := v.Type()
	k := t.Key()
	if k.Kind() != reflect.String {
		return fmt.Errorf("jsonb: unsupported data type %q for map key, it must be string", k.Kind())
	}
	if v.IsNil() {
		w.WriteString(null)
		return nil
	}
	w.WriteByte('{')
	if v.Len() == 0 {
		w.WriteByte('}')
		return nil
	}
	keys := reflext.MapKeys(v.MapKeys())
	sort.Sort(keys)
	length := len(keys)
	for i := 0; i < length; i++ {
		if i > 0 {
			w.WriteByte(',')
		}
		k := keys[i]
		vv := v.MapIndex(k)
		w.WriteRune('"')
		escapeString(w, k.String())
		w.WriteRune('"')
		w.WriteByte(':')
		encoder, err := enc.registry.LookupEncoder(vv)
		if err != nil {
			return err
		}
		if err := encoder(w, vv); err != nil {
			return err
		}
	}
	w.WriteByte('}')
	return nil
}
