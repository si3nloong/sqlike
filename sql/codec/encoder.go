package codec

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/x/reflext"
	"github.com/si3nloong/sqlike/v2/x/spatial"

	"github.com/si3nloong/sqlike/jsonb"
)

// DefaultEncoders :
type DefaultEncoders struct {
	codec *Registry
}

// EncodeByte :
func (enc DefaultEncoders) EncodeByte(_ context.Context, v reflect.Value) (interface{}, error) {
	b := v.Bytes()
	if b == nil {
		return make([]byte, 0), nil
	}
	x := base64.StdEncoding.EncodeToString(b)
	return []byte(x), nil
}

// EncodeRawBytes :
func (enc DefaultEncoders) EncodeRawBytes(_ context.Context, v reflect.Value) (interface{}, error) {
	return sql.RawBytes(v.Bytes()), nil
}

// EncodeJSONRaw :
func (enc DefaultEncoders) EncodeJSONRaw(_ context.Context, v reflect.Value) (interface{}, error) {
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
func (enc DefaultEncoders) EncodeStringer(_ context.Context, v reflect.Value) (interface{}, error) {
	x := v.Interface().(fmt.Stringer)
	return x.String(), nil
}

// EncodePointerStringer :
func (enc DefaultEncoders) EncodePointerStringer(_ context.Context, v reflect.Value) (interface{}, error) {
	x := v.Addr().Interface().(fmt.Stringer)
	return x.String(), nil
}

// EncodeTime :
func (enc DefaultEncoders) EncodeTime(_ context.Context, v reflect.Value) (interface{}, error) {
	x := v.Interface().(time.Time)
	// if x.IsZero() {
	// 	x, _ = time.Parse(time.RFC3339, "1970-01-01T08:00:00Z")
	// 	return x, nil
	// }
	// convert to UTC before storing into DB
	return x.UTC(), nil
}

// EncodeSpatial :
func (enc DefaultEncoders) EncodeSpatial(st spatial.Type) db.ValueEncoder {
	return func(ctx context.Context, v reflect.Value) (interface{}, error) {
		if reflext.IsZero(v) {
			return nil, nil
		}
		x := v.Interface().(orb.Geometry)
		var srid uint
		// if sf != nil {
		// 	tag, ok := sf.Tag().LookUp("srid")
		// 	if ok {
		// 		integer, _ := strconv.Atoi(tag)
		// 		if integer > 0 {
		// 			srid = uint(integer)
		// 		}
		// 	}
		// }
		return spatial.Geometry{
			Type: st,
			SRID: srid,
			WKT:  wkt.MarshalString(x),
		}, nil
	}
}

// EncodeString :
func (enc DefaultEncoders) EncodeString(sf context.Context, v reflect.Value) (interface{}, error) {
	str := v.String()
	if str == "" && sf != nil {
		// FIXME:
		// tag := sf.Tag()
		// if val, ok := tag.LookUp("enum"); ok {
		// 	enums := strings.Split(val, "|")
		// 	if len(enums) > 0 {
		// 		return enums[0], nil
		// 	}
		// }
	}
	return str, nil
}

// EncodeBool :
func (enc DefaultEncoders) EncodeBool(_ context.Context, v reflect.Value) (interface{}, error) {
	return v.Bool(), nil
}

// EncodeInt :
func (enc DefaultEncoders) EncodeInt(_ context.Context, v reflect.Value) (interface{}, error) {
	return v.Int(), nil
}

// EncodeUint :
func (enc DefaultEncoders) EncodeUint(_ context.Context, v reflect.Value) (interface{}, error) {
	return v.Uint(), nil
}

// EncodeFloat :
func (enc DefaultEncoders) EncodeFloat(_ context.Context, v reflect.Value) (interface{}, error) {
	return v.Float(), nil
}

// EncodePtr :
func (enc *DefaultEncoders) EncodePtr(sf context.Context, v reflect.Value) (interface{}, error) {
	if !v.IsValid() || v.IsNil() {
		return nil, nil
	}
	v = v.Elem()
	encoder, err := enc.codec.LookupEncoder(v)
	if err != nil {
		return nil, err
	}
	return encoder(sf, v)
}

// EncodeStruct :
func (enc DefaultEncoders) EncodeStruct(_ context.Context, v reflect.Value) (interface{}, error) {
	return jsonb.Marshal(v)
}

// EncodeArray :
func (enc DefaultEncoders) EncodeArray(_ context.Context, v reflect.Value) (interface{}, error) {
	return jsonb.Marshal(v)
}

// EncodeMap :
func (enc DefaultEncoders) EncodeMap(_ context.Context, v reflect.Value) (interface{}, error) {
	if v.IsNil() {
		return string("null"), nil
	}

	// t := v.Type()
	// k := t.Key()
	// if k.Kind() != reflect.String {
	// 	return nil, fmt.Errorf("codec: unsupported data type %q for map key, it must be string", k.Kind())
	// }
	// k = t.Elem()
	// if !isBaseType(k) {
	// 	return nil, fmt.Errorf("codec: unsupported data type %q for map value", k.Kind())
	// }
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
