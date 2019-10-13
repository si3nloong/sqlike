package rsql

import (
	"errors"
	"log"
	"reflect"

	"github.com/iancoleman/strcase"
	"github.com/si3nloong/sqlike/reflext"
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
	mapper       *reflext.Struct
	zero         reflect.Value
	lexer        *lexmachine.Lexer
	Parser       RSQLParser
	FormatColumn FormatFunc
	DefaultLimit uint
	MaxLimit     uint
}

// NewParser :
func NewParser(it interface{}) (*Parser, error) {
	t := reflext.Deref(reflect.TypeOf(it))
	if t.Kind() != reflect.Struct {
		return nil, errors.New("rsql: entity must be struct")
	}

	mapper := reflext.NewMapperFunc("rsql", strcase.ToLowerCamel)
	lexer := lexmachine.NewLexer()
	dl := newDefaultTokenLexer()
	dl.addActions(lexer)

	p := new(Parser)
	p.mapper = mapper.CodecByType(t)
	p.lexer = lexer
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
func (p *Parser) ParseQuery(b []byte) (interface{}, error) {
	log.Println(string(b))
	lxr, err := p.lexer.Scanner(b)
	if err != nil {
		return nil, err
	}
	scan := &Scanner{Scanner: lxr}
	scan.ParseToken()
	log.Println(scan.TC, scan.Text)
	return nil, nil
}
