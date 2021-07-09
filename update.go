package sqlike

import (
	"context"

	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/options"
	"github.com/si3nloong/sqlike/v2/sql"
	"github.com/si3nloong/sqlike/v2/sql/dialect"
	sqldriver "github.com/si3nloong/sqlike/v2/sql/driver"
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
)

// UpdateOne :
func (tb *Table) UpdateOne(
	ctx context.Context,
	act actions.UpdateOneStatement,
	opts ...*options.UpdateOneOptions,
) (int64, error) {
	x := new(actions.UpdateOneActions)
	if act != nil {
		*x = *(act.(*actions.UpdateOneActions))
	}
	opt := new(options.UpdateOneOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}

	x.Limit(1)
	return update(
		ctx,
		tb.dbName,
		tb.name,
		tb.driver,
		tb.dialect,
		tb.logger,
		&x.UpdateActions,
		&opt.UpdateOptions,
	)
}

// Update :
func (tb *Table) Update(
	ctx context.Context,
	act actions.UpdateStatement,
	opts ...*options.UpdateOptions,
) (int64, error) {
	x := new(actions.UpdateActions)
	if act != nil {
		*x = *(act.(*actions.UpdateActions))
	}
	opt := new(options.UpdateOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	return update(
		ctx,
		tb.dbName,
		tb.name,
		tb.driver,
		tb.dialect,
		tb.logger,
		x,
		opt,
	)
}

func update(
	ctx context.Context,
	dbName, tbName string,
	driver sqldriver.Driver,
	dialect dialect.Dialect,
	logger sql.Logger,
	act *actions.UpdateActions,
	opt *options.UpdateOptions,
) (int64, error) {
	if act.Database == "" {
		act.Database = dbName
	}
	if act.Table == "" {
		act.Table = tbName
	}
	if len(act.Values) < 1 {
		return 0, ErrNoValueUpdate
	}
	stmt := sqlstmt.AcquireStmt(dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	if err := dialect.Update(stmt, act); err != nil {
		return 0, err
	}
	result, err := sqldriver.Execute(
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
