package dialect

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterDialect(t *testing.T) {
	t.Run("RegisterDialect with nil dialect should panics", func(t *testing.T) {
		require.Panics(t, func() {
			RegisterDialect("", nil)
		})
	})
}

func TestGetDialectByDriver(t *testing.T) {
	t.Run("GetDialectByDriver with unknown driver, it should returns common dialect", func(t *testing.T) {
		d := GetDialectByDriver("unknown")
		require.NotNil(t, d)
		require.Contains(t, []string{"*common.commonSQL"}, reflect.TypeOf(d).String())
	})
}
