package sqlike

import (
	"context"
	"reflect"

	"github.com/si3nloong/sqlike/core"
	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/logs"
	sqlcore "github.com/si3nloong/sqlike/sqlike/sql/core"
	sqldriver "github.com/si3nloong/sqlike/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sqlike/sql/expr"
	"golang.org/x/xerrors"
)

// DestroyOne :
func (tb *Table) DestroyOne(delete interface{}) error {
	return destroyOne(
		context.Background(),
		tb.name,
		tb.driver,
		tb.dialect,
		tb.logger,
		delete,
	)
}

func destroyOne(ctx context.Context, tbName string, driver sqldriver.Driver, dialect sqlcore.Dialect, logger logs.Logger, delete interface{}) error {
	v := reflext.ValueOf(delete)
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
	f, exists := cdc.Names["$Key"]
	if !exists {
		return xerrors.New(`missing $Key field`)
	}

	x := new(actions.DeleteActions)
	x.Table = tbName
	fv := mapper.FieldByIndexesReadOnly(v, f.Index)
	x.Where(expr.Equal(f.Path, fv.Interface()))
	x.Limit(1)
	stmt, err := dialect.Delete(x)
	if err != nil {
		return err
	}
	result, err := sqldriver.Execute(
		ctx,
		driver,
		stmt,
		logger,
	)
	if err != nil {
		return err
	}
	if affected, _ := result.RowsAffected(); affected <= 0 {
		return xerrors.New("unable to modify entity")
	}
	return err
}

// DeleteMany :
func (tb *Table) DeleteMany(act actions.DeleteStatement) (int64, error) {
	return deleteMany(
		context.Background(),
		tb.name,
		tb.driver,
		tb.dialect,
		tb.logger,
		act,
	)
}

func deleteMany(ctx context.Context, tbName string, driver sqldriver.Driver, dialect sqlcore.Dialect, logger logs.Logger, act actions.DeleteStatement) (int64, error) {
	x := new(actions.DeleteActions)
	if act != nil {
		*x = *(act.(*actions.DeleteActions))
	}
	if x.Table == "" {
		x.Table = tbName
	}

	if len(x.Conditions) < 1 {
		return 0, xerrors.New("sqlike: no condition is not allow for delete")
	}

	stmt, err := dialect.Delete(x)
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
