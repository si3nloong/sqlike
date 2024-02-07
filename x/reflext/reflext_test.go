package reflext

import (
	"reflect"
	"testing"
	"time"

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
	t.Run("Get with unknown tag", func(t *testing.T) {
		v = tag.Get("unknown")
		require.Equal(t, "", v)

		v, ok = tag.LookUp("unknown")
		require.Equal(t, "", v)
		require.False(t, ok)
	})

	// existing tag with no value
	t.Run("Get with tag", func(t *testing.T) {
		v, ok = tag.Option("omitempty")
		require.True(t, ok)
		require.Empty(t, v)

		v, ok = tag.Option("no")
		require.False(t, ok)
		require.Empty(t, v)
	})

	t.Run("Get with tag of default value", func(t *testing.T) {
		// existing tag with value
		v, ok = tag.Option("size")
		require.True(t, ok)
		require.Equal(t, "20", v)

		v, ok = tag.Option("size")
		require.True(t, ok)
		require.Equal(t, "20", v)
	})
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
			// if multiple tag is defined, the last one will override it
			Name   string `db:"columnName" sqlike:",default=TEST,size=10"`
			Nested struct {
				embeddedStruct `db:"embed"`
				Enum           enum
			}
			Other    string    `db:"renamed"`
			DateTime time.Time `unknowntag:"dt"`
			embeddedStruct
			Override float64 `sql:"public.ID"`
		}
	)

	typeof = reflect.TypeOf(i)
	codec = getCodec(typeof, []string{"sqlike", "db", "sql"}, func(s string) string {
		return s
	})

	require.Equal(t, len(codec.Fields()), 16)

	t.Run("Properties", func(t *testing.T) {
		props := codec.Properties()
		require.Equal(t, len(props), 6)

		t.Run("Property: Name", func(t *testing.T) {
			p1 := props[0]
			require.Equal(t, "columnName", p1.Name())
			require.False(t, p1.IsNullable())
			require.False(t, p1.IsEmbedded())
			require.Equal(t, reflect.TypeOf(""), p1.Type())
			require.Equal(t, []int{0}, p1.Index())
			require.Nil(t, p1.Parent())
			require.Empty(t, p1.Children())
		})

		t.Run("Property: Other", func(t *testing.T) {
			p2 := props[2]
			require.Equal(t, "renamed", p2.Name())
			require.False(t, p2.IsNullable())
			require.False(t, p2.IsEmbedded())
			require.Equal(t, reflect.TypeOf(""), p2.Type())
			require.Equal(t, []int{2}, p2.Index())
			require.Nil(t, p2.Parent())
			require.Empty(t, p2.Children())
		})

		t.Run("Property: embeddedStruct", func(t *testing.T) {
			p4 := props[4]
			require.Equal(t, "test.id", p4.Name())
			require.False(t, p4.IsNullable())
			require.False(t, p4.IsEmbedded())
			require.Equal(t, reflect.TypeOf(int64(0)), p4.Type())
			require.Equal(t, []int{4, 0, 0}, p4.Index())
			require.NotNil(t, p4.Parent())
			require.Empty(t, p4.Children())

			// because being overwritten
			p5 := props[5]
			require.Equal(t, "public.ID", p5.Name())
			require.False(t, p5.IsNullable())
			require.False(t, p5.IsEmbedded())
			require.Equal(t, reflect.TypeOf(float64(0)), p5.Type())
			require.Equal(t, []int{5}, p5.Index())
			require.Nil(t, p5.Parent())
			require.Empty(t, p5.Children())
		})
	})

	t.Run("LookUpFieldByName with existed field", func(t *testing.T) {
		fi, ok := codec.LookUpFieldByName("renamed")
		require.True(t, ok)
		require.NotNil(t, fi)

		fi, ok = codec.LookUpFieldByName("public.ID")
		require.True(t, ok)
		require.NotNil(t, fi)

		fi, ok = codec.LookUpFieldByName("test.id")
		require.True(t, ok)
		require.NotNil(t, fi)

		fi, ok = codec.LookUpFieldByName("Nested.Enum")
		require.True(t, ok)
		require.NotNil(t, fi)

		fi, ok = codec.LookUpFieldByName("Nested.embed.test.id")
		require.True(t, ok)
		require.NotNil(t, fi)

		fi, ok = codec.LookUpFieldByName("columnName")
		require.True(t, ok)
		require.NotNil(t, fi)
		tag := fi.Tag()
		v, _ := tag.Option("default")
		require.Equal(t, "TEST", v)
		v, _ = tag.Option("size")
		require.Equal(t, "10", v)
		require.Equal(t, "columnName", fi.Name())
		require.Equal(t, "Name", tag.FieldName())
		require.Equal(t, reflect.TypeOf(""), fi.Type())
		require.Empty(t, fi.Children())
		require.Nil(t, fi.Parent())

		fi, ok = codec.LookUpFieldByName("dt")
		require.Nil(t, fi)
		require.False(t, ok)

		fi, ok = codec.LookUpFieldByName("DateTime")
		require.NotNil(t, fi)
		require.True(t, ok)
		v, _ = tag.Option("default")
		require.Equal(t, "TEST", v)
		require.Equal(t, reflect.TypeOf(time.Time{}), fi.Type())
		require.Empty(t, fi.Children())
		require.Nil(t, fi.Parent())
	})

	typeof = reflect.TypeOf(recursiveStruct{})
	codec = getCodec(typeof, []string{"sqlike"}, func(s string) string {
		return s
	})

	require.Equal(t, len(codec.fields), 2)
	require.Equal(t, len(codec.properties), 2)

	t.Run("Check existed fields", func(t *testing.T) {
		fi, ok := codec.LookUpFieldByName("Name")
		require.True(t, ok)
		require.NotNil(t, fi)

		fi, ok = codec.LookUpFieldByName("Recursive")
		require.True(t, ok)
		require.NotNil(t, fi)
	})

	t.Run("Check non-existed fields", func(t *testing.T) {
		fi, ok := codec.LookUpFieldByName("xName")
		require.False(t, ok)
		require.Nil(t, fi)

		fi, ok = codec.LookUpFieldByName("xRecursive")
		require.False(t, ok)
		require.Nil(t, fi)
	})

}

func TestParseTag(t *testing.T) {
	v := reflect.TypeOf(tagStruct{})

	t.Run("parseTag with MultipleTag", func(t *testing.T) {
		f, _ := v.FieldByName("ID")
		tag := parseTag(f, []string{"db", "sqlike"})
		// require.Empty(t, tag.Name())
		require.Equal(t, "ID", tag.FieldName())
		v, ok := tag.LookUp("db")
		require.True(t, ok)
		require.Equal(t, `,default=40`, v)
		require.Equal(t, `,default=40`, tag.Get("db"))

		v, ok = tag.Option("unknown")
		require.False(t, ok)
		require.Empty(t, v)
		v, ok = tag.Option("default")
		require.True(t, ok)
		require.Equal(t, "40", v)
	})

	t.Run("parseTag with skip tag", func(t *testing.T) {
		f, _ := v.FieldByName("Skip")
		tag := parseTag(f, []string{"db", "sqlike"})
		require.Equal(t, "-", tag.Name())
		require.Equal(t, "Skip", tag.FieldName())
		require.Equal(t, `-`, tag.Get("sqlike"))

		v, ok := tag.Option("missing")
		require.False(t, ok)
		require.Empty(t, v)
	})
}
