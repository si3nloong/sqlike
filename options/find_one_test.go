package options

import (
	"testing"

	"github.com/si3nloong/sqlike/v2/x/primitive"
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
			opt.SetLockMode(LockForShare())
			require.Equal(it, primitive.Lock{Type: primitive.LockForShare}, opt.Lock)
		}

		{
			opt.SetLockMode(LockForUpdate())
			require.Equal(it, primitive.Lock{Type: primitive.LockForUpdate}, opt.Lock)
		}

		{
			// default lock
			optOne := FindOne()
			require.Equal(it, primitive.Lock{}, optOne.Lock)
		}
	})
}
