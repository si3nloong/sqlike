package sqlstmt

import (
	"fmt"
	"sync"

	"reflect"

	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/x/reflext"
)

// BuildStatementFunc :
type BuildStatementFunc func(stmt db.Stmt, it any) error

// StatementBuilder :
type StatementBuilder struct {
	mu       sync.Mutex
	builders map[any]BuildStatementFunc
}

// NewStatementBuilder :
func NewStatementBuilder() *StatementBuilder {
	return &StatementBuilder{
		builders: make(map[any]BuildStatementFunc),
	}
}

// SetBuilder :
func (sb *StatementBuilder) SetBuilder(it any, p BuildStatementFunc) {
	sb.mu.Lock()
	sb.builders[it] = p
	sb.mu.Unlock()
}

// LookupBuilder :
func (sb *StatementBuilder) LookupBuilder(t reflect.Type) (blr BuildStatementFunc, ok bool) {
	blr, ok = sb.builders[t]
	return
}

// BuildStatement :
func (sb *StatementBuilder) BuildStatement(stmt db.Stmt, it any) error {
	v := reflext.ValueOf(it)
	if cb, ok := sb.builders[v.Type()]; ok {
		return cb(stmt, it)
	}
	if cb, ok := sb.builders[v.Kind()]; ok {
		return cb(stmt, it)
	}

	return fmt.Errorf("sqlstmt: invalid data type support %v", v.Type())
}
