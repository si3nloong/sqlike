package sqlstmt

import (
	"sync"

	"reflect"

	"github.com/si3nloong/sqlike/reflext"
	"golang.org/x/xerrors"
)

// BuildStatementFunc :
type BuildStatementFunc func(stmt *Statement, it interface{}) error

// StatementBuilder :
type StatementBuilder struct {
	mutex    sync.Mutex
	builders map[interface{}]BuildStatementFunc
}

// NewStatementBuilder :
func NewStatementBuilder() *StatementBuilder {
	return &StatementBuilder{
		builders: make(map[interface{}]BuildStatementFunc),
	}
}

// SetBuilder :
func (sb *StatementBuilder) SetBuilder(it interface{}, p BuildStatementFunc) {
	sb.mutex.Lock()
	defer sb.mutex.Unlock()
	sb.builders[it] = p
}

// LookupBuilder :
func (sb *StatementBuilder) LookupBuilder(t reflect.Type) (blr BuildStatementFunc, isOk bool) {
	blr, isOk = sb.builders[t]
	return
}

// BuildStatement :
func (sb *StatementBuilder) BuildStatement(stmt *Statement, it interface{}) error {
	v := reflext.ValueOf(it)
	if x, isOk := sb.builders[v.Type()]; isOk {
		return x(stmt, it)
	}
	if x, isOk := sb.builders[v.Kind()]; isOk {
		return x(stmt, it)
	}

	return xerrors.Errorf("invalid data type support %v", v.Type())
}
