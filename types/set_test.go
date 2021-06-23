package types

import (
	"reflect"
	"testing"

	"github.com/si3nloong/sqlike/sql"
	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {

	set := Set{"a", "b", "c", "d"}

	t.Run("DataType", func(it *testing.T) {
		col := set.ColumnDataType(sql.Context("", "").
			SetField(field{
				name: "Set",
				t:    reflect.TypeOf(set),
			}))

		require.Equal(it, "Set", col.Name)
		require.Equal(it, "SET", col.DataType)
		require.Equal(it, "SET('')", col.Type)
		require.Nil(it, col.DefaultValue)
		require.True(it, col.Nullable)
		require.Equal(it, "utf8mb4", *col.Charset)
		require.Equal(it, "utf8mb4_0900_ai_ci", *col.Collation)
	})

	t.Run("driver.Valuer with nil value", func(it *testing.T) {
		var set Set
		v, err := set.Value()
		require.NoError(it, err)
		require.Nil(it, v)
	})

	t.Run("driver.Valuer with value", func(it *testing.T) {
		v, err := set.Value()
		require.NoError(it, err)
		require.Equal(it, "a,b,c,d", v)
	})

	t.Run("Scan Set with sql.Scanner", func(it *testing.T) {
		var set2 Set
		err := set2.Scan("a,b,c,d")
		require.NoError(it, err)
		require.Equal(it, set, set2)

		set2 = Set{}
		err = set2.Scan([]byte("a,b,c,d"))
		require.NoError(t, err)
		require.Equal(t, set, set2)
	})

}
