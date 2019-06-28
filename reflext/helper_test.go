package reflext

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHelper(t *testing.T) {
	var (
		ptr         *string
		multiptrint *****int
		nilSlice    []string
		nilMap      map[string]interface{}
		v           reflect.Value
	)

	{
		require.True(t, IsNull(reflect.ValueOf(ptr)))
		require.True(t, IsNull(reflect.ValueOf(nilSlice)))
		require.True(t, IsNull(reflect.ValueOf(nilMap)))
	}

	{
		rt := reflect.TypeOf(multiptrint)
		rt = Deref(rt)
		// require.Equal(t, reflect.TypeOf(int(1)), multiptrint)
	}

	{
		v = ValueOf(multiptrint)
		require.Equal(t, reflect.ValueOf(multiptrint), v)
		require.Equal(t, reflect.ValueOf(multiptrint), ValueOf(reflect.ValueOf(multiptrint)))
	}

}
