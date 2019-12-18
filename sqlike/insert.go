package sqlike

import (
	"context"
	"database/sql"
	"reflect"

	"errors"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql/codec"
	sqldialect "github.com/si3nloong/sqlike/sql/dialect"
	sqldriver "github.com/si3nloong/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sqlike/logs"
	"github.com/si3nloong/sqlike/sqlike/options"
)

// InsertOne :
func (tb *Table) InsertOne(src interface{}, opts ...*options.InsertOneOptions) (sql.Result, error) {
	opt := new(options.InsertOneOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	return insertOne(
		context.Background(),
		tb.dbName,
		tb.name,
		tb.pk,
		tb.registry,
		tb.driver,
		tb.dialect,
		tb.logger,
		src,
		opt,
	)
}

func insertOne(ctx context.Context, dbName, tbName, pk string, registry *codec.Registry, driver sqldriver.Driver, dialect sqldialect.Dialect, logger logs.Logger, src interface{}, opt *options.InsertOneOptions) (sql.Result, error) {
	v := reflect.ValueOf(src)
	if !v.IsValid() {
		return nil, ErrInvalidInput
	}
	t := v.Type()
	if !reflext.IsKind(t, reflect.Ptr) {
		return nil, ErrUnaddressableEntity
	}

	if v.IsNil() {
		return nil, ErrNilEntity
	}

	mapper := reflext.DefaultMapper
	cdc := mapper.CodecByType(t)
	columns, fields := skipColumns(cdc.Properties, opt.Omits)
	length := len(columns)
	if length < 1 {
		return nil, ErrEmptyFields
	}

	values := make([][]interface{}, 1, 1)
	rows := make([]interface{}, length, length)
	for i, sf := range fields {
		val, err := encodeValue(mapper, registry, sf, v)
		if err != nil {
			return nil, err
		}
		rows[i] = val
	}
	values[0] = rows

	stmt := dialect.InsertInto(dbName, tbName, pk, columns, values, &opt.InsertOptions)
	return sqldriver.Execute(
		ctx,
		driver,
		stmt,
		getLogger(logger, opt.Debug),
	)
}

// Insert :
func (tb *Table) Insert(src interface{}, opts ...*options.InsertOptions) (sql.Result, error) {
	opt := new(options.InsertOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	return insertMany(
		context.Background(),
		tb.dbName,
		tb.name,
		tb.pk,
		tb.registry,
		tb.driver,
		tb.dialect,
		tb.logger,
		src,
		opt,
	)
}

func insertMany(ctx context.Context, dbName, tbName, pk string, registry *codec.Registry, driver sqldriver.Driver, dialect sqldialect.Dialect, logger logs.Logger, src interface{}, opt *options.InsertOptions) (sql.Result, error) {
	v := reflext.ValueOf(src)
	if !v.IsValid() {
		return nil, ErrInvalidInput
	}

	v = reflext.Indirect(v)
	t := v.Type()
	if !reflext.IsKind(t, reflect.Array) && !reflext.IsKind(t, reflect.Slice) {
		return nil, errors.New("sqlike: insert only support array or slice of entity")
	}

	records := v.Len()
	if records < 1 {
		return nil, ErrInvalidInput
	}

	t = reflext.Deref(t.Elem())
	if !reflext.IsKind(t, reflect.Struct) {
		return nil, ErrUnaddressableEntity
	}

	mapper := reflext.DefaultMapper
	cdc := mapper.CodecByType(t)
	columns, fields := skipColumns(cdc.Properties, opt.Omits)
	length := len(columns)
	if length < 1 {
		return nil, ErrEmptyFields
	}

	values := make([][]interface{}, records, records)
	v = reflext.Indirect(v)
	for i := 0; i < records; i++ {
		vi := reflext.Indirect(v.Index(i))
		rows := make([]interface{}, length, length)
		for j, sf := range fields {
			val, err := encodeValue(mapper, registry, sf, vi)
			if err != nil {
				return nil, err
			}
			rows[j] = val
		}
		values[i] = rows
	}
	stmt := dialect.InsertInto(dbName, tbName, pk, columns, values, opt)
	return sqldriver.Execute(
		context.Background(),
		driver,
		stmt,
		getLogger(logger, opt.Debug),
	)
}

func encodeValue(mapper *reflext.Mapper, registry *codec.Registry, sf *reflext.StructField, v reflect.Value) (interface{}, error) {
	fv := mapper.FieldByIndexesReadOnly(v, sf.Index)
	if _, ok := sf.Tag.LookUp("auto_increment"); ok && reflext.IsZero(fv) {
		return nil, nil
	}
	encoder, err := registry.LookupEncoder(fv)
	if err != nil {
		return nil, err
	}
	return encoder(sf, fv)
}
