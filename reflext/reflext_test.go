package reflext

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type Enum string

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

		require.Equal(t, len(codec.fields), 13)
		require.Equal(t, len(codec.properties), 4)
		require.NotNil(t, codec.names["Name"])
		require.NotNil(t, codec.names["Nested.Enum"])
	}

	{
		typeof = reflect.TypeOf(recursiveStruct{})
		codec = getCodec(typeof, "sqlike", nil)

		require.Equal(t, len(codec.fields), 2)
		require.Equal(t, len(codec.properties), 2)
		require.NotNil(t, codec.names["Name"])
		require.NotNil(t, codec.names["Recursive"])
	}
}
