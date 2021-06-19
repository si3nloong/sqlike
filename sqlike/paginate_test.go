package sqlike

import (
	"context"
	"testing"

	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/x/primitive"
	"github.com/stretchr/testify/require"
)

func TestPaginate(t *testing.T) {
	var (
		err error
		pg  *Paginator
		ctx = context.Background()
	)

	tb := Table{pk: "ID"}

	t.Run("Ascending", func(ti *testing.T) {
		pg, err = tb.Paginate(
			ctx,
			actions.Paginate().
				Where(
					expr.Equal("A", "testing"),
					expr.Between("B", 1, 100),
					expr.IsNull("C"),
				).
				OrderBy(
					expr.Asc("A"),
					expr.Asc("B"),
				),
		)

		require.NoError(ti, err)
		require.ElementsMatch(ti, []primitive.Sort{
			expr.Asc("A"),
			expr.Asc("B"),
			expr.Asc("ID"),
		}, pg.action.Sorts)
	})

	t.Run("Descending", func(ti *testing.T) {
		pg, err = tb.Paginate(
			ctx,
			actions.Paginate().
				OrderBy(
					expr.Desc("A"),
					expr.Desc("B"),
				),
		)

		require.NoError(ti, err)
		require.ElementsMatch(ti, []primitive.Sort{
			expr.Desc("A"),
			expr.Desc("B"),
			expr.Desc("ID"),
		}, pg.action.Sorts)
	})

	t.Run("Ascending & Descending", func(ti *testing.T) {
		pg, err = tb.Paginate(
			ctx,
			actions.Paginate().
				OrderBy(
					expr.Asc("A"),
					expr.Desc("B"),
					expr.Asc("C"),
				),
		)

		require.NoError(ti, err)
		require.ElementsMatch(ti, []primitive.Sort{
			expr.Asc("A"),
			expr.Desc("B"),
			expr.Asc("C"),
			expr.Asc("ID"),
		}, pg.action.Sorts)
	})

}
