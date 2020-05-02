package sqlstmt

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Stmt :
type Stmt interface {
	fmt.Stringer
	io.StringWriter
	io.ByteWriter
	AppendArg(arg interface{})
	AppendArgs(args []interface{})
}

// Formatter :
type Formatter interface {
	Format(it interface{}) string
	Var(i int) string
}

// Statement :
type Statement struct {
	strings.Builder
	start   time.Time
	elapsed time.Duration
	fmt     Formatter
	c       int
	args    []interface{}
}

// NewStatement :
func NewStatement(fmt Formatter) (sm *Statement) {
	sm = new(Statement)
	sm.fmt = fmt
	return
}

// Args :
func (sm *Statement) Args() []interface{} {
	return sm.args
}

// AppendArg :
func (sm *Statement) AppendArg(arg interface{}) {
	sm.args = append(sm.args, arg)
	sm.c = len(sm.args)
}

// AppendArgs :
func (sm *Statement) AppendArgs(args []interface{}) {
	sm.args = append(sm.args, args...)
	sm.c = len(sm.args)
}

// Format :
func (sm Statement) Format(state fmt.State, verb rune) {
	str := sm.String()
	if !state.Flag('+') {
		state.Write([]byte(str))
		return
	}
	i := 1
	args := sm.Args()
	for {
		idx := strings.Index(str, sm.fmt.Var(i))
		if idx < 0 {
			state.Write([]byte(str))
			break
		}
		state.Write([]byte(str[:idx]))
		state.Write([]byte(sm.fmt.Format(args[0])))
		str = str[idx+1:]
		args = args[1:]
		i++
	}
}

// StartTimer :
func (sm *Statement) StartTimer() {
	sm.start = time.Now()
}

// StopTimer :
func (sm *Statement) StopTimer() {
	sm.elapsed = time.Since(sm.start)
}

// TimeElapsed :
func (sm *Statement) TimeElapsed() time.Duration {
	if sm.elapsed < 0 {
		sm.StopTimer()
	}
	return sm.elapsed
}

// Reset : implement resetter as strings.Builer
func (sm *Statement) Reset() {
	sm.args = []interface{}{}
	sm.Builder.Reset()
	sm.fmt = nil
}
