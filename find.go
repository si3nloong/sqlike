package sqlike

import (
	"context"
	"database/sql"

	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/options"
	sqlx "github.com/si3nloong/sqlike/v2/sql"
	"github.com/si3nloong/sqlike/v2/sql/codec"
	"github.com/si3nloong/sqlike/v2/sql/dialect"
	sqldriver "github.com/si3nloong/sqlike/v2/sql/driver"
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
	"github.com/si3nloong/sqlike/v2/sqlike/logs"
	"github.com/si3nloong/sqlike/v2/x/reflext"
)

// SingleResult : single result is an interface implementing apis as similar as driver.Result
type SingleResult interface {
	Scan(dest ...interface{}) error
	Decode(dest interface{}) error
	Columns() []string
	ColumnTypes() ([]*sql.ColumnType, error)
	Error() error
}

// FindOne : find single record on the table, you should alway check the return error to ensure it have result return.
func (tb *Table) FindOne(
	ctx context.Context,
	act actions.SelectOneStatement,
	opts ...*options.FindOneOptions,
) SingleResult {
	x := new(actions.FindOneActions)
	if act != nil {
		*x = *(act.(*actions.FindOneActions))
	}
	opt := new(options.FindOneOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	x.Limit(1)
	csr := find(
		ctx,
		tb.dbName,
		tb.name,
		tb.client.cache,
		tb.codec,
		tb.driver,
		tb.dialect,
		tb.logger,
		&x.FindActions,
		&opt.FindOptions,
		opt.FindOptions.LockMode,
	)
	csr.close = true
	if csr.err != nil {
		return csr
	}
	if !csr.Next() {
		csr.err = sql.ErrNoRows
	}
	return csr
}

// Find : find multiple records on the table.
func (tb *Table) Find(
	ctx context.Context,
	act actions.SelectStatement,
	opts ...*options.FindOptions,
) (*Result, error) {
	x := new(actions.FindActions)
	if act != nil {
		*x = *(act.(*actions.FindActions))
	}
	opt := new(options.FindOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	// has limit and limit value is zero
	if !opt.NoLimit && x.Count < 1 {
		x.Limit(100)
	}
	csr := find(
		ctx,
		tb.dbName,
		tb.name,
		tb.client.cache,
		tb.codec,
		tb.driver,
		tb.dialect,
		tb.logger,
		x,
		opt,
		opt.LockMode,
	)
	if csr.err != nil {
		return nil, csr.err
	}
	return csr, nil
}

func find(
	ctx context.Context,
	dbName, tbName string,
	cache reflext.StructMapper,
	cdc codec.Codecer,
	driver sqldriver.Driver,
	dialect dialect.Dialect,
	logger logs.Logger,
	act *actions.FindActions,
	opt *options.FindOptions,
	lock options.LockMode,
) *Result {
	if act.Database == "" {
		act.Database = dbName
	}
	if act.Table == "" {
		act.Table = tbName
	}
	rslt := new(Result)
	rslt.ctx = sqlx.Context(act.Database, act.Table)
	rslt.cache = cache
	rslt.codec = cdc

	stmt := sqlstmt.AcquireStmt(dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	if err := dialect.Select(stmt, act, lock); err != nil {
		rslt.err = err
		return rslt
	}
	rows, err := sqldriver.Query(
		ctx,
		getDriverFromContext(ctx, driver),
		stmt,
		getLogger(logger, opt.Debug),
	)
	if err != nil {
		rslt.err = err
		return rslt
	}
	rslt.rows = rows
	rslt.columnTypes, rslt.err = rows.ColumnTypes()
	if rslt.err != nil {
		defer rslt.rows.Close()
	}
	for _, col := range rslt.columnTypes {
		rslt.columns = append(rslt.columns, col.Name())
	}
	return rslt
}
