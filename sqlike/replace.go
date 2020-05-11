package sqlike

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sqlike/options"
)

// ReplaceOne :
func (tb *Table) ReplaceOne(ctx context.Context, src interface{}, opts ...*options.InsertOneOptions) (sql.Result, error) {
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
