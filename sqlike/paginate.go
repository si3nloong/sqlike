package sqlike

import (
	"context"
	"log"

	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
)

// Paginate :
func (tb *Table) Paginate(act actions.PaginateStatement, opts ...*options.PaginateOptions) (*Result, error) {
	x := new(actions.PaginateActions)
	if act != nil {
		*x = *(act.(*actions.PaginateActions))
	}

	opt := new(options.PaginateOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	log.Println("DEBUG Paginate")
	if opt.Cursor != nil {
		fa := actions.FindOne().Select().Where(
			expr.Equal(tb.pk, opt.Cursor),
		).(*actions.FindOneActions)
		result := find(
			context.Background(),
			tb.name,
			tb.driver,
			tb.dialect,
			tb.logger,
			&fa.FindActions,
			&options.FindOptions{Debug: opt.Debug},
			0,
		)
		log.Println(result)
	}

	return nil, nil
}
