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
	err = table.DropIfExits()
	require.NoError(t, err)

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
		err = result.All(&jss)
		require.NoError(t, err)
	}

	// advance query
	{
		type output struct {
			Raw string
		}

		result, err = table.Find(
			actions.Find().
				Select(
					expr.As(
						expr.Function("REPLACE", expr.Column("Raw"), "message", "msg"),
						"Raw",
					),
				).
				Where(
					// expr.JSONContains(
					// 	"IntArr",
					// 	expr.Column("IntAtt"),
					// ),
					expr.Equal(
						expr.JSONColumn("Raw", "message"),
						"ok",
					),
				).
				OrderBy(
					expr.Desc("$Key"),
				),
			options.Find().
				SetDebug(true).
				SetNoLimit(true),
		)
		require.NoError(t, err)

		arr := []output{}
		err = result.All(&arr)
		require.True(t, len(arr) > 0)
		require.Equal(t, `{"msg": "ok"}`, arr[0].Raw)
		require.NoError(t, err)
	}
}

func newJSONStruct() (js jsonStruct) {
	js.Text = "TEXT"
	js.Raw = []byte(`{"message":"ok"}`)
	js.StrArr = []string{"a", "b", "c", "d", "e", "f"}
	js.IntArr = []int{100, 16, -2, 88, 32, -47, 25}
	return
}
