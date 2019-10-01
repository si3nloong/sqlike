package filters

import (
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/si3nloong/sqlike/sql/codec"
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// defaultFilterParser :
type defaultFilterParser struct {
	parser   *Parser
	registry *codec.Registry
}

func (p *defaultFilterParser) ParseFilter(param *Params, val string) {
	// ([A-Za-z0-9\.\$\_\@\-]+)(\=\=|\!\=|\>\=|\>|\<\=|\<|\=\?|\!\?|\=\@|\!\@)(.+)
	var re = regexp.MustCompile(`(?i)\(?[A-Z0-9\@\#\-\_\.]+`)
	log.Println("Filter =>", val)
	log.Println("000:", re.FindStringSubmatch(val))
	return

	for len(val) > 0 {
		path := val
		log.Println("000:", re.FindStringSubmatch(path))
		// and := false
		if i := strings.IndexAny(val, ",|"); i >= 0 {
			// and = path[i] == ','
			path, val = path[:i], path[i+1:]
		} else {
			val = ""
		}
		path, _ = url.QueryUnescape(path)
		matches := re.FindStringSubmatch(path)
		if len(matches) != 4 {
			continue
		}
		field, operator, value := matches[1], matches[2], matches[3]
		// log.Println(field, operator, value)
		it, _ := p.filterValue(field, operator, value)
		param.Filters = append(param.Filters, it)
	}
	// log.Println("Final :", param.Filters)
}

func (p *defaultFilterParser) filterValue(name, operator, value string) (interface{}, *FieldError) {

	log.Println("Filter :::", name, operator, value)
	var (
		col primitive.Column
		v   interface{}
	)

	f, ok := p.parser.mapper.Names[name]
	if !ok {
		return nil, &FieldError{Name: name}
	}

	if _, ok := f.Tag.LookUp("filter"); !ok {
		return nil, &FieldError{Name: name}
	}

	col = expr.Column(p.parser.columnName(f))

	switch operator {
	case "==":
		v = expr.Equal(col, v)

	case "!=":
		v = expr.NotEqual(col, v)
	case ">":
		v = expr.GreaterThan(col, v)
	case "<":
		v = expr.LesserThan(col, v)
	case ">=":
		v = expr.GreaterOrEqual(col, v)
	case "<=":
		v = expr.LesserOrEqual(col, v)
	case "=?":
		v = expr.In(col, v)
	case "!?":
		v = expr.NotIn(col, v)
	case "=@":
		v = expr.Like(col, v)
	case "!@":
		v = expr.NotLike(col, v)
	default:
		return nil, &FieldError{}
	}

	return v, nil
}
