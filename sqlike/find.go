package sqlike

import (
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/si3nloong/sqlike/sqlike/sql/codec"
	sqldriver "github.com/si3nloong/sqlike/sqlike/sql/driver"
	"golang.org/x/xerrors"
)

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
	csr := find(tb, &x.FindActions)
	csr.Next()
	return csr
}

// Find :
func (tb *Table) Find(act actions.SelectStatement, opts ...*options.FindOptions) (*Cursor, error) {
	x := new(actions.FindActions)
	if act != nil {
		*x = *(act.(*actions.FindActions))
	}
	csr := find(tb, x)
	if csr.err != nil {
		return nil, csr.err
	}
	return csr, nil
}

func find(tb *Table, act *actions.FindActions) *Cursor {
	if act.Table == "" {
		act.Table = tb.name
	}

	csr := new(Cursor)
	csr.registry = codec.DefaultRegistry
	stmt, err := tb.dialect.Select(act)
	if err != nil {
		csr.err = err
		return csr
	}
	rows, err := sqldriver.Query(
		tb.driver,
		stmt,
		tb.logger,
	)
	if err != nil {
		csr.err = err
		return csr
	}
	csr.rows = rows
	csr.columns, csr.err = rows.Columns()
	return csr
}
