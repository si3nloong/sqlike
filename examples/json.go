package examples

import (
	"sort"
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
	err = table.DropIfExists()
	require.NoError(t, err)

	// migrate
	{
		table.MustMigrate(jsonStruct{})
	}

	jss := [...]jsonStruct{
		newJSONStruct(),
		newJSONStruct(),
		newJSONStruct(),
	}

	{
		_, err = table.Insert(&jss,
			options.Insert().
				SetDebug(true))
		require.NoError(t, err)
	}

	{
		var o struct {
			Text         string
			Message      string
			QuoteMessage string
			ObjKeys      []string
		}

		extr := expr.JSON_EXTRACT(expr.Column("Raw"), "$.message")
		err = table.FindOne(
			actions.FindOne().Select(
				expr.As(expr.JSON_QUOTE(expr.Column("Text")), "Text"),
				expr.JSON_UNQUOTE(extr),
				extr,
				expr.JSON_KEYS(expr.Column("Raw")),
			),
			options.FindOne().SetDebug(true),
		).Scan(&o.Text, &o.Message, &o.QuoteMessage, &o.ObjKeys)
		require.NoError(t, err)
		require.Equal(t, `"TEXT"`, o.Text)
		require.Equal(t, `ok`, o.Message)
		require.Equal(t, `"ok"`, o.QuoteMessage)

		sort.Strings(o.ObjKeys)
		require.ElementsMatch(t, []string{
			"amountInCents", "category", "message", "status", "type",
		}, o.ObjKeys)
	}

	// advance query
	{
		type output struct {
			Raw     string
			Message string
		}

		result, err = table.Find(
			actions.Find().
				Select(
					expr.As(
						expr.Func("REPLACE", expr.Column("Raw"), "message", "msg"),
						"Raw",
					),
					expr.As(
						expr.JSON_UNQUOTE(expr.JSONColumn("Raw", "message")),
						"Message",
					),
				).
				Where(
					expr.JSON_CONTAINS(
						expr.Column("StrArr"),
						expr.JSON_QUOTE("a"),
					),
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
		require.True(t, len(arr[0].Raw) > 3)
		require.Equal(t, "ok", arr[0].Message)
		require.NoError(t, err)
	}
}

func newJSONStruct() (js jsonStruct) {
	js.Text = "TEXT"
	js.Raw = []byte(`{"message":"ok","type":"TNG","category":"EWALLET","status":"SUCCESS","amountInCents":1000}`)
	js.StrArr = []string{"a", "b", "c", "d", "e", "f"}
	js.IntArr = []int{100, 16, -2, 88, 32, -47, 25}
	return
}
