package sqlstmt

import (
	"sync"

	"github.com/si3nloong/sqlike/reflext"
	"golang.org/x/xerrors"
)

// ParseStatementFunc :
type ParseStatementFunc func(stmt *Statement, it interface{}) error

// StatementParser :
type StatementParser struct {
	mutex   sync.Mutex
	parsers map[interface{}]ParseStatementFunc
}

// NewStatementParser :
func NewStatementParser() *StatementParser {
	return &StatementParser{
		parsers: make(map[interface{}]ParseStatementFunc),
	}
}

// SetParser :
func (sp *StatementParser) SetParser(it interface{}, p ParseStatementFunc) {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()
	sp.parsers[it] = p
}

// LookupParser :
func (sp *StatementParser) LookupParser(it interface{}) (parser ParseStatementFunc, isOk bool) {
	v := reflext.ValueOf(it)
	parser, isOk = sp.parsers[v.Type()]
	return
}

// BuildStatement :
func (sp *StatementParser) BuildStatement(stmt *Statement, it interface{}) error {
	v := reflext.ValueOf(it)
	if x, isOk := sp.parsers[v.Type()]; isOk {
		return x(stmt, it)
	}
	if x, isOk := sp.parsers[v.Kind()]; isOk {
		return x(stmt, it)
	}

	return xerrors.Errorf("invalid data type support %v", v.Type())
}
