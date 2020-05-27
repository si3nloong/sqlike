package options

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteOneOptions(t *testing.T) {
	opt := DeleteOne()

	{
		opt.SetDebug(true)
		require.True(t, opt.Debug)
	}

	{
		opt.SetDebug(false)
		require.False(t, opt.Debug)
	}
}
