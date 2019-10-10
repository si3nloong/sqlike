package reflext

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type Enum string

type normalStruct struct {
	ID      int64 `sqlike:"$Key"`
	private bool
	Name    string
	Num     int
}

type pointerStruct struct {
	ID   int64
	Name *string
	Num  *int
}

type recursiveStruct struct {
	Name      string
	Recursive *recursiveStruct
}

type tagStruct struct {
	ID   int64  `sqlike:"id"`
	Skip string `sqlike:"-"`
}

type PublicStruct struct {
	ID string
}

type embeddedStruct struct {
	tagStruct    `sqlike:"test"`
	PublicStruct `sqlike:"public"`
}

func TestCodec(t *testing.T) {
	var (
		typeof reflect.Type
		codec  *Struct
		i      struct {
			Name   string
			Nested struct {
				embeddedStruct
				Enum Enum
			}
			embeddedStruct
		}
	)

	{
		typeof = reflect.TypeOf(i)
		codec = getCodec(typeof, "sqlike", nil)

		require.Equal(t, len(codec.Fields), 13)
		require.Equal(t, len(codec.Properties), 4)
		require.NotNil(t, codec.Names["Name"])
		require.NotNil(t, codec.Names["Nested.Enum"])
	}

	{
		typeof = reflect.TypeOf(recursiveStruct{})
		codec = getCodec(typeof, "sqlike", nil)

		require.Equal(t, len(codec.Fields), 2)
		require.Equal(t, len(codec.Properties), 2)
		require.NotNil(t, codec.Names["Name"])
		require.NotNil(t, codec.Names["Recursive"])
	}
}
