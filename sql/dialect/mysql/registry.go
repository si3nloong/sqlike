package mysql

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"time"

	"cloud.google.com/go/civil"
	"github.com/paulmach/orb"
	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/internal/spatial"
	"github.com/si3nloong/sqlike/v2/sql/codec"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

func buildDefaultRegistry() db.Codecer {
	rg := codec.NewRegistry()
	dec := DefaultDecoders{rg}
	enc := DefaultEncoders{rg}
	rg.RegisterTypeCodec(reflect.TypeOf([]byte{}), enc.EncodeByte, dec.DecodeByte)
	rg.RegisterTypeCodec(reflect.TypeOf(language.Tag{}), enc.EncodeStringer, dec.DecodeLanguage)
	rg.RegisterTypeCodec(reflect.TypeOf(currency.Unit{}), enc.EncodeStringer, dec.DecodeCurrency)
	rg.RegisterTypeCodec(reflect.TypeOf(civil.Time{}), enc.EncodeStringer, dec.DecodeTime)
	rg.RegisterTypeCodec(reflect.TypeOf(civil.Date{}), enc.EncodeStringer, dec.DecodeDate)
	rg.RegisterTypeCodec(reflect.TypeOf(time.Time{}), enc.EncodeDateTime, dec.DecodeDateTime)
	// rg.RegisterTypeCodec(reflect.TypeOf(time.Location{}), enc.EncodeTime, dec.DecodeTime)
	rg.RegisterTypeCodec(reflect.TypeOf(sql.RawBytes{}), enc.EncodeRawBytes, dec.DecodeRawBytes)
	rg.RegisterTypeCodec(reflect.TypeOf(json.RawMessage{}), enc.EncodeJSONRaw, dec.DecodeJSONRaw)
	rg.RegisterTypeCodec(reflect.TypeOf(orb.Point{}), enc.EncodeSpatial(spatial.Point), dec.DecodePoint)
	rg.RegisterTypeCodec(reflect.TypeOf(orb.LineString{}), enc.EncodeSpatial(spatial.LineString), dec.DecodeLineString)
	// rg.RegisterTypeCodec(reflect.TypeOf(orb.Polygon{}), enc.EncodeSpatial(spatial.Polygon), dec.DecodePolygon)
	// rg.RegisterTypeCodec(reflect.TypeOf(orb.MultiPoint{}), enc.EncodeSpatial(spatial.MultiPoint), dec.DecodeMultiPoint)
	// rg.RegisterTypeCodec(reflect.TypeOf(orb.MultiLineString{}), enc.EncodeSpatial(spatial.MultiLineString), dec.DecodeMultiLineString)
	// rg.RegisterTypeCodec(reflect.TypeOf(orb.MultiPolygon{}), enc.EncodeSpatial(spatial.MultiPolygon), dec.DecodeMultiPolygon)
	rg.RegisterKindCodec(reflect.String, enc.EncodeString, dec.DecodeString)
	rg.RegisterKindCodec(reflect.Bool, enc.EncodeBool, dec.DecodeBool)
	rg.RegisterKindCodec(reflect.Int, enc.EncodeInt, dec.DecodeInt)
	rg.RegisterKindCodec(reflect.Int8, enc.EncodeInt, dec.DecodeInt)
	rg.RegisterKindCodec(reflect.Int16, enc.EncodeInt, dec.DecodeInt)
	rg.RegisterKindCodec(reflect.Int32, enc.EncodeInt, dec.DecodeInt)
	rg.RegisterKindCodec(reflect.Int64, enc.EncodeInt, dec.DecodeInt)
	rg.RegisterKindCodec(reflect.Uint, enc.EncodeUint, dec.DecodeUint)
	rg.RegisterKindCodec(reflect.Uint8, enc.EncodeUint, dec.DecodeUint)
	rg.RegisterKindCodec(reflect.Uint16, enc.EncodeUint, dec.DecodeUint)
	rg.RegisterKindCodec(reflect.Uint32, enc.EncodeUint, dec.DecodeUint)
	rg.RegisterKindCodec(reflect.Uint64, enc.EncodeUint, dec.DecodeUint)
	rg.RegisterKindCodec(reflect.Float32, enc.EncodeFloat, dec.DecodeFloat)
	rg.RegisterKindCodec(reflect.Float64, enc.EncodeFloat, dec.DecodeFloat)
	rg.RegisterKindCodec(reflect.Ptr, enc.EncodePtr, dec.DecodePtr)
	rg.RegisterKindCodec(reflect.Struct, enc.EncodeStruct, dec.DecodeStruct)
	rg.RegisterKindCodec(reflect.Array, enc.EncodeArray, dec.DecodeArray)
	rg.RegisterKindCodec(reflect.Slice, enc.EncodeArray, dec.DecodeArray)
	rg.RegisterKindCodec(reflect.Map, enc.EncodeMap, dec.DecodeMap)
	return rg
}
