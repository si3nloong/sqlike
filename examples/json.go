package examples

import (
	"context"
	"database/sql"
	"encoding/json"
	"sort"
	"testing"

	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/stretchr/testify/require"
)

// JSONExamples :
func JSONExamples(ctx context.Context, t *testing.T, db *sqlike.Database) {
	var (
		err    error
		result *sqlike.Result
	)

	table := db.Table("JSON")
	err = table.DropIfExists(ctx)
	require.NoError(t, err)

	// migrate
	{
		table.MustMigrate(ctx, jsonStruct{})
	}

	jss := [...]jsonStruct{
		newJSONStruct(),
		newJSONStruct(),
		newJSONStruct(),
	}

	{
		_, err = table.Insert(
			ctx,
			&jss,
			options.Insert().
				SetDebug(true))
		require.NoError(t, err)
	}

	// JSON_EXTRACT, JSON_QUOTE, JSON_UNQUOTE and JSON_KEYS
	{
		var (
			id int64 = 1
			o  struct {
				Text         string
				Message      string
				QuoteMessage string
				ObjKeys      []string
			}
		)

		extr := expr.JSON_EXTRACT(expr.Column("Raw"), "$.message")
		err = table.FindOne(
			ctx,
			actions.FindOne().
				Select(
					expr.As(expr.JSON_QUOTE(expr.Column("Text")), "Text"),
					expr.JSON_UNQUOTE(extr),
					extr,
					expr.JSON_KEYS(expr.Column("Raw")),
				).
				Where(
					expr.Equal("$Key", id),
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

		// col = JSON_SET(col, '$.counter', JSON_EXTRACT(col, '$.counter') + 1)
		if affected, err := table.UpdateOne(
			ctx,
			actions.UpdateOne().
				Where(
					expr.Equal("$Key", id),
				).
				Set(
					expr.ColumnValue("Raw",
						expr.JSON_SET(
							expr.Column("Raw"),
							"$.amountInCents",
							expr.Raw("JSON_EXTRACT(`Raw`, '$.amountInCents') + 7689"),
						),
					),
				),
			options.UpdateOne().SetDebug(true),
		); err != nil {
			require.NoError(t, err)
		} else if affected < 1 {
			require.Greater(t, int64(0), affected)
		}

		var raw json.RawMessage
		if err := table.FindOne(
			ctx,
			actions.FindOne().
				Select("Raw").
				Where(
					expr.Equal("$Key", id),
				),
		).Scan(&raw); err != nil {
			require.NoError(t, err)
		}

		var output struct {
			Amount float64 `json:"amountInCents"`
		}

		if err := json.Unmarshal(raw, &output); err != nil {
			require.NoError(t, err)
		}

		require.Equal(t, float64(8689), output.Amount)
	}

	//
	{
		type rawData struct {
			ID           int64  `json:"id"`
			UnderscoreID int64  `json:"_id"`
			Text         string `json:"text"`
			Emoji        string `json:"emoji"`
			Flag         bool   `json:"flag"`
		}

		var (
			id     int64
			output jsonStruct
			js     = newJSONStruct()
			res    sql.Result
		)

		js.Raw = json.RawMessage(`{
			"_id": 100,
			"text": "Hello world",
			"emoji"  : "💩💩💩"
		}`)
		res, err = table.InsertOne(ctx, &js)
		require.NoError(t, err)
		id, err = res.LastInsertId()
		require.NoError(t, err)

		if err := table.FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("$Key", id),
				),
		).Decode(&output); err != nil {
			require.NoError(t, err)
		}

		var raw rawData
		err = json.Unmarshal(output.Raw, &raw)
		require.NoError(t, err)
		require.Equal(t, int64(0), raw.ID)
		require.Equal(t, int64(100), raw.UnderscoreID)
		require.Equal(t, "💩💩💩", raw.Emoji)

		// rename key `_id` to `id` of json object
		if affected, err := table.UpdateOne(
			ctx,
			actions.UpdateOne().
				Where(
					expr.Equal("$Key", id),
				).
				Set(
					expr.ColumnValue("Raw",
						expr.JSON_INSERT(
							expr.JSON_REMOVE(expr.Column("Raw"), "$._id"),
							"$.id",
							expr.JSON_EXTRACT(expr.Column("Raw"), "$._id"),
						),
					),
				),
			options.UpdateOne().SetDebug(true),
		); err != nil {
			require.NoError(t, err)
		} else if affected > 0 {
			require.Greater(t, affected, int64(0))
		}

		// reset the data
		output = jsonStruct{}
		raw = rawData{}

		if err := table.FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("$Key", id),
				),
		).Decode(&output); err != nil {
			require.NoError(t, err)
		}

		err = json.Unmarshal(output.Raw, &raw)
		require.NoError(t, err)
		require.Equal(t, int64(100), raw.ID)
		require.Equal(t, int64(0), raw.UnderscoreID)
		require.Equal(t, "💩💩💩", raw.Emoji)
	}

	// advance query
	{
		type output struct {
			Raw     string
			Message string
		}

		result, err = table.Find(
			ctx,
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
