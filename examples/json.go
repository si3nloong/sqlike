package examples

import (
	"testing"

	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/stretchr/testify/require"
)

// JSONExamples :
func JSONExamples(t *testing.T, db *sqlike.Database) {
	var (
		err    error
		result *sqlike.Result
	)

	table := db.Table("JSON")
	table.DropIfExits()

	// migrate
	{
		table.MustMigrate(jsonStruct{})
	}

	jss := [...]jsonStruct{
		newJSONStruct(),
	}

	{
		_, err = table.Insert(&jss,
			options.Insert().
				SetDebug(true))
		require.NoError(t, err)
	}

	{
		result, err = table.Find(
			actions.Find().Select(
				expr.JSONQuote("Text"),
			),
			options.Find().SetDebug(true),
		)
		require.NoError(t, err)
		result.All(&jss)
	}

	// advance query
	{
		// result, err = table.Find(
		// 	actions.Find().
		// 		Where(
		// 			expr.JSONContains(
		// 				json.RawMessage(`[0,17,9]`),
		// 				expr.Column("TinyUint"),
		// 			),
		// 			expr.JSONContains(
		// 				"JSONRaw",
		// 				json.RawMessage(`{"test":"hello world"}`),
		// 			),
		// 		).
		// 		OrderBy(
		// 			expr.Desc("$Key"),
		// 		),
		// 	options.Find().
		// 		SetDebug(true).
		// 		SetNoLimit(true),
		// )
		// require.NoError(t, err)
	}
}

func newJSONStruct() (js jsonStruct) {
	js.Text = "TEXT"
	js.Raw = []byte(`{"message":"ok"}`)
	js.StrArr = []string{"a", "b", "c", "d", "e", "f"}
	js.IntArr = []int{100, 16, -2, 88, 32, -47, 25}
	return
}
