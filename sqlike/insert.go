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

	arr := reflect.MakeSlice(reflect.SliceOf(t), 0, 1)
	arr = reflect.Append(arr, v)
	return insertMany(
		context.Background(),
		tb.dbName,
		tb.name,
		tb.pk,
		tb.registry,
		tb.driver,
		tb.dialect,
		tb.logger,
		arr.Interface(),
		&opt.InsertOptions,
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

	if v.Len() < 1 {
		return nil, ErrInvalidInput
	}

	t = reflext.Deref(t.Elem())
	if !reflext.IsKind(t, reflect.Struct) {
		return nil, ErrUnaddressableEntity
	}

	mapper := reflext.DefaultMapper
	cdc := mapper.CodecByType(t)
	fields := skipColumns(cdc.Properties, opt.Omits)
	if len(fields) < 1 {
		return nil, ErrEmptyFields
	}

	stmt, err := dialect.InsertInto(
		dbName,
		tbName,
		pk,
		mapper,
		registry,
		fields,
		v,
		opt,
	)
	if err != nil {
		return nil, err
	}
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
