package examples

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/si3nloong/sqlike/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

// FindExamples :
func FindExamples(t *testing.T, db *sqlike.Database) {
	var (
		// result sql.Result
		result *sqlike.Result
		ns     normalStruct
		err    error
	)

	emoji := `🤕`
	long := `プログラミングは素晴らしい力です。
	やらないのはもったいない。
	悩んでいるなら、Progateでやってみよう。
	無料で始められる、初心者向け学習サイト。
	`

	uid, _ := uuid.Parse("e7977246-910a-11e9-844d-6c96cfd87a51")
	ts, _ := time.Parse("2006-01-02 15:04:05", "2008-01-28 10:25:33")
	b := []byte(`abcd1234`)
	lang := language.Japanese
	langs := []language.Tag{
		language.AmericanEnglish,
		language.BrazilianPortuguese,
		language.Malay,
		language.ModernStandardArabic,
		language.Korean,
		language.Japanese,
	}

	table := db.Table("NormalStruct")

	{
		ns = normalStruct{}
		ns.ID = uid
		ns.Emoji = emoji
		ns.Byte = b
		ns.LongStr = long
		ns.TinyInt = -88
		ns.SmallInt = -16829
		ns.BigInt = -1298738901289381212
		ns.Uint = 1683904243
		ns.SmallUint = 188
		ns.MediumUint = 121373123
		ns.BigUint = 1298738901289381212
		ns.Float32 = 10.6789
		ns.Float64 = 19833.6789
		ns.GeoPoint = [2]float64{10.12331, 384.7899003}
		ns.JSONRaw = json.RawMessage(`{"test":"hello world"}`)
		ns.Enum = Failed
		ns.Map = make(map[string]int)
		ns.Map["one"] = 1
		ns.Map["three"] = 3
		ns.Map["eleven"] = 11
		ns.DateTime = ts
		ns.Timestamp = ts
		ns.Language = lang
		ns.Languages = langs

		_, err = table.InsertOne(&ns)
		require.NoError(t, err)
	}

	// Find one record by primary key
	{
		ns = normalStruct{}
		err = table.FindOne(
			actions.FindOne().Where(
				expr.Equal("$Key", uid),
			),
			options.FindOne().SetDebug(true),
		).Decode(&ns)
		require.NoError(t, err)

		require.Equal(t, uid, ns.ID)
		require.Equal(t, emoji, ns.Emoji)
		require.Equal(t, long, ns.LongStr)
		require.Equal(t, ts, ns.Timestamp)
		require.Equal(t, b, ns.Byte)
		require.Equal(t, int8(-88), ns.TinyInt)
		require.Equal(t, int16(-16829), ns.SmallInt)
		require.Equal(t, int64(-1298738901289381212), ns.BigInt)
		require.Equal(t, uint(1683904243), ns.Uint)
		require.Equal(t, uint16(188), ns.SmallUint)
		require.Equal(t, uint32(121373123), ns.MediumUint)
		require.Equal(t, uint64(1298738901289381212), ns.BigUint)
		require.Equal(t, float32(10.6789), ns.Float32)
		require.Equal(t, float64(19833.6789), ns.Float64)
		require.Equal(t, Enum("FAILED"), ns.Enum)
		require.Equal(t, types.GeoPoint{10.12331, 384.7899003}, ns.GeoPoint)
		require.Equal(t, map[string]int{
			"one":    1,
			"three":  3,
			"eleven": 11,
		}, ns.Map)
		require.Equal(t, lang, ns.Language)
		require.Equal(t, langs, ns.Languages)
		require.Equal(t, json.RawMessage(`{"test":"hello world"}`), ns.JSONRaw)
	}

	// Find one record by primary key
	{
		ns = normalStruct{}
		err = table.FindOne(
			actions.FindOne().Where(
				expr.Equal("$Key", "1000"),
			),
		).Decode(&ns)
		require.Equal(t, err, sqlike.ErrNoRows)
	}

	// Find multiple records by where condition
	{
		ns = normalStruct{}
		result, err = table.Find(
			actions.Find().Where(
				expr.Between("TinyInt", 1, 100),
				expr.In("Enum", []Enum{
					Success,
					Failed,
					Unknown,
				}),
			),
			options.Find().SetDebug(true),
		)
		require.NoError(t, err)
	}

	{
		ns = normalStruct{}
		result, err = table.Find(
			actions.Find().Select("Emoji"),
			options.Find().SetDebug(true),
		)
		require.NoError(t, err)
		var emojis []string
		err = result.ScanSlice(&emojis)
		require.NoError(t, err)
		require.ElementsMatch(t, []string{
			`🤕`,
			`😀 😁 😂 🤣 😃 😄 😅 😆 😉 😊`,
			`😀 😁 😂 🤣 😃 😄 😅 😆 😉 😊`,
			`😀 😁 😂 🤣 😃 😄 😅 😆 😉 😊`,
			`🤕`,
		}, emojis)
	}

	{
		ns = normalStruct{}
		result, err = table.Find(
			actions.Find().
				Where(
					expr.In("$Key", actions.Find().
						Select("$Key").
						From("sqlike", "NormalStruct").
						Where(
							expr.Between("Tinyint", 1, 100),
						).
						OrderBy(
							expr.Desc("Timestamp"),
						),
					),
				).
				OrderBy(
					expr.Field("Enum", []Enum{
						Success,
						Failed,
						Unknown,
					}),
				),
			options.Find().SetDebug(true),
		)
		require.NoError(t, err)
		nss := []normalStruct{}

		err = result.All(&nss)
		require.NoError(t, err)
	}

	// query with Like
	{
		symbol := "Hal%o%()#$\\%^&_"
		ns = normalStruct{}
		err = table.FindOne(
			actions.FindOne().
				Where(
					expr.Like("FullText", symbol+"%"),
				),
			options.FindOne().SetDebug(true),
		).Decode(&ns)
		require.NoError(t, err)
		require.Equal(t, symbol, ns.FullText)
	}

	// Aggregation
	{
		ns = normalStruct{}
		result, err = table.Find(
			actions.Find().
				Select(
					expr.As("Enum", "A"),
					expr.As(expr.Count("$Key"), "B"),
					expr.Average("MediumInt"),
					expr.As(expr.Sum("SmallInt"), "C"),
					expr.Max("BigInt"),
					expr.As(expr.Min("BigInt"), "D"),
				).
				GroupBy(
					"Enum",
					"$Key",
				).
				OrderBy(
					expr.Desc("$Key"),
				),
			options.Find().
				SetDebug(true).
				SetNoLimit(true),
		)
		require.NoError(t, err)
		require.ElementsMatch(t,
			[]string{
				"A", "B", "AVG(`MediumInt`)",
				"C", "MAX(`BigInt`)", "D",
			}, result.Columns())
	}
}

// FindErrorExamples :
func FindErrorExamples(t *testing.T, db *sqlike.Database) {
	var (
		err error
	)

	{
		_, err = db.Table("unknown_table").Find(nil, options.Find().SetDebug(true))
		require.Error(t, err)

		err = db.Table("NormalStruct").
			FindOne(nil, options.FindOne().
				SetDebug(true)).Decode(nil)
		require.Error(t, err)
	}
}
