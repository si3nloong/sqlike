package sqlike

import (
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	sqldriver "github.com/si3nloong/sqlike/sqlike/sql/driver"
)

// FindOne :
// func (tb *Table) FindOne(act actions.SelectStatement, opts ...options.FindOneOptions) *Result {
// 	x := new(actions.FindActions)
// 	if act != nil {
// 		*x = *(act.(*actions.FindActions))
// 	}
// 	x.Limit(1)
// 	cursor, err := tb.Find(x)
// 	return &Result{err: err, csr: cursor}
// }

// Find :
func (tb *Table) Find(act actions.SelectStatement, opts ...*options.FindOptions) (*Cursor, error) {
	x := new(actions.FindActions)
	if act != nil {
		*x = *(act.(*actions.FindActions))
	}
	if x.Table == "" {
		x.Table = tb.name
	}

	stmt, err := tb.dialect.Select(x)
	if err != nil {
		return nil, err
	}
	rows, err := sqldriver.Query(
		tb.driver,
		stmt,
		tb.logger,
	)
	if err != nil {
		return nil, err
	}
	csr := &Cursor{
		rows:     rows,
		// registry: tb.registry,
	}
	csr.columns, csr.err = rows.Columns()
	return csr, nil
}
