package rsql

import (
	"net/url"
	"strings"

	"github.com/si3nloong/sqlike/sql/expr"
)

func (p *Parser) parseSelect(values map[string]string, params *Params) (errs Errors) {
	val, ok := values[p.SelectTag]
	delete(values, p.SelectTag)
	if !ok || len(val) < 1 {
		return
	}

	paths := strings.Split(val, ",")
	for _, v := range paths {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		v, err := url.QueryUnescape(v)
		if err != nil {
			errs = append(errs, &FieldError{Name: v, Module: p.SelectTag})
			continue
		}
		f, ok := p.mapper.Names[v]
		if !ok {
			errs = append(errs, &FieldError{Name: v, Module: p.SelectTag})
			continue
		}
		if _, ok := f.Tag.LookUp("select"); !ok {
			errs = append(errs, &FieldError{Name: v, Module: p.SelectTag})
			continue
		}
		params.Selects = append(params.Selects, expr.Column(f.Name))
	}
	return
}
