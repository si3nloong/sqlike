package sqlike

import (
	"context"
	"database/sql"

	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/logs"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/si3nloong/sqlike/sql/codec"
	sqldialect "github.com/si3nloong/sqlike/sql/dialect"
	sqldriver "github.com/si3nloong/sqlike/sql/driver"
)

// SingleResult :
type SingleResult interface {
	Decode(dest interface{}) error
	Error() error
}

// FindOne :
func (tb *Table) FindOne(act actions.SelectOneStatement, opts ...*options.FindOneOptions) SingleResult {
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
		context.Background(),
		tb.name,
		tb.driver,
		tb.dialect,
		tb.logger,
		&x.FindActions,
		&opt.FindOptions,
		options.NoLock,
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
func (tb *Table) Find(act actions.SelectStatement, opts ...*options.FindOptions) (*Result, error) {
	x := new(actions.FindActions)
	if act != nil {
		*x = *(act.(*actions.FindActions))
	}
	opt := new(options.FindOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	csr := find(
		context.Background(),
		tb.name,
		tb.driver,
		tb.dialect,
		tb.logger,
		x,
		opt,
		options.NoLock,
	)
	if csr.err != nil {
		return nil, csr.err
	}
	return csr, nil
}

func find(ctx context.Context, tbName string, driver sqldriver.Driver, dialect sqldialect.Dialect, logger logs.Logger, act *actions.FindActions, opt *options.FindOptions, lock options.LockMode) *Result {
	if act.Table == "" {
		act.Table = tbName
	}
	// has limit and limit value is zero
	if !opt.NoLimit && act.Record < 1 {
		act.Limit(100)
	}
	csr := new(Result)
	csr.registry = codec.DefaultRegistry
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
	csr.columns, csr.err = rows.Columns()
	return csr
}
