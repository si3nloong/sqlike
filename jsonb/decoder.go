package jsonb

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/si3nloong/sqlike/reflext"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

// date format :
var (
	ddmmyyyy       = regexp.MustCompile(`^\d{4}\-\d{2}\-\d{2}$`)
	ddmmyyyyhhmmss = regexp.MustCompile(`^\d{4}\-\d{2}\-\d{2}\s\d{2}\:\d{2}:\d{2}$`)
)

// DefaultDecoder :
type DefaultDecoder struct {
	registry *Registry
}

// DecodeByte :
func (dec DefaultDecoder) DecodeByte(r *Reader, v reflect.Value) error {
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

// DecodeLanguage :
func (dec DefaultDecoder) DecodeLanguage(r *Reader, v reflect.Value) error {
	str, err := r.ReadString()
	if err != nil {
		return err
	}
	if str == "" {
		return nil
	}
	l, err := language.Parse(str)
	if err != nil {
		v.Set(reflect.ValueOf(language.Tag{}))
		return err
	}
	v.Set(reflect.ValueOf(l))
	return nil
}

// DecodeCurrency :
func (dec DefaultDecoder) DecodeCurrency(r *Reader, v reflect.Value) error {
	str, err := r.ReadString()
	if err != nil {
		return err
	}
	if str == "" {
		v.Set(reflect.ValueOf(currency.Unit{}))
		return nil
	}
	cur, err := currency.ParseISO(str)
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(cur))
	return nil
}

// DecodeTime :
func (dec DefaultDecoder) DecodeTime(r *Reader, v reflect.Value) error {
	b, err := r.ReadBytes()
	if err != nil {
		return err
	}
	str := string(b)
	if str == null || str == `""` {
		reflext.Set(v, reflect.ValueOf(time.Time{}))
		return nil
	}
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return errors.New("jsonb: invalid format of date")
	}
	str = string(b[1 : len(b)-1])
	var x time.Time
	x, err = DecodeTime(str)
	if err != nil {
		return err
	}
	reflext.Set(v, reflect.ValueOf(x))
	return nil
}

// DecodeJSONRaw :
func (dec DefaultDecoder) DecodeJSONRaw(r *Reader, v reflect.Value) error {
	v.SetBytes(r.Bytes())
	return nil
}

// DecodeJSONNumber :
func (dec DefaultDecoder) DecodeJSONNumber(r *Reader, v reflect.Value) error {
	x, err := r.ReadNumber()
	if err != nil {
		return err
	}
	v.SetString(x.String())
	return nil
}

// DecodeString :
func (dec DefaultDecoder) DecodeString(r *Reader, v reflect.Value) error {
	x, err := r.ReadEscapeString()
	if err != nil {
		return err
	}
	v.SetString(x)
	return nil
}

// DecodeBool :
func (dec DefaultDecoder) DecodeBool(r *Reader, v reflect.Value) error {
	x, err := r.ReadBoolean()
	if err != nil {
		return err
	}
	v.SetBool(x)
	return nil
}

// DecodeInt :
func (dec DefaultDecoder) DecodeInt(quote bool) ValueDecoder {
	return func(r *Reader, v reflect.Value) error {
		var (
			str string
			err error
		)
		if quote {
			str, err = r.ReadString()
			if err != nil {
				return err
			}
		} else {
			num, err := r.ReadNumber()
			if err != nil {
				return err
			}
			str = num.String()
		}
		x, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		if v.OverflowInt(x) {
			return errors.New("integer overflow")
		}
		v.SetInt(x)
		return nil
	}
}

// DecodeUint :
func (dec DefaultDecoder) DecodeUint(quote bool) ValueDecoder {
	return func(r *Reader, v reflect.Value) error {
		var (
			str string
			err error
		)
		if quote {
			str, err = r.ReadString()
			if err != nil {
				return err
			}
		} else {
			num, err := r.ReadNumber()
			if err != nil {
				return err
			}
			str = num.String()
		}
		x, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		if v.OverflowUint(x) {
			return errors.New("unsigned integer overflow")
		}
		v.SetUint(x)
		return nil
	}
}

// DecodeFloat :
func (dec DefaultDecoder) DecodeFloat(r *Reader, v reflect.Value) error {
	num, err := r.ReadNumber()
	if err != nil {
		return err
	}
	x, err := strconv.ParseFloat(num.String(), 64)
	if err != nil {
		return err
	}
	if v.OverflowFloat(x) {
		return errors.New("jsonb: float overflow")
	}
	v.SetFloat(x)
	return nil
}

// DecodePtr :
func (dec *DefaultDecoder) DecodePtr(r *Reader, v reflect.Value) error {
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
func (dec *DefaultDecoder) DecodeStruct(r *Reader, v reflect.Value) error {
	mapper := reflext.DefaultMapper
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
func (dec *DefaultDecoder) DecodeArray(r *Reader, v reflect.Value) error {
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
func (dec *DefaultDecoder) DecodeSlice(r *Reader, v reflect.Value) error {
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
func (dec *DefaultDecoder) DecodeMap(r *Reader, v reflect.Value) error {
	if r.IsNull() {
		v.Set(reflect.Zero(v.Type()))
		return r.skipNull()
	}

	var (
		decodeKey ValueDecoder
		err       error
		t         = v.Type()
		kind      = t.Key().Kind()
	)

	switch kind {
	case reflect.String:
		decodeKey = dec.registry.kindDecoders[reflect.String]
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		decodeKey = dec.DecodeInt(true)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		decodeKey = dec.DecodeUint(true)
	default:
		key := t.Key()
		if key.Kind() != reflect.Ptr {
			key = reflect.PtrTo(key)
		}
		if !key.Implements(textUnmarshaler) {
			return fmt.Errorf("jsonb: unsupported data type of map key, %q", t.Key().Kind())
		}
		decodeKey = textUnmarshalerDecoder()
	}

	decodeValue, err := dec.registry.LookupDecoder(t.Elem())
	if err != nil {
		return err
	}
	x := reflect.MakeMap(t)
	if err := r.ReadObject(func(it *Reader, k string) error {
		ki := reflext.Zero(t.Key())
		err = decodeKey(NewReader([]byte(strconv.Quote(k))), ki)
		if err != nil {
			return err
		}
		vi := reflext.Zero(t.Elem())
		err = decodeValue(it, vi)
		if err != nil {
			return err
		}
		x.SetMapIndex(ki, vi)
		return nil
	}); err != nil {
		return err
	}
	v.Set(x)
	return nil
}

// DecodeInterface :
func (dec DefaultDecoder) DecodeInterface(r *Reader, v reflect.Value) error {
	x, err := r.ReadValue()
	if err != nil {
		return err
	}
	if x != nil {
		v.Set(reflect.ValueOf(x))
	}
	return nil
}

// DecodeTime :
func DecodeTime(str string) (t time.Time, err error) {
	switch {
	case ddmmyyyy.MatchString(str):
		t, err = time.Parse("2006-01-02", str)
	case ddmmyyyyhhmmss.MatchString(str):
		t, err = time.Parse("2006-01-02 15:04:05", str)
	default:
		t, err = time.Parse(time.RFC3339Nano, str)
	}
	return
}
