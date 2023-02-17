package sqlike

import (
	"context"
	"database/sql"

	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/options"
	sqlx "github.com/si3nloong/sqlike/v2/sql"
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
	"github.com/si3nloong/sqlike/v2/x/reflext"
)

// SingleResult : single result is an interface implementing apis as similar as driver.Result
type SingleResult interface {
	Scan(dest ...any) error
	Decode(dest any) error
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
	rslt := find(
		ctx,
		tb.dbName,
		tb.name,
		tb.client.cache,
		tb.driver,
		tb.dialect,
		tb.logger,
		&x.FindActions,
		&opt.FindOptions,
	)
	rslt.close = true
	if rslt.err != nil {
		return rslt
	}
	if !rslt.Next() {
		rslt.err = sql.ErrNoRows
	}
	return rslt
}

// Find : find multiple records on the table.
func (tb *Table) Find(
	ctx context.Context,
	act actions.SelectStatement,
	opts ...*options.FindOptions,
) (*Rows, error) {
	x := new(actions.FindActions)
	if act != nil {
		*x = *(act.(*actions.FindActions))
	}
	opt := new(options.FindOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	// has limit and limit value is zero
	if !opt.NoLimit && x.RowCount < 1 {
		x.Limit(100)
	}
	csr := find(
		ctx,
		tb.dbName,
		tb.name,
		tb.client.cache,
		tb.driver,
		tb.dialect,
		tb.logger,
		x,
		opt,
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
	driver db.Driver,
	dialect db.Dialect,
	logger db.Logger,
	act *actions.FindActions,
	opt *options.FindOptions,
) *Rows {
	if act.Database == "" {
		act.Database = dbName
	}
	if act.Table == "" {
		act.Table = tbName
	}
	rslt := new(Rows)
	rslt.ctx = sqlx.Context(act.Database, act.Table)
	rslt.cache = cache
	rslt.dialect = dialect

	stmt := sqlstmt.AcquireStmt(dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	if err := dialect.Select(stmt, act, opt.Lock); err != nil {
		rslt.err = err
		return rslt
	}
	rows, err := db.Query(
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
