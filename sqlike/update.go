package sqlike

import (
	"context"
	"reflect"

	"github.com/si3nloong/sqlike/core"
	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/logs"
	"github.com/si3nloong/sqlike/sqlike/options"
	sqlcore "github.com/si3nloong/sqlike/sqlike/sql/core"
	sqldriver "github.com/si3nloong/sqlike/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sqlike/sql/expr"
	"golang.org/x/xerrors"
)

// ModifyOne :
func (tb *Table) ModifyOne(update interface{}, opts ...*options.ModifyOneOptions) error {
	return modifyOne(
		tb.name,
		tb.dialect,
		tb.driver,
		tb.logger,
		update,
		opts,
	)
}

func modifyOne(tbName string, dialect sqlcore.Dialect, driver sqldriver.Driver, logger logs.Logger, update interface{}, opts []*options.ModifyOneOptions) error {
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

	mapper := core.DefaultMapper
	cdc := mapper.CodecByType(t)
	if _, exists := cdc.Names["$Key"]; !exists {
		return xerrors.New(`missing $Key field`)
	}

	_, fields := skipGeneratedColumns(cdc.NameFields)
	x := new(actions.UpdateActions)
	x.Table = tbName

	for _, sf := range fields {
		fv := mapper.FieldByIndexesReadOnly(v, sf.Index)
		if sf.Path == "$Key" {
			x.Where(expr.Equal("$Key", fv.Interface()))
			continue
		}
		x.Set(sf.Path, fv.Interface())
	}

	opt := new(options.ModifyOneOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}

	x.Limit(1)
	stmt, err := dialect.Update(x)
	if err != nil {
		return err
	}

	result, err := sqldriver.Execute(
		context.Background(),
		driver,
		stmt,
		getLogger(logger, opt.Debug),
	)
	if err != nil {
		return err
	}
	if affected, _ := result.RowsAffected(); affected <= 0 {
		return xerrors.New("unable to modify entity")
	}
	return err
}

// UpdateOne :
func (tb *Table) UpdateOne(act actions.UpdateOneStatement, opts ...*options.UpdateOneOptions) (int64, error) {
	x := new(actions.UpdateOneActions)
	if act != nil {
		*x = *(act.(*actions.UpdateOneActions))
	}
	if x.Table == "" {
		x.Table = tb.name
	}
	opt := new(options.UpdateOneOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}

	x.Limit(1)
	return update(
		context.Background(),
		tb.driver,
		tb.dialect,
		getLogger(tb.logger, opt.Debug),
		&x.UpdateActions,
	)
}

// UpdateMany :
func (tb *Table) UpdateMany(act actions.UpdateStatement, opts ...*options.UpdateManyOptions) (int64, error) {
	x := new(actions.UpdateActions)
	if act != nil {
		*x = *(act.(*actions.UpdateActions))
	}
	if x.Table == "" {
		x.Table = tb.name
	}
	opt := new(options.UpdateManyOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	return update(
		context.Background(),
		tb.driver,
		tb.dialect,
		getLogger(tb.logger, opt.Debug),
		x,
	)
}

func update(ctx context.Context, driver sqldriver.Driver, dialect sqlcore.Dialect, logger logs.Logger, act *actions.UpdateActions) (int64, error) {
	if len(act.Values) < 1 {
		return 0, ErrNoValueUpdate
	}
	stmt, err := dialect.Update(act)
	if err != nil {
		return 0, err
	}
	result, err := sqldriver.Execute(
		ctx,
		driver,
		stmt,
		logger,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
