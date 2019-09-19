package filters

import (
	"errors"
	"log"
	"net/url"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql/expr"
)

const (
	defaultLimit    = uint(20)
	defaultMaxLimit = uint(100)
)

// FormatFunc :
type FormatFunc func(string) string

// Parser :
type Parser struct {
	SelectTag     string
	FilterTag     string
	SortTag       string
	LimitTag      string
	Strict        bool
	mapper        *reflext.Struct
	v             reflect.Value
	FormatInput   FormatFunc
	FormatColumn  FormatFunc
	IgnoreError   bool
	MultiplexSort bool
	DefaultLimit  uint
	MaxLimit      uint
}

// NewParser :
func NewParser(tagName string, it interface{}) (*Parser, error) {
	t := reflext.Deref(reflect.TypeOf(it))
	if t.Kind() != reflect.Struct {
		return nil, errors.New("invalid model expected, it must be struct")
	}

	toLower := strcase.ToLowerCamel
	mapper := reflext.NewMapperFunc(tagName, nil, toLower)
	return &Parser{
		SelectTag:    "$select",
		FilterTag:    "$filter",
		SortTag:      "$sort",
		mapper:       mapper.CodecByType(t),
		v:            reflext.Zero(t),
		FormatInput:  toLower,
		FormatColumn: nil,
		DefaultLimit: defaultLimit,
		MaxLimit:     defaultMaxLimit,
	}, nil
}

// MustNewParser :
func MustNewParser(tagName string, it interface{}) *Parser {
	p, err := NewParser(tagName, it)
	if err != nil {
		panic(err)
	}
	return p
}

// ParseQuery :
func (p *Parser) ParseQuery(query string) (*Params, error) {
	values, _ := url.ParseQuery(query)
	// if err != nil {
	// return nil, err
	// }

	var (
		param  = new(Params)
		errs   Errors
		length int
		// names  []string
	)

	mod := "select"
	vals, ok := values[p.SelectTag]
	length = len(vals)
	if !ok || length < 1 {
		goto Filter
	}

	param.Selects = make([]interface{}, 0)
	vals = strings.Split(vals[length-1], ",")
	for _, v := range vals {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		f, ok := p.mapper.Names[v]
		if !ok {
			errs = append(errs, &FieldError{Name: v, Module: mod})
			continue
		}
		if _, ok := f.Tag.LookUp(mod); !ok {
			errs = append(errs, &FieldError{Name: v, Module: mod})
			continue
		}
		param.Selects = append(param.Selects, expr.Column(p.columnName(f)))
	}

Filter:
	mod = "filter"
	vals, ok = values[p.FilterTag]
	length = len(vals)
	if !ok || length < 1 {
		goto Sort
	}

	log.Println(values)

Sort:
	mod = "sort"
	vals, ok = values[p.SortTag]
	length = len(vals)
	if !ok || length < 1 {
		goto Limit
	}

	vals = strings.Split(vals[length-1], ",")
	for _, v := range vals {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		dir := v[0]
		if dir == '-' {
			v = v[1:]
		}
		f, ok := p.mapper.Names[v]
		if !ok {
			errs = append(errs, &FieldError{Name: v, Module: mod})
			continue
		}
		if _, ok := f.Tag.LookUp(mod); !ok {
			errs = append(errs, &FieldError{Name: v, Module: mod})
			continue
		}
		sort := expr.Asc(p.columnName(f))
		if dir == '-' {
			sort = expr.Desc(p.columnName(f))
		}
		param.Sorts = append(param.Sorts, sort)
	}

Limit:
	mod = "sort"
	vals, ok = values[p.LimitTag]
	length = len(vals)
	if !ok || length < 1 {
		goto End
	}

	// u64, err = strconv.ParseUint(values[length-1], 10, 64)

End:
	if param.Limit > p.MaxLimit {
		param.Limit = p.MaxLimit
	}

	log.Println("Debug ====================>")
	log.Println("Error :", len(errs), errs)
	log.Println("Select :", param.Selects)
	log.Println("Filter :", param.Filters)
	log.Println("Sort :", param.Sorts)
	log.Println("Limit :", param.Limit)

	return param, errs
}

// SetStrict :
func (p *Parser) SetStrict(strict bool) *Parser {
	p.Strict = strict
	return p
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
