package sqlstmt

import (
	"fmt"
	"strings"
)

// Formatter :
type Formatter interface {
	Format(it interface{}) string
}

// Statement :
type Statement struct {
	strings.Builder
	fmt  Formatter
	c    int
	args []interface{}
}

// NewStatement :
func NewStatement(fmt Formatter) (stmt *Statement) {
	stmt = new(Statement)
	stmt.fmt = fmt
	return
}

// Args :
func (stmt *Statement) Args() []interface{} {
	return stmt.args
}

// AppendArg :
func (stmt *Statement) AppendArg(arg interface{}) {
	stmt.args = append(stmt.args, arg)
	stmt.c = len(stmt.args)
}

// AppendArgs :
func (stmt *Statement) AppendArgs(args []interface{}) {
	stmt.args = append(stmt.args, args...)
	stmt.c = len(stmt.args)
}

// Format :
func (stmt Statement) Format(state fmt.State, verb rune) {
	str := stmt.String()
	args := stmt.Args()
	for {
		idx := strings.Index(str, `?`)
		if idx < 0 {
			state.Write([]byte(str))
			break
		}
		state.Write([]byte(str[:idx]))
		state.Write([]byte(stmt.fmt.Format(args[0])))
		str = str[idx+1:]
		args = args[1:]
	}
	return
}
