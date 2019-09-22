package filters

import (
	"errors"
	"reflect"

	"github.com/iancoleman/strcase"
	"github.com/si3nloong/sqlike/reflext"
)

const (
	defaultLimit    = uint(20)
	defaultMaxLimit = uint(100)
	maxUint         = ^uint(0)
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
	MultiplexSort bool
	DefaultLimit  uint
	MaxLimit      uint
	// IgnoreError   bool
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
		LimitTag:     "$limit",
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

// SetStrict :
func (p *Parser) SetStrict(strict bool) *Parser {
	p.Strict = strict
	return p
}
