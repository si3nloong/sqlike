package sqlike

import (
	"context"
	"errors"
	"reflect"

	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/options"
	"github.com/si3nloong/sqlike/v2/sql/expr"
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
	"github.com/si3nloong/sqlike/v2/x/reflext"
)

// ModifyOne :
func (tb *Table) ModifyOne(
	ctx context.Context,
	update any,
	opts ...*options.ModifyOneOptions,
) (int64, error) {
	return modifyOne(
		ctx,
		tb.dbName,
		tb.name,
		tb.pk,
		tb.client.cache,
		tb.dialect,
		tb.driver,
		tb.logger,
		update,
		opts,
	)
}

func modifyOne(
	ctx context.Context,
	dbName, tbName, pk string,
	cache reflext.StructMapper,
	dialect db.Dialect,
	driver db.Driver,
	logger db.Logger,
	update any,
	opts []*options.ModifyOneOptions,
) (int64, error) {
	v := reflext.ValueOf(update)
	if !v.IsValid() {
		return 0, ErrInvalidInput
	}

	t := v.Type()
	if !reflext.IsKind(t, reflect.Ptr) {
		return 0, ErrUnaddressableEntity
	}

	if v.IsNil() {
		return 0, ErrNilEntity
	}

	cdc := cache.CodecByType(t)
	opt := new(options.ModifyOneOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}

	fields := skipColumns(cdc.Properties(), opt.Omits)
	act := new(actions.UpdateActions)
	act.Table = tbName

	// FIXME: shouldn't use any tuple, strong type recommended
	var pkv = [2]any{}
	for _, sf := range fields {
		fv := cache.FieldByIndexesReadOnly(v, sf.Index())
		if _, ok := sf.Tag().Option("primary_key"); ok {
			if pkv[0] != nil {
				act.Set(expr.ColumnValue(pkv[0].(string), pkv[1]))
			}
			pkv[0] = sf.Name()
			pkv[1] = fv.Interface()
			continue
		}
		if sf.Name() == pk && pkv[0] == nil {
			pkv[0] = sf.Name()
			pkv[1] = fv.Interface()
			continue
		}
		act.Set(expr.ColumnValue(sf.Name(), fv.Interface()))
	}

	if pkv[0] == nil {
		return 0, errors.New("sqlike: missing primary key field")
	}

	act.Where(expr.Equal(pkv[0].(string), pkv[1])).Limit(1)
	act.Table = tbName
	act.Database = dbName

	stmt := sqlstmt.AcquireStmt(dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	if err := dialect.Update(stmt, *act); err != nil {
		return 0, err
	}

	result, err := db.Execute(
		ctx,
		getDriverFromContext(ctx, driver),
		stmt,
		getLogger(logger, opt.Debug),
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
