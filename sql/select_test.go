package sql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSelect(t *testing.T) {
	query := Select("A", "B", "C").
		From("Table").
		Where().
		OrderBy("A")

	require.Equal(t, []any{"Table"}, query.Tables)
	// require.Equal(t, primitive.Group{
	// 	Values: []any{
	// 		primitive.Column{Name: "A"},
	// 		primitive.Column{Name: "B"},
	// 		primitive.Column{Name: "C"},
	// 	},
	// }, query.Exprs)
	// require.Equal(t, primitive.Group{
	// 	Values: []any{
	// 		primitive.Column{Name: "A"},
	// 	},
	// }, query.Sorts)
}
