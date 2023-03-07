package jsonb

import (
	"bytes"
	"encoding"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/si3nloong/sqlike/v2/x/reflext"
)

// DefaultEncoder :
type DefaultEncoder struct {
	registry *Registry
}

// EncodeByte :
func (enc DefaultEncoder) EncodeByte(w JsonWriter, v reflect.Value) error {
	if v.IsNil() {
		w.WriteString(null)
		return nil
	}
	w.WriteString(`"` + base64.StdEncoding.EncodeToString(v.Bytes()) + `"`)
	return nil
}

// EncodeStringer :
func (enc DefaultEncoder) EncodeStringer(w JsonWriter, v reflect.Value) error {
	x := v.Interface().(fmt.Stringer)
	w.WriteString(strconv.Quote(x.String()))
	return nil
}

// EncodeJsonRaw :
func (enc DefaultEncoder) EncodeJsonRaw(w JsonWriter, v reflect.Value) error {
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
func (enc DefaultEncoder) EncodeTime(w JsonWriter, v reflect.Value) error {
	var temp [20]byte
	x := v.Interface().(time.Time)
	w.Write(x.UTC().AppendFormat(temp[:0], `"`+time.RFC3339Nano+`"`))
	return nil
}

// EncodeString :
func (enc DefaultEncoder) EncodeString(w JsonWriter, v reflect.Value) error {
	w.WriteByte('"')
	escapeString(w, v.String())
	w.WriteByte('"')
	return nil
}

// EncodeBool :
func (enc DefaultEncoder) EncodeBool(w JsonWriter, v reflect.Value) error {
	var temp [4]byte
	w.Write(strconv.AppendBool(temp[:0], v.Bool()))
	return nil
}

// EncodeInt :
func (enc DefaultEncoder) EncodeInt(w JsonWriter, v reflect.Value) error {
	var temp [8]byte
	w.Write(strconv.AppendInt(temp[:0], v.Int(), 10))
	return nil
}

// EncodeUint :
func (enc DefaultEncoder) EncodeUint(w JsonWriter, v reflect.Value) error {
	var temp [10]byte
	w.Write(strconv.AppendUint(temp[:0], v.Uint(), 10))
	return nil
}

// EncodeFloat :
func (enc DefaultEncoder) EncodeFloat(w JsonWriter, v reflect.Value) error {
	f64 := v.Float()
	if f64 <= 0 {
		w.WriteByte('0')
		return nil
	}
	w.WriteString(strconv.FormatFloat(f64, 'E', -1, 64))
	return nil
}

// EncodePtr :
func (enc *DefaultEncoder) EncodePtr(w JsonWriter, v reflect.Value) error {
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
func (enc *DefaultEncoder) EncodeStruct(w JsonWriter, v reflect.Value) error {
	w.WriteByte('{')
	mapper := reflext.DefaultMapper()
	cdc := mapper.CodecByType(v.Type())
	for i, sf := range cdc.Properties() {
		if i > 0 {
			w.WriteByte(',')
		}
		w.WriteString(strconv.Quote(sf.Name()) + `:`)
		fv := mapper.FieldByIndexesReadOnly(v, sf.Index())
		encoder, err := enc.registry.LookupEncoder(fv)
		if err != nil {
			return err
		}
		if err := encoder(w, fv); err != nil {
			return err
		}
	}
	w.WriteByte('}')
	return nil
}

// EncodeArray :
func (enc *DefaultEncoder) EncodeArray(w JsonWriter, v reflect.Value) error {
	if v.Kind() == reflect.Slice && v.IsNil() {
		w.WriteString(null)
		return nil
	}
	w.WriteByte('[')
	length := v.Len()
	for i := 0; i < length; i++ {
		if i > 0 {
			w.WriteByte(',')
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
	w.WriteByte(']')
	return nil
}

// EncodeInterface :
func (enc *DefaultEncoder) EncodeInterface(w JsonWriter, v reflect.Value) error {
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
func (enc *DefaultEncoder) EncodeMap(w JsonWriter, v reflect.Value) error {
	t := v.Type()
	k := t.Key()
	if v.IsNil() {
		w.WriteString(null)
		return nil
	}

	w.WriteByte('{')
	if v.Len() == 0 {
		w.WriteByte('}')
		return nil
	}

	keys := v.MapKeys()
	var encode ValueEncoder
	if k.Implements(textMarshaler) {
		encode = func(w JsonWriter, vi reflect.Value) error {
			it := vi.Interface().(encoding.TextMarshaler)
			b, err := it.MarshalText()
			if err != nil {
				return err
			}

			w.WriteByte('"')
			w.Write(b)
			w.WriteByte('"')
			return nil
		}
	} else {
		switch k.Kind() {
		case reflect.String:
			sort.SliceStable(keys, func(i, j int) bool {
				return keys[i].String() < keys[j].String()
			})
			encode = enc.registry.kindEncoders[reflect.String]
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			sort.SliceStable(keys, func(i, j int) bool {
				return keys[i].Int() < keys[j].Int()
			})
			encode = func(w JsonWriter, vi reflect.Value) error {
				w.WriteString(`"` + strconv.FormatInt(vi.Int(), 10) + `"`)
				return nil
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			sort.SliceStable(keys, func(i, j int) bool {
				return keys[i].Uint() < keys[j].Uint()
			})
			encode = func(w JsonWriter, vi reflect.Value) error {
				w.WriteString(`"` + strconv.FormatUint(vi.Uint(), 10) + `"`)
				return nil
			}
		default:
			return fmt.Errorf("jsonb: unsupported data type %q for map key, it must be string", k.Kind())
		}

	}

	// Question: do we really need to sort the key before encode?

	length := len(keys)
	for i := 0; i < length; i++ {
		if i > 0 {
			w.WriteByte(',')
		}
		k := keys[i]
		val := v.MapIndex(k)
		if err := encode(w, k); err != nil {
			return err
		}
		w.WriteByte(':')
		encoder, err := enc.registry.LookupEncoder(val)
		if err != nil {
			return err
		}
		if err := encoder(w, val); err != nil {
			return err
		}
	}
	w.WriteByte('}')
	return nil
}
