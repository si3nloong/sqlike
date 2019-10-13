package rsql

import (
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql/codec"
	"github.com/si3nloong/sqlike/util"
	"github.com/timtadh/lexmachine"
)

// FormatFunc :
type FormatFunc func(string) string

// RSQLParser :
type RSQLParser interface {
	ParseQuery(query string) (interface{}, error)
}

// Parser :
type Parser struct {
	SelectTag    string
	FilterTag    string
	SortTag      string
	LimitTag     string
	mapper       *reflext.Struct
	zero         reflect.Value
	lexer        *lexmachine.Lexer
	registry     *codec.Registry
	Parser       RSQLParser
	FormatColumn FormatFunc
	DefaultLimit uint
	MaxLimit     uint
}

// NewParser :
func NewParser(it interface{}) (*Parser, error) {
	t := reflext.Deref(reflext.TypeOf(it))
	if t.Kind() != reflect.Struct {
		return nil, errors.New("rsql: entity must be struct")
	}

	mapper := reflext.NewMapperFunc("rsql", strcase.ToLowerCamel)
	lexer := lexmachine.NewLexer()
	dl := newDefaultTokenLexer()
	dl.addActions(lexer)

	p := new(Parser)
	p.SelectTag = "$select"
	p.FilterTag = "$filter"
	p.SortTag = "$sort"
	p.LimitTag = "$limit"
	p.mapper = mapper.CodecByType(t)
	p.lexer = lexer
	p.registry = codec.DefaultRegistry
	p.Parser = dl
	p.DefaultLimit = defaultLimit
	p.MaxLimit = defaultMaxLimit
	p.zero = reflext.Zero(t)

	return p, nil
}

// MustNewParser :
func MustNewParser(it interface{}) *Parser {
	p, err := NewParser(it)
	if err != nil {
		panic(err)
	}
	return p
}

// ParseQuery :
func (p *Parser) ParseQuery(query string) (*Params, error) {
	return p.ParseQueryBytes([]byte(query))
}

// ParseQueryBytes :
func (p *Parser) ParseQueryBytes(query []byte) (*Params, error) {
	values := make(map[string]string)
	if err := parseRawQuery(values, util.UnsafeString(query)); err != nil {
		return nil, err
	}

	var (
		params = new(Params)
		err    error
		// errs   = make(Errors, 0)
	)

	log.Println(values, len(values))
	p.parseSelect(values, params)
	// p.parseFilter(values, query, params)
	p.parseSort(values, params)
	p.parseLimit(values, params)

	if err != nil {
		return nil, err
	}
	return params, nil
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
