package sqlike

import (
	"context"
	"errors"

	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/internal/primitive"
	"github.com/si3nloong/sqlike/v2/options"
	"github.com/si3nloong/sqlike/v2/sql/expr"
	"github.com/si3nloong/sqlike/v2/x/reflext"
)

// Paginate :
func (tb *Table) Paginate(
	ctx context.Context,
	act actions.PaginateStatement,
	opts ...*options.PaginateOptions,
) (*Paginator, error) {
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
	length := len(x.Sorts)
	fields := make([]any, length+1)
	sort := expr.Asc(tb.pk)
	if length > 0 {
		x := x.Sorts[length-1].(primitive.Sort)
		if x.Order == primitive.Descending {
			sort = expr.Desc(tb.pk)
		}
	}
	x.Sorts = append(x.Sorts, sort)
	for i, sf := range x.Sorts {
		fields[i] = sf.(primitive.Sort).Field
	}
	if x.RowCount == 1 {
		return nil, errors.New("sqlike: pagination required more than 1 limit")
	}
	if x.RowCount == 0 {
		x.RowCount = 100
	}
	// x.RowCount++
	return &Paginator{
		ctx:    ctx,
		table:  tb,
		fields: fields,
		action: x.FindActions,
		option: &opt.FindOptions,
	}, nil
}

// Paginator :
type Paginator struct {
	ctx    context.Context
	table  *Table
	fields []any
	values []any
	action actions.FindActions
	option *options.FindOptions
	err    error
}

// After: Paginate the records using after cursor.
//
//	if err := pg.After(context.Background(), 1); err != nil {
//		panic(err)
//	}
func (pg *Paginator) After(ctx context.Context, cursor any) (err error) {
	if pg.err != nil {
		return pg.err
	}
	if cursor == nil || reflext.IsZero(reflext.ValueOf(cursor)) {
		return ErrInvalidCursor
	}
	fa := actions.FindOne().
		Select(pg.fields...).
		Where(
			expr.Equal(pg.table.pk, cursor),
		).(*actions.FindOneActions)
	fa.Limit(1)
	result := find(
		ctx,
		pg.table.dbName,
		pg.table.name,
		pg.table.client.cache,
		getDriverFromContext(ctx, pg.table.driver),
		pg.table.dialect,
		pg.table.logger,
		&fa.FindActions,
		&options.FindOptions{Debug: pg.option.Debug},
	)
	// prevent memory leak
	defer result.Close()
	pg.values, err = result.nextValues()
	return
}

// All :
func (pg *Paginator) All(results any) error {
	if pg.err != nil {
		return pg.err
	}
	result := find(
		pg.ctx,
		pg.table.dbName,
		pg.table.name,
		pg.table.client.cache,
		pg.table.driver,
		pg.table.dialect,
		pg.table.logger,
		pg.buildAction(),
		pg.option,
	)
	return result.All(results)
}

func (pg *Paginator) buildAction() *actions.FindActions {
	action := pg.action
	if len(pg.values) < 1 {
		return &action
	}
	length := len(pg.fields)
	filters := make([]any, 0, length)
	fields := make([]any, 0)
	for i, sf := range action.Sorts {
		var v primitive.C
		val := toString(pg.values[i])
		x := sf.(primitive.Sort)
		// last sort record
		if i == length-1 {
			if x.Order == primitive.Ascending {
				fields = append(fields, expr.GreaterOrEqual(x.Field, val))
			} else {
				fields = append(fields, expr.LesserOrEqual(x.Field, val))
			}
			continue
		}
		if x.Order == primitive.Ascending {
			filters = append(filters, expr.GreaterOrEqual(x.Field, val))
			v = expr.GreaterThan(x.Field, val)
		} else {
			filters = append(filters, expr.LesserOrEqual(x.Field, val))
			v = expr.LesserThan(x.Field, val)
		}
		fields = append(fields, v)
	}
	filters = append(filters, expr.Or(fields...))
	if len(action.Conditions.Values) > 0 {
		action.Conditions.Values = append(action.Conditions.Values, primitive.And)
	}
	action.Conditions.Values = append(action.Conditions.Values, expr.And(filters...))
	return &action
}

func toString(v any) any {
	switch vi := v.(type) {
	case []byte:
		return string(vi)
	default:
		return vi
	}
}
