package primitive

import (
	"fmt"
	"strings"
)

type JoinType int

const (
	InnerJoin JoinType = iota
	CrossJoin
	LeftJoin
	RightJoin
)

type Pair [2]string

//lint:ignore U1000 to comply it's a sql clause
func (p Pair) isClause() {}

// Raw :
type Raw struct {
	Value string
}

// JSONColumn :
type JSONColumn struct {
	Column        string
	Nested        []string
	UnquoteResult bool
}

// WithQuote :
func (x JSONColumn) WithQuote() JSONColumn {
	x.UnquoteResult = true
	return x
}

func (x JSONColumn) String() string {
	nested := strings.Join(x.Nested, ".")
	operator := "->"
	if strings.HasPrefix(nested, "$.") {
		nested = "$." + nested
	}
	if x.UnquoteResult {
		operator += ">"
	}
	return fmt.Sprintf("`%s`%s'$.%s'", x.Column, operator, nested)
}

// Column :
type Column struct {
	Table string
	Name  string
}

// Alias :
type Alias struct {
	Name  string
	Alias string
}

// CastAs :
type CastAs struct {
	Value    any
	DataType DataType
}

// Func :
type Func struct {
	Name string
	Args []any
}

// Encoding :
type Encoding struct {
	Charset *string
	Column  any
	Collate string
}

// JSONFunc :
type JSONFunc struct {
	Prefix any
	Type   jsonFunction
	Args   []any
}

type Join struct {
	Type     JoinType
	SubQuery any
	On       [2]any
}

// Group :
type Group struct {
	Values []any
}

// R :
type R struct {
	From any
	To   any
}

// Field :
type Field struct {
	Name   string
	Values []any
}

// L :
type L struct {
	Field any
	IsNot bool
	Value any
}

//lint:ignore U1000 to comply it's a sql clause
func (l L) isClause() {}

// C :
type C struct {
	Field    any
	Operator Operator
	Value    any
}

//lint:ignore U1000 to comply it's a sql clause
func (c C) isClause() {}

// KV :
type KV struct {
	Field string
	Value any
}

type operator int

// operators :
const (
	Add operator = iota
	Deduct
)

// Math :
type Math struct {
	Field string
	Mode  operator
	Value uint
}

type order int

// orders :
const (
	Ascending order = iota
	Descending
)

// Nil :
type Nil struct {
	Field any
	IsNot bool
}

// Sort :
type Sort struct {
	Field any
	Order order
}

// Value :
type Value struct {
	Raw any
}

type aggregate int

// aggregation :
const (
	Sum aggregate = iota + 1
	Count
	Average
	Max
	Min
)

// Aggregate :
type Aggregate struct {
	Field any
	By    aggregate
}

// As :
type As struct {
	Field any
	Name  string
}
