package codec

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/si3nloong/sqlike/jsonb"
	"github.com/si3nloong/sqlike/types"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"

	"errors"
)

// DefaultDecoders :
type DefaultDecoders struct {
	registry *Registry
}

// DecodeByte :
func (dec DefaultDecoders) DecodeByte(it interface{}, v reflect.Value) error {
	var (
		x   []byte
		err error
	)
	switch vi := it.(type) {
	case string:
		x, err = base64.StdEncoding.DecodeString(vi)
		if err != nil {
			return err
		}
	case []byte:
		x, err = base64.StdEncoding.DecodeString(string(vi))
		if err != nil {
			return err
		}
	case nil:
		x = make([]byte, 0)
	}
	v.SetBytes(x)
	return nil
}

// DecodeRawBytes :
func (dec DefaultDecoders) DecodeRawBytes(it interface{}, v reflect.Value) error {
	var (
		x sql.RawBytes
	)
	switch vi := it.(type) {
	case []byte:
		x = sql.RawBytes(vi)
	case string:
		x = sql.RawBytes(vi)
	case sql.RawBytes:
		x = vi
	case bool:
		str := strconv.FormatBool(vi)
		x = []byte(str)
	case int64:
		str := strconv.FormatInt(vi, 10)
		x = []byte(str)
	case uint64:
		str := strconv.FormatUint(vi, 10)
		x = []byte(str)
	case float64:
		str := strconv.FormatFloat(vi, 'e', -1, 64)
		x = []byte(str)
	case time.Time:
		x = []byte(vi.Format(time.RFC3339))
	case nil:
	default:
	}
	v.SetBytes(x)
	return nil
}

// DecodeCurrency :
func (dec DefaultDecoders) DecodeCurrency(it interface{}, v reflect.Value) error {
	var (
		x   currency.Unit
		err error
	)
	switch vi := it.(type) {
	case string:
		x, err = currency.ParseISO(vi)
		if err != nil {
			return err
		}
	case []byte:
		x, err = currency.ParseISO(string(vi))
		if err != nil {
			return err
		}
	case nil:
	}
	v.Set(reflect.ValueOf(x))
	return nil
}

// DecodeLanguage :
func (dec DefaultDecoders) DecodeLanguage(it interface{}, v reflect.Value) error {
	var (
		x   language.Tag
		str string
		err error
	)
	switch vi := it.(type) {
	case string:
		str = vi
	case []byte:
		str = string(vi)
	case nil:
	default:
		return errors.New("language tag is not well-formed")
	}
	if str != "" {
		x, err = language.Parse(str)
		if err != nil {
			return err
		}
	}
	v.Set(reflect.ValueOf(x))
	return nil
}

// DecodeJSONRaw :
func (dec DefaultDecoders) DecodeJSONRaw(it interface{}, v reflect.Value) error {
	b := new(bytes.Buffer)
	switch vi := it.(type) {
	case string:
		if err := json.Compact(b, []byte(vi)); err != nil {
			return err
		}
	case []byte:
		if err := json.Compact(b, vi); err != nil {
			return err
		}
	case nil:
	}
	v.SetBytes(b.Bytes())
	return nil
}

// DecodeTime :
func (dec DefaultDecoders) DecodeTime(it interface{}, v reflect.Value) error {
	var (
		x   time.Time
		err error
	)
	switch vi := it.(type) {
	case time.Time:
		x = vi
	case string:
		x, err = types.DecodeTime(vi)
		if err != nil {
			return err
		}
	case []byte:
		x, err = types.DecodeTime(b2s(vi))
		if err != nil {
			return err
		}
	case nil:
	}
	// convert back to UTC
	v.Set(reflect.ValueOf(x.UTC()))
	return nil
}

// DecodeSpatial :
func (dec DefaultDecoders) DecodePoint(it interface{}, v reflect.Value) error {
	data, ok := it.([]byte)
	if !ok {
		return errors.New("point must be []byte")
	}

	if len(data) == 42 {
		dst := make([]byte, 21)
		_, err := hex.Decode(dst, data)
		if err != nil {
			return err
		}
		data = dst
	}

	p := orb.Point{}
	scanner := wkb.Scanner(&p)
	// if len(data) == 21 {
	// 	// the length of a point type in WKB
	// 	return scan.Scan(data[:])
	// }

	if len(data) == 25 {
		// Most likely MySQL's SRID+WKB format.
		// However, could be a line string or multipoint with only one point.
		// But those would be invalid for parsing a point.
		// return p.unmarshalWKB(data[4:])
		if err := scanner.Scan(data[4:]); err != nil {
			return err
		}
		v.Set(reflect.ValueOf(p))
		return nil
	}

	if len(data) == 0 {
		// empty data, return empty go struct which in this case
		// would be [0,0]
		return nil
	}

	return errors.New("incorrect point")
}

// DecodeLineString :
func (dec DefaultDecoders) DecodeLineString(it interface{}, v reflect.Value) error {
	data, ok := it.([]byte)
	if !ok {
		return errors.New("line string must be []byte")
	}

	if len(data) == 0 {
		return nil
	}

	p := orb.LineString{}
	scanner := wkb.Scanner(&p)
	if err := scanner.Scan(data[4:]); err != nil {
		return err
	}
	v.Set(reflect.ValueOf(p))
	return nil
}

// DecodeString :
func (dec DefaultDecoders) DecodeString(it interface{}, v reflect.Value) error {
	var x string
	switch vi := it.(type) {
	case string:
		x = vi
	case []byte:
		x = string(vi)
	case int64:
		x = strconv.FormatInt(vi, 10)
	case uint64:
		x = strconv.FormatUint(vi, 10)
	case nil:
	}
	v.SetString(x)
	return nil
}

// DecodeBool :
func (dec DefaultDecoders) DecodeBool(it interface{}, v reflect.Value) error {
	var (
		x   bool
		err error
	)
	switch vi := it.(type) {
	case []byte:
		x, err = strconv.ParseBool(b2s(vi))
		if err != nil {
			return err
		}
	case string:
		x, err = strconv.ParseBool(vi)
		if err != nil {
			return err
		}
	case bool:
		x = vi
	case int64:
		if vi == 1 {
			x = true
		}
	case uint64:
		if vi == 1 {
			x = true
		}
	case nil:
	}
	v.SetBool(x)
	return nil
}

// DecodeInt :
func (dec DefaultDecoders) DecodeInt(it interface{}, v reflect.Value) error {
	var (
		x   int64
		err error
	)
	switch vi := it.(type) {
	case []byte:
		x, err = strconv.ParseInt(b2s(vi), 10, 64)
		if err != nil {
			return err
		}
	case string:
		x, err = strconv.ParseInt(vi, 10, 64)
		if err != nil {
			return err
		}
	case int64:
		x = vi
	case uint64:
		x = int64(vi)
	case nil:
	}
	if v.OverflowInt(x) {
		return errors.New("integer overflow")
	}
	v.SetInt(x)
	return nil
}

// DecodeUint :
func (dec DefaultDecoders) DecodeUint(it interface{}, v reflect.Value) error {
	var (
		x   uint64
		err error
	)
	switch vi := it.(type) {
	case []byte:
		x, err = strconv.ParseUint(b2s(vi), 10, 64)
		if err != nil {
			return err
		}
	case string:
		x, err = strconv.ParseUint(vi, 10, 64)
		if err != nil {
			return err
		}
	case int64:
		x = uint64(vi)
	case uint64:
		x = vi
	case nil:
	}
	if v.OverflowUint(x) {
		return errors.New("unsigned integer overflow")
	}
	v.SetUint(x)
	return nil
}

// DecodeFloat :
func (dec DefaultDecoders) DecodeFloat(it interface{}, v reflect.Value) error {
	var (
		x   float64
		err error
	)
	switch vi := it.(type) {
	case []byte:
		x, err = strconv.ParseFloat(b2s(vi), 64)
		if err != nil {
			return err
		}
	case string:
		x, err = strconv.ParseFloat(vi, 64)
		if err != nil {
			return err
		}
	case float64:
		x = vi
	case int64:
		x = float64(vi)
	case uint64:
		x = float64(vi)
	case nil:

	}
	if v.OverflowFloat(x) {
		return errors.New("float overflow")
	}
	v.SetFloat(x)
	return nil
}

// DecodePtr :
func (dec *DefaultDecoders) DecodePtr(it interface{}, v reflect.Value) error {
	t := v.Type()
	if it == nil {
		v.Set(reflect.Zero(t))
		return nil
	}
	t = t.Elem()
	decoder, err := dec.registry.LookupDecoder(t)
	if err != nil {
		return err
	}
	return decoder(it, v.Elem())
}

// DecodeStruct :
func (dec *DefaultDecoders) DecodeStruct(it interface{}, v reflect.Value) error {
	var b []byte
	switch vi := it.(type) {
	case string:
		b = []byte(vi)
	case []byte:
		b = vi
	}
	return jsonb.UnmarshalValue(b, v)
}

// DecodeArray :
func (dec DefaultDecoders) DecodeArray(it interface{}, v reflect.Value) error {
	var b []byte
	switch vi := it.(type) {
	case string:
		b = []byte(vi)
	case []byte:
		b = vi
	}
	return jsonb.UnmarshalValue(b, v)
}

// DecodeMap :
func (dec DefaultDecoders) DecodeMap(it interface{}, v reflect.Value) error {
	var b []byte
	switch vi := it.(type) {
	case string:
		b = []byte(vi)
	case []byte:
		b = vi
	}
	return jsonb.UnmarshalValue(b, v)
}
