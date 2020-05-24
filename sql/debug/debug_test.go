package debug

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDebug(t *testing.T) {
	err := ToSQL(nil)
	require.NoError(t, err)
}
