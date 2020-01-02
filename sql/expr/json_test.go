package expr

import (
	"encoding/json"
	"testing"

	"github.com/si3nloong/sqlike/sqlike/primitive"
	"github.com/stretchr/testify/require"
)

func TestJSON(t *testing.T) {
	var (
		it interface{}
	)

	// json_quote
	{
		it = JSON_QUOTE("a")
		require.Equal(t, primitive.JSONFunc{
			Type: primitive.JSONQuote,
			Arguments: []interface{}{
				primitive.Column{Name: "a"},
			},
		}, it)
	}

	// json_contains
	{
		it = JSON_CONTAINS("a", "b")
		require.Equal(t, primitive.JC{
			Target:    primitive.Column{Name: "a"},
			Candidate: primitive.Value{Raw: "b"},
		}, it)

		it = JSON_CONTAINS("a", Column("b"))
		require.Equal(t, primitive.JC{
			Target: primitive.Column{Name: "a"},
			Candidate: primitive.CastAs{
				Value:    primitive.Column{Name: "b"},
				DataType: primitive.JSON,
			},
		}, it)

		it = JSON_CONTAINS(
			json.RawMessage(`["A","B","C"]`),
			Column("b"),
		)
		require.Equal(t, primitive.JC{
			Target: primitive.Value{
				Raw: json.RawMessage(`["A","B","C"]`),
			},
			Candidate: primitive.CastAs{
				Value:    primitive.Column{Name: "b"},
				DataType: primitive.JSON,
			},
		}, it)
	}
}
