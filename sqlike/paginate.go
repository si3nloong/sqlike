package sqlike

import (
	"context"

	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

type Paginator struct {
	Result
}

func (p *Paginator) NextPage() bool {
	return true
}

// Paginate :
func (tb *Table) Paginate(act actions.PaginateStatement, opts ...*options.PaginateOptions) (*Paginator, error) {
	x := new(actions.PaginateActions)
	if act != nil {
		*x = *(act.(*actions.PaginateActions))
	}
	opt := new(options.PaginateOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	// sort by primary key
	x.Sorts = append(x.Sorts, expr.Desc(tb.pk))
	x.OrderBy(x.Sorts...)
	if opt.Cursor != nil {
		length := len(x.Sorts)
		fields := make([]interface{}, length, length)
		for i, sf := range x.Sorts {
			fields[i] = sf.Field
		}
		fa := actions.FindOne().Select(fields...).Where(
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
		values, err := result.nextValues()
		if err != nil {
			return nil, err
		}
		filters := make([]interface{}, 0)
		for i, sf := range x.Sorts {
			v := primitive.C{}
			if sf.Order == primitive.Ascending {
				if sf.Field != tb.pk {
					filters = append(filters, expr.GreaterOrEqual(sf.Field, values[i]))
				}
				v = expr.GreaterThan(sf.Field, values[i])
			} else {
				if sf.Field != tb.pk {
					filters = append(filters, expr.LesserOrEqual(sf.Field, values[i]))
				}
				v = expr.LesserThan(sf.Field, values[i])
			}
			fields[i] = v
		}
		filters = append(filters, expr.Or(fields...))
		x.Conditions = append(x.Conditions, expr.And(filters...))
	}
	result := find(
		context.Background(),
		tb.name,
		tb.driver,
		tb.dialect,
		tb.logger,
		&x.FindActions,
		&opt.FindOptions,
		0,
	)
	return &Paginator{
		Result: *result,
	}, nil
}
