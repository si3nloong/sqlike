package types

import (
	"database/sql"
	"reflect"
	"strconv"

	"github.com/tidwall/gjson"
	"github.com/valyala/bytebufferpool"
)

// Any : type any is a type safe data type when you get the value
type Any struct {
	kind reflect.Kind
	raw  sql.RawBytes
}

func NewAny(it interface{}) Any {
	v := reflect.ValueOf(it)
	any := Any{kind: v.Kind()}
	switch any.kind {
	case reflect.String:
		any.raw = sql.RawBytes(v.String())
	case reflect.Bool:
		any.raw = sql.RawBytes(strconv.FormatBool(v.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		any.raw = sql.RawBytes(strconv.FormatInt(v.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		any.raw = sql.RawBytes(strconv.FormatUint(v.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		any.raw = sql.RawBytes(strconv.FormatFloat(v.Float(), 'f', 10, 64))
	default:
		panic("only support base type")
	}
	return any
}

// MarshalJSON :
func (a Any) MarshalJSON() ([]byte, error) {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)
	b.WriteByte('{')
	b.WriteString(`"kind":`)
	b.WriteString(strconv.FormatInt(int64(a.kind), 10))
	b.WriteByte(',')
	b.WriteString(`"raw":`)
	b.WriteString(string(a.raw))
	b.WriteByte('}')
	return b.Bytes(), nil
}

// MarshalJSON :
func (a *Any) UnmarshalJSON(b []byte) error {
	any := Any{}
	any.kind = reflect.Kind(gjson.GetBytes(b, "kind").Uint())
	any.raw = sql.RawBytes(gjson.GetBytes(b, "raw").Raw)
	*a = any
	return nil
}

// String :
func (a Any) String() string {
	return string(a.raw)
}

// Int64 :
func (a Any) Int64() int64 {
	num, _ := strconv.ParseInt(string(a.raw), 10, 64)
	return num
}

// Uint64 :
func (a Any) Uint64() uint64 {
	num, _ := strconv.ParseUint(string(a.raw), 10, 64)
	return num
}

// Float64 :
func (a Any) Float64() float64 {
	f, _ := strconv.ParseFloat(string(a.raw), 64)
	return f
}

// Bool :
func (a Any) Bool() bool {
	flag, _ := strconv.ParseBool(string(a.raw))
	return flag
}
