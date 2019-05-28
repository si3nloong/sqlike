package internal

import (
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"bitbucket.org/SianLoong/sqlike/reflext"
	"bitbucket.org/SianLoong/sqlike/sqlike/sql/component"
	sqltype "bitbucket.org/SianLoong/sqlike/sqlike/sql/types"
	"golang.org/x/xerrors"
)

// ColumnTyper :
type ColumnTyper interface {
	DataType(driver string, sf *reflext.StructField) component.Column
}

// DataTypeFunc :
type DataTypeFunc func(sf *reflext.StructField) component.Column

// SchemaBuilder :
type SchemaBuilder struct {
	mutex    sync.Mutex
	typeMap  map[interface{}]sqltype.Type
	builders map[sqltype.Type]DataTypeFunc
}

// NewSchemaBuilder :
func NewSchemaBuilder() *SchemaBuilder {
	sb := &SchemaBuilder{
		typeMap:  make(map[interface{}]sqltype.Type),
		builders: make(map[sqltype.Type]DataTypeFunc),
	}
	sb.SetDefaultTypes()
	return sb
}

// SetType :
func (sb *SchemaBuilder) SetType(it interface{}, t sqltype.Type) {
	sb.mutex.Lock()
	defer sb.mutex.Unlock()
	sb.typeMap[it] = t
}

// SetTypeBuilder :
func (sb *SchemaBuilder) SetTypeBuilder(t sqltype.Type, builder DataTypeFunc) {
	sb.mutex.Lock()
	defer sb.mutex.Unlock()
	sb.builders[t] = builder
}

// LookUpType :
func (sb *SchemaBuilder) LookUpType(t reflect.Type) (typ sqltype.Type, exists bool) {
	t = reflext.Deref(t)
	typ, exists = sb.typeMap[t]
	return
}

// GetColumn :
func (sb *SchemaBuilder) GetColumn(sf *reflext.StructField) (component.Column, error) {
	v := sf.Zero
	if x, isOk := v.Interface().(ColumnTyper); isOk {
		return x.DataType("mysql", sf), nil
	}

	t := reflext.Deref(v.Type())
	if x, isOk := sb.typeMap[t]; isOk {
		return sb.builders[x](sf), nil
	}

	if x, isOk := sb.typeMap[t.Kind()]; isOk {
		return sb.builders[x](sf), nil
	}

	return component.Column{}, xerrors.Errorf("invalid data type support %v", t)
}

// SetDefaultTypes :
func (sb *SchemaBuilder) SetDefaultTypes() {
	sb.SetType(reflect.TypeOf([]byte{}), sqltype.Byte)
	sb.SetType(reflect.TypeOf(time.Time{}), sqltype.DateTime)
	sb.SetType(reflect.TypeOf(json.RawMessage{}), sqltype.JSON)
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
