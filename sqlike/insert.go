package sqlike

import (
	"context"
	"database/sql"
	"reflect"

	"errors"

	"github.com/si3nloong/sqlike/sql/codec"
	sqldialect "github.com/si3nloong/sqlike/sql/dialect"
	sqldriver "github.com/si3nloong/sqlike/sql/driver"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sqlike/logs"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/si3nloong/sqlike/x/reflext"
)

// InsertOne : insert single record. You should always pass in the address of input.
func (tb *Table) InsertOne(ctx context.Context, src interface{}, opts ...*options.InsertOneOptions) (sql.Result, error) {
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
		ctx,
		tb.dbName,
		tb.name,
		tb.pk,
		tb.client.cache,
		tb.codec,
		tb.driver,
		tb.dialect,
		tb.logger,
		arr.Interface(),
		&opt.InsertOptions,
	)
}

// Insert : insert multiple records. You should always pass in the address of the slice.
func (tb *Table) Insert(ctx context.Context, src interface{}, opts ...*options.InsertOptions) (sql.Result, error) {
	opt := new(options.InsertOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	return insertMany(
		ctx,
		tb.dbName,
		tb.name,
		tb.pk,
		tb.client.cache,
		tb.codec,
		tb.driver,
		tb.dialect,
		tb.logger,
		src,
		opt,
	)
}

func insertMany(ctx context.Context, dbName, tbName, pk string, cache reflext.StructMapper, cdc codec.Codecer, driver sqldriver.Driver, dialect sqldialect.Dialect, logger logs.Logger, src interface{}, opt *options.InsertOptions) (sql.Result, error) {
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

	def := cache.CodecByType(t)
	stmt := sqlstmt.AcquireStmt(dialect)
	defer sqlstmt.ReleaseStmt(stmt)

	if err := dialect.InsertInto(
		stmt,
		dbName,
		tbName,
		pk,
		cache,
		cdc,
		def.Properties(),
		v,
		opt,
	); err != nil {
		return nil, err
	}
	return sqldriver.Execute(
		ctx,
		driver,
		stmt,
		getLogger(logger, opt.Debug),
	)
}
