package options

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaginateOptions(t *testing.T) {
	opt := Paginate()

	{
		opt.SetDebug(true)
		require.True(t, opt.Debug)
	}

	{
		opt.SetDebug(false)
		require.False(t, opt.Debug)
	}
}
