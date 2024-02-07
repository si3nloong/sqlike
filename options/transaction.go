package options

import (
	"database/sql"
	"time"
)

// IsolationLevel :
type IsolationLevel = sql.IsolationLevel

// Various isolation levels that drivers may support in BeginTx.
// If a driver does not support a given isolation level an error may be returned.
//
// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
const (
	LevelDefault         = sql.LevelDefault
	LevelReadUncommitted = sql.LevelReadUncommitted
	LevelReadCommitted   = sql.LevelReadCommitted
	LevelWriteCommitted  = sql.LevelWriteCommitted
	LevelRepeatableRead  = sql.LevelRepeatableRead
	LevelSnapshot        = sql.LevelSnapshot
	LevelSerializable    = sql.LevelSerializable
	LevelLinearizable    = sql.LevelLinearizable
)

// Transaction :
func Transaction() *TransactionOptions {
	return &TransactionOptions{}
}

// TransactionOptions :
type TransactionOptions struct {
	Duration       time.Duration
	IsolationLevel IsolationLevel
	ReadOnly       bool
}

// SetTimeOut :
func (opts *TransactionOptions) SetTimeOut(duration time.Duration) *TransactionOptions {
	opts.Duration = duration
	return opts
}

// SetIsolationLevel :
func (opts *TransactionOptions) SetIsolationLevel(level IsolationLevel) *TransactionOptions {
	opts.IsolationLevel = level
	return opts
}

// SetReadOnly :
func (opts *TransactionOptions) SetReadOnly(readOnly bool) *TransactionOptions {
	opts.ReadOnly = readOnly
	return opts
}
