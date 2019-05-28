package core

// import (
// 	"bytes"
// 	"encoding/base64"
// 	"encoding/json"
// 	"reflect"
// 	"strconv"
// 	"time"

// 	"github.com/si3nloong/sqlike/core/codec"
// 	"github.com/si3nloong/sqlike/jsonb"
// )

// var (
// 	uintByte  [10]byte
// 	floatByte [10]byte
// )

// // DefaultEncoders :
// type DefaultEncoders struct {
// 	registry *codec.Registry
// }

// // SetEncoders :
// func (enc DefaultEncoders) SetEncoders(rg *codec.Registry) {
// 	rg.SetTypeEncoder(reflect.TypeOf([]byte{}), enc.EncodeByte)
// 	rg.SetTypeEncoder(reflect.TypeOf(time.Time{}), enc.EncodeTime)
// 	rg.SetTypeEncoder(reflect.TypeOf(json.RawMessage{}), enc.EncodeJSONRaw)
// 	rg.SetKindEncoder(reflect.String, enc.EncodeString)
// 	rg.SetKindEncoder(reflect.Bool, enc.EncodeBool)
// 	rg.SetKindEncoder(reflect.Int, enc.EncodeInt)
// 	rg.SetKindEncoder(reflect.Int8, enc.EncodeInt)
// 	rg.SetKindEncoder(reflect.Int16, enc.EncodeInt)
// 	rg.SetKindEncoder(reflect.Int32, enc.EncodeInt)
// 	rg.SetKindEncoder(reflect.Int64, enc.EncodeInt)
// 	rg.SetKindEncoder(reflect.Uint, enc.EncodeUint)
// 	rg.SetKindEncoder(reflect.Uint8, enc.EncodeUint)
// 	rg.SetKindEncoder(reflect.Uint16, enc.EncodeUint)
// 	rg.SetKindEncoder(reflect.Uint32, enc.EncodeUint)
// 	rg.SetKindEncoder(reflect.Uint64, enc.EncodeUint)
// 	rg.SetKindEncoder(reflect.Float32, enc.EncodeFloat)
// 	rg.SetKindEncoder(reflect.Float64, enc.EncodeFloat)
// 	rg.SetKindEncoder(reflect.Ptr, enc.EncodePtr)
// 	rg.SetKindEncoder(reflect.Struct, enc.EncodeStruct)
// 	rg.SetKindEncoder(reflect.Array, enc.EncodeArray)
// 	rg.SetKindEncoder(reflect.Slice, enc.EncodeArray)
// 	// rg.SetKindEncoder(reflect.Map, enc.EncodeMap)
// 	enc.registry = rg
// }

// // EncodeByte :
// func (enc DefaultEncoders) EncodeByte(w codec.ValueWriter, v reflect.Value) error {
// 	b64 := base64.NewEncoder(base64.StdEncoding, w)
// 	b64.Write(v.Bytes())
// 	return b64.Close()
// }

// // EncodeJSONRaw :
// func (enc DefaultEncoders) EncodeJSONRaw(w codec.ValueWriter, v reflect.Value) error {
// 	if v.IsNil() {
// 		w.Write([]byte(`null`))
// 		return nil
// 	}
// 	buf := new(bytes.Buffer)
// 	if err := json.Compact(buf, v.Bytes()); err != nil {
// 		return err
// 	}
// 	if buf.Len() == 0 {
// 		w.Write([]byte(`{}`))
// 		return nil
// 	}
// 	w.Write(buf.Bytes())
// 	return nil
// }

// // EncodeTime :
// func (enc DefaultEncoders) EncodeTime(w codec.ValueWriter, v reflect.Value) error {
// 	var temp [20]byte
// 	x := v.Interface().(time.Time)
// 	w.Write(x.UTC().AppendFormat(temp[:0], timeFormat))
// 	return nil
// }

// // EncodeString :
// func (enc DefaultEncoders) EncodeString(w codec.ValueWriter, v reflect.Value) error {
// 	w.WriteString(v.String())
// 	return nil
// }

// // EncodeBool :
// func (enc DefaultEncoders) EncodeBool(w codec.ValueWriter, v reflect.Value) error {
// 	if v.Bool() {
// 		return w.WriteByte(byte(49))
// 	}
// 	return w.WriteByte(byte(48))
// }

// // EncodeInt :
// func (enc DefaultEncoders) EncodeInt(w codec.ValueWriter, v reflect.Value) error {
// 	var temp [8]byte
// 	w.Write(strconv.AppendInt(temp[:0], v.Int(), 10))
// 	return nil
// }

// // EncodeUint :
// func (enc DefaultEncoders) EncodeUint(w codec.ValueWriter, v reflect.Value) error {
// 	var temp [10]byte
// 	w.Write(strconv.AppendUint(temp[:0], v.Uint(), 10))
// 	return nil
// }

// // EncodeFloat :
// func (enc DefaultEncoders) EncodeFloat(w codec.ValueWriter, v reflect.Value) error {
// 	var temp [10]byte
// 	// switch vi := v.Interface().(type) {
// 	// case float32:
// 	// 	dec := decimal.NewFromFloat32(vi)
// 	// 	b, err := dec.MarshalText()
// 	// 	log.Println("FLoat32:", string(b), err)
// 	// 	return dec.MarshalText()
// 	// 	// case reflect.Float64:
// 	// }

// 	// // log.Println("Result : " + fmt.Sprintf("%f", f32))
// 	// log.Println(v.Float())
// 	w.Write(strconv.AppendFloat(temp[:0], v.Float(), 'f', -1, 64))
// 	// str := strconv.FormatFloat(v.Float(), 'f', -1, 64)
// 	// log.Println(str)
// 	// b = []byte(str)
// 	return nil
// }

// // EncodePtr :
// func (enc *DefaultEncoders) EncodePtr(w codec.ValueWriter, v reflect.Value) error {
// 	if v.IsNil() {
// 		return nil
// 	}
// 	v = v.Elem()
// 	encoder, err := enc.registry.LookupEncoder(v.Type())
// 	if err != nil {
// 		return err
// 	}
// 	return encoder(w, v)
// }

// // EncodeStruct :
// func (enc DefaultEncoders) EncodeStruct(w codec.ValueWriter, v reflect.Value) error {
// 	b, err := jsonb.Marshal(v)
// 	if err != nil {
// 		return err
// 	}
// 	w.Write(b)
// 	return nil
// }

// // EncodeArray :
// func (enc DefaultEncoders) EncodeArray(w codec.ValueWriter, v reflect.Value) error {
// 	// return jsonb.Marshal(v)
// 	return nil
// }

// // EncodeMap :
// func (enc DefaultEncoders) EncodeMap(w codec.ValueWriter, v reflect.Value) error {
// 	// return jsonb.Marshal(v)
// 	return nil
// }
