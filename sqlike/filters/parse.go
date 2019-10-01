package filters

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql/expr"
)

// expression  = [ "(" ]
// ( constraint / expression )
// [ operator ( constraint / expression ) ]
// [ ")" ]
// operator    = ";" / ","

// ParseQuery :
func (p *Parser) ParseQuery(query string) (*Params, error) {
	var (
		param  = new(Params)
		err    error
		errs   Errors
		length int
		u64    uint64
		val    string
		paths  []string
		ok     bool
	)

	values := make(map[string]string)
	if err := parseRawQuery(values, query); err != nil {
		return param, append(errs, &FieldError{})
	}

	// Select fields
	{
		val, ok = values[p.SelectTag]
		if !ok || len(val) < 1 {
			goto Filter
		}

		paths = strings.Split(val, ",")
		for _, v := range paths {
			v = strings.TrimSpace(v)
			if len(v) == 0 {
				continue
			}
			v, err = url.QueryUnescape(v)
			if err != nil {
				errs = append(errs, &FieldError{})
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
			param.Selects = append(param.Selects, expr.Column(p.columnName(f)))
		}
	}

	// Filter fields
Filter:
	{
		val, ok = values[p.FilterTag]
		if !ok || len(val) < 1 {
			goto Sort
		}

		log.Println("Value :::", val)

		p.Filter.ParseFilter(param, val)
	}

	// Sort fields
Sort:
	{
		val, ok = values[p.SortTag]
		length = len(val)
		if !ok || length < 1 {
			goto Limit
		}

		paths = strings.Split(val, ",")
		for _, v := range paths {
			v = strings.TrimSpace(v)
			if len(v) == 0 {
				continue
			}
			v, err = url.QueryUnescape(v)
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
			param.Sorts = append(param.Sorts, sort)
		}
	}

	// Limit record
Limit:
	{
		val, ok = values[p.LimitTag]
		if !ok || len(val) < 1 {
			goto End
		}

		u64, err = strconv.ParseUint(val, 10, 64)
		if err != nil {
			errs = append(errs, &FieldError{Value: "invalid value, " + val, Module: p.LimitTag})
			goto End
		}
		if u64 > uint64(maxUint) {
			errs = append(errs, &FieldError{Value: fmt.Sprintf("overflow unsigned integer, %d", u64), Module: p.LimitTag})
			goto End
		}
		param.Limit = uint(u64)
		if param.Limit > p.MaxLimit {
			param.Limit = p.MaxLimit // prevent toxic limit
			errs = append(errs, &FieldError{Value: "maximum value of limit", Module: p.LimitTag})
		}
	}

End:

	log.Println("Debug ====================>")
	log.Println("Error :", len(errs), errs)
	log.Println("Select :", param.Selects)
	log.Println("Filter :", param.Filters)
	log.Println("Sort :", param.Sorts)
	log.Println("Limit :", param.Limit)

	if len(errs) > 0 {
		return param, errs
	}
	return param, nil
}

func (p *Parser) columnName(f *reflext.StructField) string {
	name, ok := f.Tag.LookUp("column")
	if ok {
		return strings.TrimSpace(name)
	}
	name = f.Name
	if p.FormatColumn != nil {
		return p.FormatColumn(name)
	}
	return name
}

func parseRawQuery(m map[string]string, query string) (err error) {
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		m[key] = value
	}
	return err
}
