package options

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteOptions(t *testing.T) {
	opt := Delete()

	{
		opt.SetDebug(true)
		require.True(t, opt.Debug)
	}

	{
		opt.SetDebug(false)
		require.False(t, opt.Debug)
	}
}
