package actions

import (
	"testing"

	"github.com/si3nloong/sqlike/v2/sql/expr"
	"github.com/stretchr/testify/require"
)

func TestActions(t *testing.T) {
	require.Equal(t, new(FindOneActions), FindOne())
	require.Equal(t, new(FindActions), Find())
	require.Equal(t, new(PaginateActions), Paginate())
	require.Equal(t, new(UpdateOneActions), UpdateOne())
	require.Equal(t, new(UpdateActions), Update())
	require.Equal(t, new(DeleteOneActions), DeleteOne())
	require.Equal(t, new(DeleteActions), Delete())

	foAction := new(FindOneActions)
	foAction.OrderBy(expr.Asc("Aa"), expr.Desc("Bb"), expr.Desc("Cc"))
	require.ElementsMatch(t, []interface{}{
		expr.Asc("Aa"),
		expr.Desc("Bb"),
		expr.Desc("Cc"),
	}, foAction.Sorts)

	dlAction := new(DeleteActions)
	dlAction.Limit(12)
	require.Equal(t, uint(12), dlAction.Record)
	dlAction.OrderBy(expr.Asc("A"), expr.Desc("B"))
	require.ElementsMatch(t, []interface{}{
		expr.Asc("A"),
		expr.Desc("B"),
	}, dlAction.Sorts)
}
