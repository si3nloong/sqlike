package primitive

import (
	"fmt"
	"strings"
)

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
	Value    interface{}
	DataType DataType
}

// Func :
type Func struct {
	Name      string
	Arguments []interface{}
}

// JSONFunc :
type JSONFunc struct {
	Type      jsonFunction
	Arguments []interface{}
}

// Group :
type Group struct {
	Values []interface{}
}

// Col :
type Col string

// R :
type R struct {
	From interface{}
	To   interface{}
}

// Field :
type Field struct {
	Name   string
	Values []interface{}
}

// L :
type L struct {
	Field interface{}
	IsNot bool
	Value interface{}
}

// C :
type C struct {
	Field    interface{}
	Operator Operator
	Value    interface{}
}

// KV :
type KV struct {
	Field string
	Value interface{}
}

type operator int

// operators :
const (
	Add operator = iota
	Deduct
)

// Math :
type Math struct {
	Field Col
	Mode  operator
	Value int
}

type order int

// orders :
const (
	Ascending order = iota
	Descending
)

// Nil :
type Nil struct {
	Field interface{}
	IsNot bool
}

// Sort :
type Sort struct {
	Field string
	Order order
}

// Value :
type Value struct {
	Raw interface{}
}

// JC :
type JC struct {
	Target    interface{}
	Candidate interface{}
	Path      *string
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
	Field interface{}
	By    aggregate
}

// As :
type As struct {
	Field interface{}
	Name  string
}
