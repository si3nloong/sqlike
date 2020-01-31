package sqlike

import (
	"context"
	"database/sql"

	"github.com/si3nloong/sqlike/sql/codec"
	sqldialect "github.com/si3nloong/sqlike/sql/dialect"
	sqldriver "github.com/si3nloong/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/logs"
	"github.com/si3nloong/sqlike/sqlike/options"
)

// SingleResult :
type SingleResult interface {
	Scan(dest ...interface{}) error
	Decode(dest interface{}) error
	Columns() []string
	ColumnTypes() ([]*sql.ColumnType, error)
	Error() error
}

// FindOne :
func (tb *Table) FindOne(ctx context.Context, act actions.SelectOneStatement, opts ...*options.FindOneOptions) SingleResult {
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
		tb.registry,
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

// Find :
func (tb *Table) Find(ctx context.Context, act actions.SelectStatement, opts ...*options.FindOptions) (*Result, error) {
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
		tb.registry,
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

func find(ctx context.Context, dbName, tbName string, registry *codec.Registry, driver sqldriver.Driver, dialect sqldialect.Dialect, logger logs.Logger, act *actions.FindActions, opt *options.FindOptions, lock options.LockMode) *Result {
	if act.Database == "" {
		act.Database = dbName
	}
	if act.Table == "" {
		act.Table = tbName
	}
	csr := new(Result)
	csr.registry = registry
	stmt, err := dialect.Select(act, lock)
	if err != nil {
		csr.err = err
		return csr
	}
	rows, err := sqldriver.Query(
		ctx,
		driver,
		stmt,
		getLogger(logger, opt.Debug),
	)
	if err != nil {
		csr.err = err
		return csr
	}
	csr.rows = rows
	csr.columnTypes, csr.err = rows.ColumnTypes()
	for _, ct := range csr.columnTypes {
		csr.columns = append(csr.columns, ct.Name())
	}
	return csr
}
