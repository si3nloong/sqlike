package expr

import (
	"encoding/json"
	"testing"

	"github.com/si3nloong/sqlike/v2/x/primitive"
	"github.com/stretchr/testify/require"
)

func TestJSON(t *testing.T) {
	var (
		it interface{}
	)

	t.Run("JSON_QUOTE", func(tst *testing.T) {
		it = JSON_QUOTE("a")
		require.Equal(tst, primitive.JSONFunc{
			Type: primitive.JSON_QUOTE,
			Args: []interface{}{
				primitive.Value{Raw: "a"},
			},
		}, it)
	})

	t.Run("JSON_CONTAINS", func(tst *testing.T) {
		it = JSON_CONTAINS(Column("a"), Column("b"))
		require.Equal(tst, primitive.JSONFunc{
			Type: primitive.JSON_CONTAINS,
			Args: []interface{}{
				primitive.Column{Name: "a"},
				primitive.Column{Name: "b"},
			},
		}, it)

		it = JSON_CONTAINS(`["a", "b"]`, Column("b"))
		require.Equal(tst, primitive.JSONFunc{
			Type: primitive.JSON_CONTAINS,
			Args: []interface{}{
				primitive.Value{Raw: `["a", "b"]`},
				primitive.Column{Name: "b"},
			},
		}, it)

		raw := json.RawMessage(`["A","B","C"]`)
		it = JSON_CONTAINS(raw, Column("b"))
		require.Equal(tst, primitive.JSONFunc{
			Type: primitive.JSON_CONTAINS,
			Args: []interface{}{
				primitive.Value{Raw: raw},
				primitive.Column{Name: "b"},
			},
		}, it)
	})
}
