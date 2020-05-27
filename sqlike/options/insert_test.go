package options

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInsertOptions(t *testing.T) {
	opt := Insert()

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

	t.Run("SetMode", func(it *testing.T) {
		{
			opt.SetMode(InsertIgnore)
			require.Equal(it, InsertIgnore, opt.Mode)
		}

		{
			opt.SetMode(InsertOnDuplicate)
			require.Equal(it, InsertOnDuplicate, opt.Mode)
		}

		// default insert mode
		{
			ot := Insert()
			require.Equal(it, insertMode(0), ot.Mode)
		}
	})

	t.Run("SetOmitFields", func(it *testing.T) {
		opt.SetOmitFields("test", "__c__")
		require.ElementsMatch(it, []string{"test", "__c__"}, opt.Omits)
	})

}
