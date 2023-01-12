package mysql

import (
	"bytes"
	"context"
	"database/sql"
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

	sqlx "github.com/si3nloong/sqlike/v2/sql"

	"github.com/si3nloong/sqlike/v2/jsonb"
)

// DefaultEncoders :
type DefaultEncoders struct {
	codec *codec.Registry
}

// EncodeByte :
func (enc DefaultEncoders) EncodeByte(_ context.Context, v reflect.Value) (any, error) {
	b := v.Bytes()
	if b == nil {
		return make([]byte, 0), nil
	}
	x := base64.StdEncoding.EncodeToString(b)
	return []byte(x), nil
}

// EncodeRawBytes :
func (enc DefaultEncoders) EncodeRawBytes(_ context.Context, v reflect.Value) (any, error) {
	return sql.RawBytes(v.Bytes()), nil
}

// EncodeJSONRaw :
func (enc DefaultEncoders) EncodeJSONRaw(_ context.Context, v reflect.Value) (any, error) {
	if v.IsNil() {
		return []byte("null"), nil
	}
	buf := new(bytes.Buffer)
	if err := json.Compact(buf, v.Bytes()); err != nil {
		return nil, err
	}
	if buf.Len() == 0 {
		return []byte(`{}`), nil
	}
	return json.RawMessage(buf.Bytes()), nil
}

// EncodeStringer :
func (enc DefaultEncoders) EncodeStringer(_ context.Context, v reflect.Value) (any, error) {
	x := v.Interface().(fmt.Stringer)
	return x.String(), nil
}

// EncodeTime :
func (enc DefaultEncoders) EncodeDateTime(ctx context.Context, v reflect.Value) (any, error) {
	x := v.Interface().(time.Time)
	// TODO:
	// If `CURRENT_TIMESTAMP` is define, we will pass `nil` value
	// elseif the datetime is zero, we will pass the current datetime of the machine
	// else we will pass the value of it instead
	//
	// And we should handle the TIMESTAMP type as well
	f := sqlx.GetField(ctx)
	def, ok := f.Tag().LookUp("default")
	if ok && strings.EqualFold(def, "CURRENT_TIMESTAMP") {
		return nil, nil
	} else if x.IsZero() {
		x, _ = time.Parse(time.RFC3339, "1970-01-01T08:00:00Z")
		return x, nil
	}
	// convert to UTC before storing into DB
	return x.UTC(), nil
}

// EncodeSpatial :
func (enc DefaultEncoders) EncodeSpatial(st spatial.Type) db.ValueEncoder {
	return func(ctx context.Context, v reflect.Value) (any, error) {
		if reflext.IsZero(v) {
			return nil, nil
		}

		f := sqlx.GetField(ctx)
		x := v.Interface().(orb.Geometry)
		// var srid uint
		tag, ok := f.Tag().Option("srid")
		if ok {
			integer, _ := strconv.Atoi(tag)
			if integer > 0 {
				// srid = uint(integer)
			}
		}

		// switch vi.Type {
		// case spatial.Point:
		// 	stmt.WriteString("ST_PointFromText")
		// case spatial.LineString:
		// 	stmt.WriteString("ST_LineStringFromText")
		// case spatial.Polygon:
		// 	stmt.WriteString("ST_PolygonFromText")
		// case spatial.MultiPoint:
		// 	stmt.WriteString("ST_MultiPointFromText")
		// case spatial.MultiLineString:
		// 	stmt.WriteString("ST_MultiLineStringFromText")
		// case spatial.MultiPolygon:
		// 	stmt.WriteString("ST_MultiPolygonFromText")
		// default:
		// }

		// stmt.WriteString("(?")
		// if vi.SRID > 0 {
		// 	stmt.WriteString(fmt.Sprintf(",%d", vi.SRID))
		// }
		// stmt.WriteByte(')')
		// stmt.AppendArgs(vi.WKT)
		// log.Println(x)
		return sql.RawBytes(`ST_PointFromText(` + wkt.MarshalString(x) + `)`), nil
	}
}

// EncodeString :
func (enc DefaultEncoders) EncodeString(ctx context.Context, v reflect.Value) (any, error) {
	str := v.String()
	f := sqlx.GetField(ctx)
	if str == "" {
		tag := f.Tag()
		if val, ok := tag.Option("enum"); ok {
			enums := strings.Split(val, "|")
			if len(enums) > 0 {
				return enums[0], nil
			}
		}
	}
	return str, nil
}

// EncodeBool :
func (enc DefaultEncoders) EncodeBool(_ context.Context, v reflect.Value) (any, error) {
	return v.Bool(), nil
}

// EncodeInt :
func (enc DefaultEncoders) EncodeInt(ctx context.Context, v reflect.Value) (any, error) {
	f := sqlx.GetField(ctx)
	if _, ok := f.Tag().Option("auto_increment"); ok {
		return nil, nil
	}
	return v.Int(), nil
}

// EncodeUint :
func (enc DefaultEncoders) EncodeUint(ctx context.Context, v reflect.Value) (any, error) {
	f := sqlx.GetField(ctx)
	if _, ok := f.Tag().Option("auto_increment"); ok {
		return nil, nil
	}
	return v.Uint(), nil
}

// EncodeFloat :
func (enc DefaultEncoders) EncodeFloat(_ context.Context, v reflect.Value) (any, error) {
	return v.Float(), nil
}

// EncodePtr :
func (enc *DefaultEncoders) EncodePtr(ctx context.Context, v reflect.Value) (any, error) {
	if !v.IsValid() || v.IsNil() {
		return nil, nil
	}
	v = v.Elem()
	encoder, err := enc.codec.LookupEncoder(v)
	if err != nil {
		return nil, err
	}
	return encoder(ctx, v)
}

// EncodeStruct :
func (enc DefaultEncoders) EncodeStruct(_ context.Context, v reflect.Value) (any, error) {
	return jsonb.Marshal(v)
}

// EncodeArray :
func (enc DefaultEncoders) EncodeArray(_ context.Context, v reflect.Value) (any, error) {
	return jsonb.Marshal(v)
}

// EncodeMap :
func (enc DefaultEncoders) EncodeMap(_ context.Context, v reflect.Value) (any, error) {
	if v.IsNil() {
		return string("null"), nil
	}
	return jsonb.Marshal(v)
}

func isBaseType(t reflect.Type) bool {
	for {
		k := t.Kind()
		switch k {
		case reflect.String:
			return true
		case reflect.Bool:
			return true
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return true
		case reflect.Float32, reflect.Float64:
			return true
		case reflect.Ptr:
			t = t.Elem()
		default:
			return false
		}
	}
}
