package dialect

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterDialect(t *testing.T) {
	require.Panics(t, func() {
		RegisterDialect("", nil)
	})
}

func TestGetDialectByDriver(t *testing.T) {
	require.Nil(t, GetDialectByDriver(""))
	require.Nil(t, GetDialectByDriver("unknown"))
}
