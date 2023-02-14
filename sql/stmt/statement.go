package sqlstmt

import (
	"fmt"
	"strings"
	"time"
)

// Formatter :
type Formatter interface {
	Format(it any) string
	Var(i int) string
}

// Statement :
type Statement struct {
	*strings.Builder
	start   time.Time
	elapsed time.Duration

	// SQL formatter (different driver different formatter)
	fmt Formatter

	// arguments count
	c int

	// SQL query arguments
	args []any
}

func (sm *Statement) Pos() int {
	return sm.c
}

// Args :
func (sm *Statement) Args() []any {
	return sm.args
}

// AppendArgs :
func (sm *Statement) AppendArgs(args ...any) {
	sm.args = append(sm.args, args...)
	sm.c = len(sm.args)
}

// WriteAppendArgs :
func (sm *Statement) WriteAppendArgs(query string, args ...any) {
	sm.Builder.WriteString(query)
	sm.args = append(sm.args, args...)
	sm.c = len(sm.args)
}

// Format :
func (sm *Statement) Format(state fmt.State, verb rune) {
	if sm.fmt == nil {
		panic("missing formatter, unable to debug")
	}

	str := sm.String()
	if !state.Flag('+') {
		state.Write([]byte(str))
		return
	}

	var (
		i    = 1
		args = sm.Args()
		idx  int
	)
	for {
		idx = strings.Index(str, sm.fmt.Var(i))
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
	sm.args = nil
	sm.c = 0
	sm.start = time.Time{}
	sm.elapsed = time.Duration(0)
	sm.Builder.Reset()
}
