package options

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLockMode(t *testing.T) {
	opt := FindOne().SetLockMode(LockForRead)
	require.Equal(t, LockForRead, opt.LockMode)
}
