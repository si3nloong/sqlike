package options

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindOne(t *testing.T) {
	opt := FindOne()

	t.Run("SetDebug", func(it *testing.T) {
		{
			opt.SetDebug(true)
			require.True(it, opt.Debug)
		}

		{
			opt.SetDebug(false)
			require.False(it, opt.Debug)
		}
	})

	t.Run("SetNoLimit", func(it *testing.T) {
		{
			opt.SetNoLimit(true)
			require.True(it, opt.NoLimit)
		}

		{
			opt.SetNoLimit(false)
			require.False(it, opt.NoLimit)
		}
	})

	t.Run("SetOmitFields", func(it *testing.T) {
		opt.SetOmitFields("A", "_underscore", "cdf")
		require.ElementsMatch(it, []string{
			"A", "_underscore", "cdf",
		}, opt.OmitFields)
	})

	t.Run("SetLockMode", func(it *testing.T) {
		{
			opt.SetLockMode(LockForRead)
			require.Equal(it, LockForRead, opt.LockMode)
		}

		{
			opt.SetLockMode(LockForUpdate)
			require.Equal(it, LockForUpdate, opt.LockMode)
		}

		{
			// default lock
			ot := FindOne()
			require.Equal(it, LockMode(0), ot.LockMode)
		}
	})
}
