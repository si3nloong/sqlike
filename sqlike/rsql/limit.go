package rsql

import (
	"fmt"
	"strconv"
)

func (p *Parser) parseLimit(values map[string]string, params *Params) (errs Errors) {
	val, ok := values[p.LimitTag]
	delete(values, p.LimitTag)
	if !ok || len(val) < 1 {
		return nil
	}

	u64, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		errs = append(errs, &FieldError{Value: "invalid value, " + val, Module: p.LimitTag})
		return errs
	}
	if u64 > uint64(maxUint) {
		errs = append(errs, &FieldError{Value: fmt.Sprintf("overflow unsigned integer, %d", u64), Module: p.LimitTag})
		return errs
	}
	params.Limit = uint(u64)
	if params.Limit > p.MaxLimit {
		params.Limit = p.MaxLimit // prevent toxic query (limit)
		errs = append(errs, &FieldError{Value: "maximum value of limit", Module: p.LimitTag})
		return errs
	}
	return nil
}
