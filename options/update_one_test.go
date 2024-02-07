package options

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateOneOptions(t *testing.T) {
	opt := UpdateOne()

	{
		opt.SetDebug(true)
		require.True(t, opt.Debug)
	}

	{
		opt.SetDebug(false)
		require.False(t, opt.Debug)
	}
}
