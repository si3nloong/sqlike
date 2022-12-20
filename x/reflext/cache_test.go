package reflext

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type dbStruct struct {
	Name  string `db:"name" sqlike:""`
	skip  bool
	Email *string
}

func TestMapper(t *testing.T) {
	var (
		mapper = NewMapperFunc(100, []string{"db"}, func(s string) string {
			return s
		})
		ok bool
	)

	require.NotNil(t, DefaultMapper())

	tmp := dbStruct{Name: "John"}
	v := reflect.ValueOf(&tmp)
	typeof := v.Type()
	fv := mapper.FieldByName(v, "name")
	require.NotNil(t, fv)

	t.Run("FieldByName with non-existed fields should panic", func(t *testing.T) {
		require.Panics(t, func() {
			mapper.FieldByName(reflect.ValueOf(0), "unknown")
			mapper.FieldByName(reflect.ValueOf(""), "unknown")
		})
		require.Panics(t, func() {
			mapper.FieldByName(v, "unknown")
		})
	})

	t.Run("TraversalsByName", func(t *testing.T) {
		require.ElementsMatch(t, [][]int{{0}}, mapper.TraversalsByName(typeof, []string{"name"}))
		require.ElementsMatch(t, [][]int{{2}}, mapper.TraversalsByName(typeof, []string{"Email"}))
	})

	// FieldByIndexesReadOnly will not initialise the field even if it's nil
	t.Run("FieldByIndexesReadOnly", func(t *testing.T) {
		fv := mapper.FieldByIndexesReadOnly(v, []int{0})
		require.Equal(t, reflect.String, fv.Kind())
		require.Equal(t, "John", fv.Interface().(string))

		fv = mapper.FieldByIndexesReadOnly(v, []int{2})
		require.Nil(t, fv.Interface())

		require.Panics(t, func() {
			mapper.FieldByIndexesReadOnly(v, []int{1000000})
		})
	})

	// FieldByIndexes will initialise if the field is nil
	t.Run("FieldByIndexes", func(t *testing.T) {
		fv := mapper.FieldByIndexes(v, []int{2})
		require.NotNil(t, fv.Interface())
		require.Equal(t, "", fv.Elem().Interface().(string))

		require.Panics(t, func() {
			mapper.FieldByIndexes(v, []int{1000000})
		})
	})

	{
		fv, ok := mapper.LookUpFieldByName(v, "name")
		require.True(t, ok)
		require.Equal(t, "John", fv.Interface().(string))

		fv, ok = mapper.LookUpFieldByName(v, "unknown")
		require.False(t, ok)
		require.Equal(t, v.Elem(), fv)
	}

	codec := mapper.CodecByType(v.Type())

	// lookup an existed field
	t.Run("LookUpFieldByName with existed field", func(t *testing.T) {
		_, ok = codec.LookUpFieldByName("name")
		require.True(t, ok)
	})

	// lookup unexists field
	t.Run("LookUpFieldByName with non-existed field", func(t *testing.T) {
		_, ok = codec.LookUpFieldByName("Unknown")
		require.False(t, ok)
	})

	// lookup private field
	t.Run("LookUpFieldByName with private field", func(t *testing.T) {
		_, ok = codec.LookUpFieldByName("skip")
		require.False(t, ok)
	})
}
