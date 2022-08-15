package reflext

import (
	"log"
	"reflect"
	"strconv"
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

type directEmbed struct {
	ID string
	No int
}

type embeddedStruct struct {
	tagStruct    `json:"test" sqlike:"test"`
	PublicStruct `json:"public" sqlike:"public"`
	directEmbed
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
			path:  "",
			t:     reflect.TypeOf(""),
			null:  false,
			embed: false,
		}
	)

	require.Equal(t, "", sf.Name())
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
			// If multiple tag is defined, the last one will override it
			// In this case, the field name will be `Name`, option value will be `default=TEST`
			Name   string `db:"columnName" sqlike:",default=TEST"`
			Nested struct {
				embeddedStruct
				Enum enum
			}
			embeddedStruct
		}

		/*
			{
				Name,
				Nested: {
					Enum
				}
				test
				public.ID,
				ID,

			}
		*/
	)

	// b, _ := json.Marshal(i)
	// log.Println(string(b))
	// panic("")

	t.Run("Parse codec with embedded struct", func(t *testing.T) {
		var (
			f  FieldInfo
			ok bool
		)

		typeof = reflect.TypeOf(i)
		codec = getCodec(typeof, []string{"sqlike", "db"}, nil)

		log.Println("debug start =======================>")
		for _, f := range codec.fields {
			log.Println(f.Name(), f.Type())
		}
		log.Println("debug ended =======================>")

		log.Println("debug start =======================>")
		for _, f := range codec.properties {
			log.Println(f.Name(), f.Type())
		}
		log.Println("debug ended =======================>")

		require.Equal(t, 19, len(codec.fields))
		require.Equal(t, 5, len(codec.properties))

		f, ok = codec.LookUpFieldByName("columnName")
		require.True(t, ok)
		require.NotNil(t, f)

		v, ok := f.Tag().LookUp("default")
		require.True(t, ok)
		require.Equal(t, "TEST", v)

		f, ok = codec.LookUpFieldByName("Nested.Enum")
		require.True(t, ok)
		require.NotNil(t, f)
	})

	t.Run("Parse codec with recursive struct", func(t *testing.T) {
		var f FieldInfo

		typeof = reflect.TypeOf(recursiveStruct{})
		codec = getCodec(typeof, []string{"sqlike"}, nil)

		require.Equal(t, len(codec.fields), 2)
		require.Equal(t, len(codec.properties), 2)

		f, _ = codec.LookUpFieldByName("Name")
		require.NotNil(t, f)

		f, _ = codec.LookUpFieldByName("Recursive")
		require.NotNil(t, f)
	})
}

func TestParseTag(t *testing.T) {

	t.Run("Parse tag with multiple tag keys", func(t *testing.T) {
		tag := parseTag(reflect.StructField{
			Tag: `db:"Name" sql:"Name2"`,
		}, []string{"db", "sql"}, strconv.Quote)

		require.Equal(t, "Name2", tag.Name())
		require.Equal(t, "Name2", tag.Name())
		require.Equal(t, "Name2", tag.Name())
		require.Equal(t, "Name2", tag.Name())
	})

}
