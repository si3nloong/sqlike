package sqlike

import (
	"context"
	"errors"
	"reflect"

	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/options"
	"github.com/si3nloong/sqlike/v2/sql/dialect"
	sqldriver "github.com/si3nloong/sqlike/v2/sql/driver"
	"github.com/si3nloong/sqlike/v2/sql/expr"
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
	"github.com/si3nloong/sqlike/v2/x/reflext"
)

// ModifyOne :
func (tb *Table) ModifyOne(
	ctx context.Context,
	update interface{},
	opts ...*options.ModifyOneOptions,
) error {
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
	dialect dialect.Dialect,
	driver sqldriver.Driver,
	logger db.Logger,
	update interface{},
	opts []*options.ModifyOneOptions,
) error {
	v := reflext.ValueOf(update)
	if !v.IsValid() {
		return ErrInvalidInput
	}

	t := v.Type()
	if !reflext.IsKind(t, reflect.Ptr) {
		return ErrUnaddressableEntity
	}

	if v.IsNil() {
		return ErrNilEntity
	}

	cdc := cache.CodecByType(t)
	opt := new(options.ModifyOneOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}

	fields := skipColumns(cdc.Properties(), opt.Omits)
	x := new(actions.UpdateActions)
	x.Table = tbName

	var pkv = [2]interface{}{}
	for _, sf := range fields {
		fv := cache.FieldByIndexesReadOnly(v, sf.Index())
		if _, ok := sf.Tag().LookUp("primary_key"); ok {
			if pkv[0] != nil {
				x.Set(expr.ColumnValue(pkv[0].(string), pkv[1]))
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
		x.Set(expr.ColumnValue(sf.Name(), fv.Interface()))
	}

	if pkv[0] == nil {
		return errors.New("sqlike: missing primary key field")
	}

	x.Where(expr.Equal(pkv[0], pkv[1]))
	x.Limit(1)
	x.Table = tbName
	x.Database = dbName

	stmt := sqlstmt.AcquireStmt(dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	if err := dialect.Update(stmt, x); err != nil {
		return err
	}

	result, err := sqldriver.Execute(
		ctx,
		getDriverFromContext(ctx, driver),
		stmt,
		getLogger(logger, opt.Debug),
	)
	if err != nil {
		return err
	}
	if !opt.NoStrict {
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affected < 1 {
			return ErrNoRecordAffected
		}
	}
	return nil
}
