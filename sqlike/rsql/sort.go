package rsql

import (
	"net/url"
	"strings"

	"github.com/si3nloong/sqlike/sql/expr"
)

func (p *Parser) parseSort(values map[string]string, params *Params) (errs Errors) {
	val, ok := values[p.SortTag]
	if !ok || len(val) < 1 {
		return nil
	}

	paths := strings.Split(val, ",")
	for _, v := range paths {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		v, err := url.QueryUnescape(v)
		if err != nil {
			errs = append(errs, &FieldError{Module: p.SortTag})
			continue
		}
		desc := v[0] == '-'
		if desc {
			v = v[1:]
		}
		f, ok := p.mapper.Names[v]
		if !ok {
			errs = append(errs, &FieldError{Name: v, Module: p.SortTag})
			continue
		}
		if _, ok := f.Tag.LookUp("sort"); !ok {
			errs = append(errs, &FieldError{Name: v, Module: p.SortTag})
			continue
		}
		sort := expr.Asc(p.columnName(f))
		if desc {
			sort = expr.Desc(p.columnName(f))
		}
		params.Sorts = append(params.Sorts, sort)
	}
	return nil
}
