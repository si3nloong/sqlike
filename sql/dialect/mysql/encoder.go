package mysql

import (
	"bytes"
	"database/sql"
	"encoding"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/internal/spatial"
	"github.com/si3nloong/sqlike/v2/sql/codec"
	"github.com/si3nloong/sqlike/v2/x/reflext"

	"github.com/si3nloong/sqlike/v2/jsonb"

	"github.com/gofrs/uuid/v5"
)

// DefaultEncoders :
type DefaultEncoders struct {
	codec *codec.Registry
}

func (enc DefaultEncoders) EncodeUUID(d db.SqlDriver, v reflect.Value, opts map[string]string) (string, []any, error) {
	if v.IsZero() && v.CanSet() {
		id, err := uuid.NewV7()
		if err != nil {
			return "", nil, err
		}
		switch vi := v.Addr().Interface().(type) {
		case []byte:
			v.Set(reflect.ValueOf(id.Bytes()))
		case string:
			v.Set(reflect.ValueOf(id.String()))
		case sql.Scanner:
			if err := vi.Scan(id.String()); err != nil {
				return "", nil, err
			}
		case encoding.TextUnmarshaler:
			if err := vi.UnmarshalText(id.Bytes()); err != nil {
				return "", nil, err
			}
		}
		return d.Var(1), []any{id.Bytes()}, nil
	}
	fName := `UUID_TO_BIN(` + d.Var(1) + `)`
	switch vi := v.Interface().(type) {
	case []byte:
		// v.Set()
		return fName, []any{string(vi)}, nil
	case string:
		return fName, []any{vi}, nil
	case encoding.BinaryMarshaler:
		id, err := vi.MarshalBinary()
		if err != nil {
			return "", nil, err
		}
		return d.Var(1), []any{id}, nil

	// case fmt.Stringer:
	// 	return fName, []any{vi.String()}, nil
	// case driver.Valuer:
	// 	uuid, err := vi.Value()
	// 	if err != nil {
	// 		return "", nil, err
	// 	}
	// 	return fName, []any{uuid}, nil
	default:
		return "", nil, fmt.Errorf(`sqlike: invalid data type %v for UUID`, v.Type())
	}
}

func (enc DefaultEncoders) EncodeTime(d db.SqlDriver, v reflect.Value, opts map[string]string) (string, []any, error) {
	vi := v.Interface().(time.Time)
	if vi.IsZero() {
		fName := `CURRENT_TIMESTAMP`
		size, ok := opts["size"]
		if !ok {
			return fName, nil, nil
		}
		n, _ := strconv.Atoi(size)
		if n <= 0 {
			return fName, nil, nil
		}
		return fName + "(" + size + ")", nil, nil
	}
	return d.Var(1), []any{vi.UTC()}, nil
}

// EncodeSpatial :
func (enc DefaultEncoders) EncodeSpatial(t spatial.Type) db.ValueEncoder {
	return func(d db.SqlDriver, v reflect.Value, opts map[string]string) (string, []any, error) {
		x := v.Interface().(orb.Geometry)
		fName := "ST_PointFromText"
		switch t {
		case spatial.Point:
			fName = "ST_PointFromText"
		case spatial.LineString:
			fName = "ST_LineStringFromText"
		case spatial.Polygon:
			fName = "ST_PolygonFromText"
		case spatial.MultiPoint:
			fName = "ST_MultiPointFromText"
		case spatial.MultiLineString:
			fName = "ST_MultiLineStringFromText"
		case spatial.MultiPolygon:
			fName = "ST_MultiPolygonFromText"
		default:
		}
		fName = fName + "(" + d.Var(1)
		if v, ok := opts["srid"]; ok {
			if n, _ := strconv.Atoi(v); n > 0 {
				fName = fName + "," + v
			}
		}
		return fName + ")", []any{wkt.MarshalString(x)}, nil
	}
}

// EncodeByte :
func (enc DefaultEncoders) EncodeByte(d db.SqlDriver, v reflect.Value, opts map[string]string) (string, []any, error) {
	if v.IsNil() {
		return d.Var(1), []any{nil}, nil
	}
	b := reflext.Indirect(v).Bytes()
	if b == nil {
		return d.Var(1), []any{[]byte{}}, nil
	}
	x := base64.StdEncoding.EncodeToString(b)
	return d.Var(1), []any{[]byte(x)}, nil
}

// EncodeRawBytes :
func (enc DefaultEncoders) EncodeRawBytes(d db.SqlDriver, v reflect.Value, opts map[string]string) (string, []any, error) {
	return d.Var(1), []any{sql.RawBytes(v.Bytes())}, nil
}

// EncodeJSONRaw :
func (enc DefaultEncoders) EncodeJSONRaw(d db.SqlDriver, v reflect.Value, opts map[string]string) (string, []any, error) {
	if v.IsNil() {
		return d.Var(1), []any{[]byte("null")}, nil
	}
	buf := new(bytes.Buffer)
	if err := json.Compact(buf, v.Bytes()); err != nil {
		return "", nil, err
	}
	if buf.Len() == 0 {
		return d.Var(1), []any{[]byte(`{}`)}, nil
	}
	return d.Var(1), []any{json.RawMessage(buf.Bytes())}, nil
}

// EncodeStringer :
func (enc DefaultEncoders) EncodeStringer(d db.SqlDriver, v reflect.Value, opts map[string]string) (string, []any, error) {
	vi := v.Interface().(fmt.Stringer)
	return d.Var(1), []any{vi.String()}, nil
}

func (enc DefaultEncoders) EncodeDateTime(d db.SqlDriver, v reflect.Value, opts map[string]string) (string, []any, error) {
	vi := v.Interface().(time.Time)
	if vi.IsZero() {
		vi = time.Now()
		if v.CanSet() {
			v.Set(reflect.ValueOf(vi))
		}
	}
	// convert to UTC before storing into DB
	return d.Var(1), []any{vi.UTC()}, nil
}

// // EncodeTime :
// func (enc DefaultEncoders) EncodeDateTime(ctx context.Context, v reflect.Value) (any, error) {
// 	x := v.Interface().(time.Time)
// 	// TODO:
// 	// If `CURRENT_TIMESTAMP` is define, we will pass `nil` value
// 	// elseif the datetime is zero, we will pass the current datetime of the machine
// 	// else we will pass the value of it instead
// 	//
// 	// And we should handle the TIMESTAMP type as well
// 	f := sqlx.GetField(ctx)
// 	// def, _ := f.Tag().LookUp("default")
// 	v2, _ := f.Tag().Option("default")
// 	if strings.Contains(v2, "CURRENT_TIMESTAMP") {
// 		return nil, nil
// 	} else if x.IsZero() {
// 		x, _ = time.Parse(time.RFC3339, "1970-01-01T08:00:00Z")
// 		return x, nil
// 	}
// 	// convert to UTC before storing into DB
// 	return x.UTC(), nil
// }

// EncodeString :
func (enc DefaultEncoders) EncodeString(d db.SqlDriver, v reflect.Value, opts map[string]string) (string, []any, error) {
	str := v.String()
	if str == "" {
		if val, ok := opts["enum"]; ok {
			enums := strings.Split(val, "|")
			if len(enums) > 0 {
				return d.Var(1), []any{enums[0]}, nil
			}
		}
	}
	return d.Var(1), []any{v.String()}, nil
}

// EncodeBool :
func (enc DefaultEncoders) EncodeBool(d db.SqlDriver, v reflect.Value, opts map[string]string) (string, []any, error) {
	return d.Var(1), []any{v.Bool()}, nil
}

// EncodeInt :
func (enc DefaultEncoders) EncodeInt(d db.SqlDriver, v reflect.Value, opts map[string]string) (string, []any, error) {
	if _, ok := opts["auto_increment"]; ok {
		return d.Var(1), []any{nil}, nil
	}
	return d.Var(1), []any{v.Int()}, nil
}

// EncodeUint :
func (enc DefaultEncoders) EncodeUint(d db.SqlDriver, v reflect.Value, opts map[string]string) (string, []any, error) {
	if _, ok := opts["auto_increment"]; ok {
		return d.Var(1), []any{nil}, nil
	}
	return d.Var(1), []any{v.Uint()}, nil
}

// EncodeFloat :
func (enc DefaultEncoders) EncodeFloat(d db.SqlDriver, v reflect.Value, opts map[string]string) (string, []any, error) {
	return d.Var(1), []any{v.Float()}, nil
}

// EncodePtr :
func (enc *DefaultEncoders) EncodePtr(d db.SqlDriver, v reflect.Value, opts map[string]string) (string, []any, error) {
	if !v.IsValid() || v.IsNil() {
		return d.Var(1), []any{nil}, nil
	}
	v = v.Elem()
	encoder, err := enc.codec.LookupEncoder(v)
	if err != nil {
		return "", nil, err
	}
	return encoder(d, v, opts)
}

// EncodeStruct :
func (enc DefaultEncoders) EncodeStruct(d db.SqlDriver, v reflect.Value, opts map[string]string) (string, []any, error) {
	b, err := jsonb.Marshal(v)
	if err != nil {
		return "", nil, err
	}
	return d.Var(1), []any{b}, nil
}

// EncodeArray :
func (enc DefaultEncoders) EncodeArray(d db.SqlDriver, v reflect.Value, opts map[string]string) (string, []any, error) {
	b, err := jsonb.Marshal(v)
	if err != nil {
		return "", nil, err
	}
	return d.Var(1), []any{b}, nil
}

// EncodeMap :
func (enc DefaultEncoders) EncodeMap(d db.SqlDriver, v reflect.Value, opts map[string]string) (string, []any, error) {
	if v.IsNil() {
		return d.Var(1), []any{"null"}, nil
	}
	b, err := jsonb.Marshal(v)
	if err != nil {
		return "", nil, err
	}
	return d.Var(1), []any{b}, nil
}
