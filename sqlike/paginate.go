package sqlike

import (
	"context"

	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

type Paginator struct {
	table  *Table
	fields []interface{}
	values []interface{}
	action actions.FindActions
	option *options.FindOptions
	err    error
}

// Paginate :
func (tb *Table) Paginate(act actions.PaginateStatement, opts ...*options.PaginateOptions) (*Paginator, error) {
	x := new(actions.PaginateActions)
	if act != nil {
		*x = *(act.(*actions.PaginateActions))
	}
	if x.Table == "" {
		x.Table = tb.name
	}
	opt := new(options.PaginateOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	// sort by primary key
	x.Sorts = append(x.Sorts, expr.Desc(tb.pk))
	x.OrderBy(x.Sorts...)
	length := len(x.Sorts)
	fields := make([]interface{}, length, length)
	for i, sf := range x.Sorts {
		fields[i] = sf.Field
	}
	return &Paginator{
		table:  tb,
		fields: fields,
		action: x.FindActions,
		option: &opt.FindOptions,
	}, nil
}

func (pg *Paginator) NextPage(cursor interface{}) (err error) {
	if cursor != nil {
		fa := actions.FindOne().Select(pg.fields...).Where(
			expr.Equal(pg.table.pk, cursor),
		).(*actions.FindOneActions)
		fa.Limit(1)
		result := find(
			context.Background(),
			pg.table.name,
			pg.table.driver,
			pg.table.dialect,
			pg.table.logger,
			&fa.FindActions,
			&options.FindOptions{Debug: pg.option.Debug},
			0,
		)
		pg.values, err = result.nextValues()
		if err != nil {
			return
		}
	}
	return
}

// All :
func (pg *Paginator) All(results interface{}) error {
	action := pg.action
	if len(pg.values) > 0 {
		length := len(pg.fields)
		filters := make([]interface{}, 0, length)
		fields := make([]interface{}, length, length)
		for i, sf := range action.Sorts {
			v := primitive.C{}
			val := toString(pg.values[i])
			if sf.Order == primitive.Ascending {
				if sf.Field != pg.table.pk {
					filters = append(filters, expr.GreaterOrEqual(sf.Field, val))
				}
				v = expr.GreaterThan(sf.Field, val)
			} else {
				if sf.Field != pg.table.pk {
					filters = append(filters, expr.LesserOrEqual(sf.Field, val))
				}
				v = expr.LesserThan(sf.Field, val)
			}
			fields[i] = v
		}
		filters = append(filters, expr.Or(fields...))
		if len(action.Conditions) > 0 {
			action.Conditions = append(action.Conditions, primitive.And)
		}
		action.Conditions = append(action.Conditions, expr.And(filters...))
	}
	result := find(
		context.Background(),
		pg.table.name,
		pg.table.driver,
		pg.table.dialect,
		pg.table.logger,
		&action,
		pg.option,
		0,
	)
	return result.All(results)
}

func toString(v interface{}) interface{} {
	switch vi := v.(type) {
	case []byte:
		return string(vi)
	default:
		return vi
	}
}
