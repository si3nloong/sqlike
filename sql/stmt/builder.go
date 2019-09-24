package sqlstmt

import (
	"fmt"
	"sync"

	"reflect"

	"github.com/si3nloong/sqlike/reflext"
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
func (sb *StatementBuilder) LookupBuilder(t reflect.Type) (blr BuildStatementFunc, ok bool) {
	blr, ok = sb.builders[t]
	return
}

// BuildStatement :
func (sb *StatementBuilder) BuildStatement(stmt *Statement, it interface{}) error {
	v := reflext.ValueOf(it)
	if x, ok := sb.builders[v.Type()]; ok {
		return x(stmt, it)
	}
	if x, ok := sb.builders[v.Kind()]; ok {
		return x(stmt, it)
	}

	return fmt.Errorf("invalid data type support %v", v.Type())
}
