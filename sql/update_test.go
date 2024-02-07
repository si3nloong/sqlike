package sql

import (
	"testing"

	"github.com/si3nloong/sqlike/v2/sql/expr"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {

	t.Run("Update on table `sqlike`", func(t *testing.T) {
		query := Update("sqlike").
			Where().
			OrderBy("A")

		require.Equal(t, "sqlike", query.Table)
	})

	t.Run("Update on table `sql`", func(t *testing.T) {
		type table string
		query := Update(table("sql")).
			Where().
			OrderBy("A")

		require.Equal(t, "sql", query.Table)
	})

	t.Run("Update on table `sqlike`.`A`", func(t *testing.T) {
		query := Update(expr.Pair("sqlike", "A")).
			Where().
			OrderBy("A")

		require.Equal(t, "sqlike", query.Database)
		require.Equal(t, "A", query.Table)
	})
}
