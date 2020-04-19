package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	gouuid "github.com/satori/go.uuid"
	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql/driver"
	sqltype "github.com/si3nloong/sqlike/sql/type"
	"github.com/si3nloong/sqlike/sqlike/columns"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

// DataTyper :
type DataTyper interface {
	DataType(info driver.Info, sf *reflext.StructField) columns.Column
}

// DataTypeFunc :
type DataTypeFunc func(sf *reflext.StructField) columns.Column

// Builder :
type Builder struct {
	mutex    sync.Mutex
	typeMap  map[interface{}]sqltype.Type
	builders map[sqltype.Type]DataTypeFunc
}

// NewBuilder :
func NewBuilder() *Builder {
	sb := &Builder{
		typeMap:  make(map[interface{}]sqltype.Type),
		builders: make(map[sqltype.Type]DataTypeFunc),
	}
	sb.SetDefaultTypes()
	return sb
}

// SetType :
func (sb *Builder) SetType(it interface{}, t sqltype.Type) {
	sb.mutex.Lock()
	defer sb.mutex.Unlock()
	sb.typeMap[it] = t
}

// SetTypeBuilder :
func (sb *Builder) SetTypeBuilder(t sqltype.Type, builder DataTypeFunc) {
	sb.mutex.Lock()
	defer sb.mutex.Unlock()
	sb.builders[t] = builder
}

// LookUpType :
func (sb *Builder) LookUpType(t reflect.Type) (typ sqltype.Type, exists bool) {
	t = reflext.Deref(t)
	typ, exists = sb.typeMap[t]
	return
}

// GetColumn :
func (sb *Builder) GetColumn(info driver.Info, sf *reflext.StructField) (columns.Column, error) {
	t := reflext.Deref(sf.Type)
	v := reflect.New(t)
	if x, ok := v.Interface().(DataTyper); ok {
		return x.DataType(info, sf), nil
	}

	if x, ok := sb.typeMap[t]; ok {
		return sb.builders[x](sf), nil
	}

	if x, ok := sb.typeMap[t.Kind()]; ok {
		return sb.builders[x](sf), nil
	}

	return columns.Column{}, fmt.Errorf("schema: invalid data type support %v", t)
}

// SetDefaultTypes :
func (sb *Builder) SetDefaultTypes() {
	sb.SetType(reflect.TypeOf([]byte{}), sqltype.Byte)
	sb.SetType(reflect.TypeOf(uuid.UUID{}), sqltype.UUID)
	sb.SetType(reflect.TypeOf(gouuid.UUID{}), sqltype.UUID)
	sb.SetType(reflect.TypeOf(language.Tag{}), sqltype.String)
	sb.SetType(reflect.TypeOf(currency.Unit{}), sqltype.Char)
	sb.SetType(reflect.TypeOf(time.Time{}), sqltype.DateTime)
	sb.SetType(reflect.TypeOf(json.RawMessage{}), sqltype.JSON)
	sb.SetType(reflect.TypeOf(orb.Point{}), sqltype.Point)
	sb.SetType(reflect.TypeOf(orb.LineString{}), sqltype.LineString)
	sb.SetType(reflect.TypeOf(orb.Polygon{}), sqltype.Polygon)
	sb.SetType(reflect.TypeOf(orb.MultiPoint{}), sqltype.MultiPoint)
	sb.SetType(reflect.TypeOf(orb.MultiLineString{}), sqltype.MultiLineString)
	sb.SetType(reflect.TypeOf(orb.MultiPolygon{}), sqltype.MultiPolygon)
	sb.SetType(reflect.String, sqltype.String)
	sb.SetType(reflect.Bool, sqltype.Bool)
	sb.SetType(reflect.Int, sqltype.Int)
	sb.SetType(reflect.Int8, sqltype.Int8)
	sb.SetType(reflect.Int16, sqltype.Int16)
	sb.SetType(reflect.Int32, sqltype.Int32)
	sb.SetType(reflect.Int64, sqltype.Int64)
	sb.SetType(reflect.Uint, sqltype.Uint)
	sb.SetType(reflect.Uint8, sqltype.Uint8)
	sb.SetType(reflect.Uint16, sqltype.Uint16)
	sb.SetType(reflect.Uint32, sqltype.Uint32)
	sb.SetType(reflect.Uint64, sqltype.Uint64)
	sb.SetType(reflect.Float32, sqltype.Float32)
	sb.SetType(reflect.Float64, sqltype.Float64)
	sb.SetType(reflect.Struct, sqltype.Struct)
	sb.SetType(reflect.Array, sqltype.Array)
	sb.SetType(reflect.Slice, sqltype.Slice)
	sb.SetType(reflect.Map, sqltype.Map)
}
