package options

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModifyOneOptions(t *testing.T) {
	opt := ModifyOne()

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

	t.Run("SetOmitFields", func(it *testing.T) {
		opt.SetOmitFields("A", "cc")
		require.ElementsMatch(it, []string{"A", "cc"}, opt.Omits)
	})

	t.Run("SetStrict", func(it *testing.T) {
		opt.SetStrict(true)
		require.False(it, opt.NoStrict)

		opt.SetStrict(false)
		require.True(it, opt.NoStrict)
	})
}
