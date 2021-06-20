package options

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTransactionOptions(t *testing.T) {
	opt := Transaction()

	t.Run("SetTimeOut", func(it *testing.T) {
		timeout := time.Minute * 30
		opt.SetTimeOut(timeout)
		require.Equal(it, timeout, opt.Duration)

		timeout = time.Second * 15
		opt.SetTimeOut(timeout)
		require.Equal(it, timeout, opt.Duration)
	})

	t.Run("SetReadOnly", func(it *testing.T) {
		opt.SetReadOnly(true)
		require.True(it, opt.ReadOnly)

		opt.SetReadOnly(false)
		require.False(it, opt.ReadOnly)
	})

	t.Run("SetReadOnly", func(it *testing.T) {

		// LevelDefault         = sql.LevelDefault
		// LevelReadUncommitted = sql.LevelReadUncommitted
		// LevelReadCommitted   = sql.LevelReadCommitted
		// LevelWriteCommitted  = sql.LevelWriteCommitted
		// LevelRepeatableRead  = sql.LevelRepeatableRead
		// LevelSnapshot        = sql.LevelSnapshot
		// LevelSerializable    = sql.LevelSerializable
		// LevelLinearizable    = sql.LevelLinearizable
		opt.SetIsolationLevel(sql.LevelDefault)
		require.Equal(it, sql.LevelDefault, opt.IsolationLevel)

		opt.SetIsolationLevel(sql.LevelReadUncommitted)
		require.Equal(it, sql.LevelReadUncommitted, opt.IsolationLevel)

		opt.SetIsolationLevel(sql.LevelReadCommitted)
		require.Equal(it, sql.LevelReadCommitted, opt.IsolationLevel)

		opt.SetIsolationLevel(sql.LevelWriteCommitted)
		require.Equal(it, sql.LevelWriteCommitted, opt.IsolationLevel)

		opt.SetIsolationLevel(sql.LevelRepeatableRead)
		require.Equal(it, sql.LevelRepeatableRead, opt.IsolationLevel)

		opt.SetIsolationLevel(sql.LevelSnapshot)
		require.Equal(it, sql.LevelSnapshot, opt.IsolationLevel)

		opt.SetIsolationLevel(sql.LevelSerializable)
		require.Equal(it, sql.LevelSerializable, opt.IsolationLevel)

		opt.SetIsolationLevel(sql.LevelLinearizable)
		require.Equal(it, sql.LevelLinearizable, opt.IsolationLevel)
	})
}
