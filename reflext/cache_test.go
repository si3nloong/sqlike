package reflext

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type dbStruct struct {
	Name string `db:"name"`
	skip bool
}

func TestMapper(t *testing.T) {
	var (
		mapper = NewMapperFunc("db", func(str string) string {
			return str
		})
		ok bool
	)

	result := mapper.CodecByType(reflect.TypeOf(dbStruct{}))

	// lookup an existed field
	{
		_, ok = result.LookUpFieldByName("name")
		require.True(t, ok)
	}

	// lookup unexists field
	{
		_, ok = result.LookUpFieldByName("Unknown")
		require.False(t, ok)
	}

	// lookup private field
	{
		_, ok = result.LookUpFieldByName("skip")
		require.False(t, ok)
	}

}
