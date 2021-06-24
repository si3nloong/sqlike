package reflext

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type enum string

type recursiveStruct struct {
	Name      string
	Recursive *recursiveStruct
}

type tagStruct struct {
	ID   int64  `sqlike:"id,omitempty" db:",default=40"`
	Skip string `sqlike:"-"`
}

type PublicStruct struct {
	ID string
}

type embeddedStruct struct {
	tagStruct    `sqlike:"test"`
	PublicStruct `sqlike:"public"`
}

func TestStructTag(t *testing.T) {
	var (
		tag = StructTag{
			fieldName: "A",
			name:      "a",
			opts: map[string]string{
				"omitempty": "",
				"size":      "20",
				"charset":   "",
			},
		}
		v  string
		ok bool
	)

	require.Equal(t, tag.FieldName(), "A")
	require.Equal(t, tag.Name(), "a")
	require.Equal(t, tag.FieldName(), "A")

	// unexists tag
	{
		v = tag.Get("unknown")
		require.Equal(t, "", v)

		v, ok = tag.LookUp("unknown")
		require.Equal(t, "", v)
		require.False(t, ok)
	}

	// existing tag with no value
	{
		v = tag.Get("omitempty")
		require.Equal(t, "", v)

		v, ok = tag.LookUp("omitempty")
		require.Equal(t, "", v)
		require.True(t, ok)
	}

	// existing tag with value
	{
		v = tag.Get("size")
		require.Equal(t, "20", v)

		v, ok = tag.LookUp("size")
		require.Equal(t, "20", v)
		require.True(t, ok)
	}
}

func TestStructField(t *testing.T) {
	var (
		sf = StructField{
			id:    "",
			idx:   []int{0, 5},
			name:  "Name",
			path:  "",
			t:     reflect.TypeOf(""),
			null:  false,
			embed: false,
		}
	)

	require.Equal(t, []int{0, 5}, sf.Index())
	require.Equal(t, reflect.TypeOf("str"), sf.Type())
	require.Nil(t, sf.Children())
	require.Nil(t, sf.Parent())
	require.False(t, sf.IsNullable())
	require.False(t, sf.IsEmbedded())
}

func TestCodec(t *testing.T) {
	var (
		typeof reflect.Type
		codec  *Struct
		i      struct {
			Name   string `db:"columnName" sqlike:",default=TEST"`
			Nested struct {
				embeddedStruct
				Enum enum
			}
			embeddedStruct
		}
	)

	{
		typeof = reflect.TypeOf(i)
		codec = getCodec(typeof, []string{"sqlike", "db"}, nil)

		require.Equal(t, len(codec.fields), 13)
		require.Equal(t, len(codec.properties), 4)
		require.NotNil(t, codec.names["columnName"])
		v, _ := codec.names["columnName"].Tag().LookUp("default")
		require.Equal(t, v, "TEST")
		require.NotNil(t, codec.names["Nested.Enum"])
	}

	{
		typeof = reflect.TypeOf(recursiveStruct{})
		codec = getCodec(typeof, []string{"sqlike"}, nil)

		require.Equal(t, len(codec.fields), 2)
		require.Equal(t, len(codec.properties), 2)
		require.NotNil(t, codec.names["Name"])
		require.NotNil(t, codec.names["Recursive"])
	}
}
