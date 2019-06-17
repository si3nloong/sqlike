package sqlike

import (
	"database/sql"
	"reflect"

	"github.com/si3nloong/sqlike/core"
	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/si3nloong/sqlike/sqlike/sql/codec"
	sqlcore "github.com/si3nloong/sqlike/sqlike/sql/core"
	sqldriver "github.com/si3nloong/sqlike/sqlike/sql/driver"
	"golang.org/x/xerrors"
)

// InsertOne :
func (tb *Table) InsertOne(src interface{}, opts ...*options.InsertOneOptions) (sql.Result, error) {
	return insertOne(
		tb.name,
		tb.driver,
		tb.dialect,
		tb.logger,
		src,
		opts,
	)
}

func insertOne(tbName string, driver sqldriver.Driver, dialect sqlcore.Dialect, logger Logger, src interface{}, opts []*options.InsertOneOptions) (sql.Result, error) {
	v := reflect.ValueOf(src)
	t := v.Type()
	if !reflext.IsKind(t, reflect.Ptr) {
		return nil, ErrUnaddressableEntity
	}

	if v.IsNil() {
		return nil, xerrors.New("entity is nil")
	}

	mapper := core.DefaultMapper
	cdc := mapper.CodecByType(t)
	columns, fields := skipGeneratedColumns(cdc.NameFields)
	length := len(columns)
	if length < 1 {
		return nil, ErrEmptyFields
	}

	values := make([][]interface{}, 1, 1)
	rows := make([]interface{}, length, length)
	for i, sf := range fields {
		val, err := encodeValue(mapper, codec.DefaultRegistry, sf, v)
		if err != nil {
			return nil, err
		}
		rows[i] = val
	}
	values[0] = rows

	opt := new(options.InsertOneOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}

	stmt := dialect.InsertInto(tbName, columns, values, opt)
	return sqldriver.Execute(
		driver,
		stmt,
		logger,
	)
}

// InsertMany :
func (tb *Table) InsertMany(srcs interface{}, opts ...*options.InsertManyOptions) (sql.Result, error) {
	v := reflext.Indirect(reflect.ValueOf(srcs))
	t := v.Type()
	if !reflext.IsKind(t, reflect.Array) && !reflext.IsKind(t, reflect.Slice) {
		return nil, xerrors.New("InsertMany only support array or slice of entity")
	}

	if v.Len() < 1 {
		return nil, xerrors.New("empty entity")
	}

	t = reflext.Deref(t.Elem())
	if !reflext.IsKind(t, reflect.Struct) {
		return nil, ErrUnaddressableEntity
	}

	mapper := core.DefaultMapper
	cdc := mapper.CodecByType(t)
	columns, fields := skipGeneratedColumns(cdc.NameFields)
	length := len(columns)
	if length < 1 {
		return nil, ErrEmptyFields
	}

	values := make([][]interface{}, 0)
	v = reflext.Indirect(v)
	count := v.Len()
	for i := 0; i < count; i++ {
		vi := reflext.Indirect(v.Index(i))
		rows := make([]interface{}, length, length)
		for j, sf := range fields {
			val, err := encodeValue(mapper, codec.DefaultRegistry, sf, vi)
			if err != nil {
				return nil, err
			}
			rows[j] = val
		}
		values = append(values, rows)
	}

	opt := new(options.InsertManyOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}

	stmt := tb.dialect.InsertInto(tb.name, columns, values, opt)
	return sqldriver.Execute(
		tb.driver,
		stmt,
		tb.logger,
	)
}

func encodeValue(mapper *reflext.Mapper, registry *codec.Registry, sf *reflext.StructField, v reflect.Value) (interface{}, error) {
	fv := mapper.FieldByIndexesReadOnly(v, sf.Index)
	if _, isOk := sf.Tag.LookUp("auto_increment"); isOk && reflext.IsZero(fv) {
		return nil, nil
	}
	encoder, err := registry.LookupEncoder(fv.Type())
	if err != nil {
		return nil, err
	}
	return encoder(sf, fv)
}
