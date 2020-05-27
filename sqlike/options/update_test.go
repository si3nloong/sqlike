package options

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateOptions(t *testing.T) {
	opt := Update()

	{
		opt.SetDebug(true)
		require.True(t, opt.Debug)
	}

	{
		opt.SetDebug(false)
		require.False(t, opt.Debug)
	}
}
