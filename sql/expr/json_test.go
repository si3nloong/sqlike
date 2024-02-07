package expr

import (
	"encoding/json"
	"testing"

	"github.com/si3nloong/sqlike/v2/internal/primitive"
	"github.com/stretchr/testify/require"
)

func TestJSON(t *testing.T) {
	var (
		it any
	)

	t.Run("JSON_QUOTE", func(t *testing.T) {
		it = JSON_QUOTE("a")
		require.Equal(t, primitive.JSONFunc{
			Type: primitive.JSON_QUOTE,
			Args: []any{
				primitive.Value{Raw: "a"},
			},
		}, it)
	})

	t.Run("JSON_CONTAINS", func(t *testing.T) {
		it = JSON_CONTAINS(Column("a"), Column("b"))
		require.Equal(t, primitive.JSONFunc{
			Type: primitive.JSON_CONTAINS,
			Args: []any{
				primitive.Column{Name: "a"},
				primitive.Column{Name: "b"},
			},
		}, it)

		it = JSON_CONTAINS(`["a", "b"]`, Column("b"))
		require.Equal(t, primitive.JSONFunc{
			Type: primitive.JSON_CONTAINS,
			Args: []any{
				primitive.Value{Raw: `["a", "b"]`},
				primitive.Column{Name: "b"},
			},
		}, it)

		raw := json.RawMessage(`["A","B","C"]`)
		it = JSON_CONTAINS(raw, Column("b"))
		require.Equal(t, primitive.JSONFunc{
			Type: primitive.JSON_CONTAINS,
			Args: []any{
				primitive.Value{Raw: raw},
				primitive.Column{Name: "b"},
			},
		}, it)
	})
}
