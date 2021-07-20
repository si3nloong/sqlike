package sql

import (
	"testing"

	"github.com/si3nloong/sqlike/v2/x/primitive"
	"github.com/stretchr/testify/require"
)

func TestSelect(t *testing.T) {
	query := Select("A", "B", "C").
		From("Table").
		Where().
		OrderBy("A")

	require.Equal(t, []interface{}{"Table"}, query.Tables)
	require.Equal(t, primitive.Group{
		Values: []interface{}{
			primitive.Column{Name: "A"},
			primitive.Column{Name: "B"},
			primitive.Column{Name: "C"},
		},
	}, query.Exprs)
	// require.Equal(t, primitive.Group{}, query.Conditions)
}
