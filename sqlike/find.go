package sqlike

import (
	"database/sql"

	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/si3nloong/sqlike/sqlike/sql/codec"
	sqlcore "github.com/si3nloong/sqlike/sqlike/sql/core"
	sqldriver "github.com/si3nloong/sqlike/sqlike/sql/driver"
	"golang.org/x/xerrors"
)

// ErrUnaddressableEntity :
var ErrUnaddressableEntity = xerrors.New("unaddressable entity")

// SingleResult :
type SingleResult interface {
	Decode(dest interface{}) error
	Error() error
}

// FindOne :
func (tb *Table) FindOne(act actions.SelectOneStatement, opts ...options.FindOneOptions) SingleResult {
	x := new(actions.FindOneActions)
	if act != nil {
		*x = *(act.(*actions.FindOneActions))
	}
	x.Limit(1)
	csr := find(
		tb.name,
		tb.driver,
		tb.dialect,
		tb.logger,
		&x.FindActions,
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
func (tb *Table) Find(act actions.SelectStatement, opts ...*options.FindOptions) (*Cursor, error) {
	x := new(actions.FindActions)
	if act != nil {
		*x = *(act.(*actions.FindActions))
	}
	csr := find(
		tb.name,
		tb.driver,
		tb.dialect,
		tb.logger,
		x,
	)
	if csr.err != nil {
		return nil, csr.err
	}
	return csr, nil
}

func find(tbName string, driver sqldriver.Driver, dialect sqlcore.Dialect, logger Logger, act *actions.FindActions) *Cursor {
	if act.Table == "" {
		act.Table = tbName
	}
	csr := new(Cursor)
	csr.registry = codec.DefaultRegistry
	stmt, err := dialect.Select(act)
	if err != nil {
		csr.err = err
		return csr
	}
	rows, err := sqldriver.Query(
		driver,
		stmt,
		logger,
	)
	if err != nil {
		csr.err = err
		return csr
	}
	csr.rows = rows
	csr.columns, csr.err = rows.Columns()
	return csr
}
