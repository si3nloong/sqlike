package examples

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/si3nloong/sqlike/v2"
	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/options"
	"github.com/si3nloong/sqlike/v2/sql"
	"github.com/si3nloong/sqlike/v2/sql/expr"
	"github.com/si3nloong/sqlike/v2/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

// FindExamples :
func FindExamples(ctx context.Context, t *testing.T, db *sqlike.Database) {
	var (
		result *sqlike.Rows
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
	jsonRaw := json.RawMessage(`{"test":"hello world"}`)
	lang := language.Japanese
	langs := []language.Tag{
		language.AmericanEnglish,
		language.BrazilianPortuguese,
		language.Malay,
		language.ModernStandardArabic,
		language.Korean,
		language.Japanese,
	}
	virtualColumn := "virtual column"
	numMap := map[string]int{
		"one":    1,
		"three":  3,
		"eleven": 11,
	}

	table := db.Table("NormalStruct")

	// insert record before find
	{
		now := time.Now()
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
		ns.JSONRaw = jsonRaw
		ns.Enum = Failed
		ns.Map = make(map[string]int)
		ns.Map["one"] = 1
		ns.Map["three"] = 3
		ns.Map["eleven"] = 11
		ns.Struct.VirtualStr = virtualColumn
		ns.Date = civil.DateOf(now)
		ns.DateTime = ts
		ns.Timestamp = ts
		ns.Date = civil.DateOf(ts)
		ns.Language = lang
		ns.Languages = langs
		ns.Set = append(ns.Set, "A", "A", "B")
		ns.CreatedAt = now
		ns.UpdatedAt = now

		_, err = table.InsertOne(
			ctx, &ns,
		)
		require.NoError(t, err)
	}

	// Find one record by primary key
	{
		ns = normalStruct{}
		result := table.FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("$Key", uid),
				),
			options.FindOne().SetDebug(true),
		)

		err = result.Decode(&ns)
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
		require.Equal(t, virtualColumn, ns.Struct.VirtualStr)

		require.Nil(t, ns.Struct.NestedNullInt)
		require.Nil(t, ns.Struct.Key)
		require.Nil(t, ns.PtrUUID)
		var nilKey *types.Key
		require.Equal(t, nilKey, ns.Struct.Key)
		require.True(t, ns.Struct.Key == nil)
		require.True(t, ns.Struct.Key == nilKey)

		require.Equal(t, numMap, ns.Map)
		require.Equal(t, lang, ns.Language)
		require.Equal(t, langs, ns.Languages)
		require.ElementsMatch(t, types.Set{"A", "B"}, ns.Set)
		require.Equal(t, json.RawMessage(`{"test":"hello world"}`), ns.JSONRaw)

		columns := []string{
			"$Key", "Key", "PtrUUID", "Date",
			"SID", "Emoji", "FullText", "LongStr", "CustomStrType",
			"EmptyByte", "Byte", "Bool",
			"Int", "TinyInt", "SmallInt", "MediumInt", "BigInt",
			"Uint", "TinyUint", "SmallUint", "MediumUint", "BigUint",
			"Float32", "Float64", "UFloat32",
			"EmptyStruct", "Struct", "VirtualColumn",
			"Struct.StoredStr", "JSONRaw", "Map",
			"DateTime", "Timestamp", "Location", "Language", "Languages",
			"Currency", "Currencies",
			"Enum", "Set",
			"CreatedAt", "UpdatedAt",
		}
		cols := result.Columns()
		sort.Strings(columns)
		sort.Strings(cols)
		require.True(t, len(cols) > 0)
		require.ElementsMatch(t, columns, cols)
	}

	// Find one with scan
	{
		var i struct {
			skip      string
			count     uint
			id        *string
			emoji     string
			customStr string
			boolean   bool
			jsonRaw   json.RawMessage
			numMap    map[string]int
		}
		ns = normalStruct{}

		/*
			SELECT
				COUNT(`$Key`) AS `c`
			FROM `sqlike`.`NormalStruct`
			WHERE
				`$Key` = "e7977246-910a-11e9-844d-6c96cfd87a51"
			LIMIT 1;
		*/
		// Scan with unmatched number of fields
		err = table.FindOne(
			ctx,
			actions.FindOne().
				Select(
					expr.As(expr.Count("$Key"), "c"),
				).
				Where(
					expr.Equal("$Key", uid),
				),
			options.FindOne().SetDebug(true),
		).Scan(&i.count, i.skip, &i.id, &i.emoji)
		require.NoError(t, err)
		require.True(t, i.count > 0)

		/*
			SELECT
				`$Key`,`Emoji`,`CustomStrType`,`Bool`,`JSONRaw`,`Map`,`Language`
			FROM `sqlike`.`NormalStruct`
			WHERE
				`$Key` = "e7977246-910a-11e9-844d-6c96cfd87a51"
			LIMIT 1;
		*/
		// Scan with fields
		err = table.FindOne(
			ctx,
			actions.FindOne().
				Select(
					expr.Func("BIN_TO_UUID", expr.Column("$Key")),
					"Emoji", "CustomStrType", "Bool",
					"JSONRaw", "Map", "Language",
				).
				Where(
					expr.Equal("$Key", uid),
				),
			options.FindOne().SetDebug(true),
		).Scan(&i.id, &i.emoji, &i.customStr, &i.boolean, &i.jsonRaw, &i.numMap)
		require.NoError(t, err)
		require.NotNil(t, i.id)
		require.Equal(t, uid.String(), *i.id)
		require.Equal(t, emoji, i.emoji)
		require.Equal(t, jsonRaw, i.jsonRaw)
		require.Equal(t, numMap, i.numMap)

		// Scan error
		err = table.FindOne(
			ctx,
			actions.FindOne().
				Select(
					"$Key",
				).
				Where(
					expr.Equal("$Key", uid),
				),
			options.FindOne().SetDebug(true),
		).Scan(i.skip)
		require.Error(t, err)
	}

	// Find one record by primary key
	{
		/*
			SELECT * FROM `sqlike`.`NormalStruct`
			WHERE
				`$Key` = "1000"
			LIMIT 1;
		*/
		ns = normalStruct{}
		err = table.FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("$Key", "1000"),
				),
		).Decode(&ns)
		require.Equal(t, err, sqlike.ErrNoRows)

	}

	// Find multiple records by where condition
	{
		/*
			SELECT * FROM `sqlike`.`NormalStruct`
			WHERE
				(
					`TinyInt` BETWEEN 1 AND 100 AND
					`Enum` IN ("SUCCESS","FAILED","UNKNOWN")
				)
			LIMIT 100;
		*/
		ns = normalStruct{}
		nss := []normalStruct{}
		result, err = table.Find(
			ctx,
			actions.Find().
				Where(
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
		err = result.All(&nss)
		require.NoError(t, err)

	}

	// Find with scan slice
	{
		/*
			SELECT `Emoji` FROM `sqlike`.`NormalStruct` LIMIT 100;
		*/
		ns = normalStruct{}
		result, err = table.Find(
			ctx,
			actions.Find().Select("Emoji"),
		)
		require.NoError(t, err)
		var emojis []string
		err = result.ScanSlice(&emojis)
		require.NoError(t, err)
		require.ElementsMatch(t, []string{
			`🤕`,
			`🥶 😱 😨 😰`,
			`😀 😁 😂 🤣 😃 😄 😅 😆 😉 😊`,
			`😀 😁 😂 🤣 😃 😄 😅 😆 😉 😊`,
			`😀 😁 😂 🤣 😃 😄 😅 😆 😉 😊`,
			`🤕`,
		}, emojis)
	}

	// Find with subquery
	{
		/*
			SELECT * FROM `sqlike`.`NormalStruct`
			WHERE (
				`$Key` IN (
					SELECT `$Key` FROM `sqlike`.`NormalStruct`
					WHERE `Tinyint` BETWEEN 1 AND 100
					ORDER BY `Timestamp` DESC
				) AND
				EXISTS (SELECT 1 FROM `sqlike`.`NormalStruct`)
			)
			ORDER BY FIELD(`Enum`,"SUCCESS","FAILED","UNKNOWN")
			LIMIT 100;
		*/
		ns = normalStruct{}
		result, err = table.Find(
			ctx,
			actions.Find().
				Where(
					expr.In("$Key", sql.Select("$Key").
						From("sqlike", "NormalStruct").
						Where(
							expr.Between("Tinyint", 1, 100),
						).
						OrderBy(
							expr.Desc("Timestamp"),
						),
					),
					expr.Exists(
						sql.Select(expr.Raw("1")).
							From("sqlike", "NormalStruct"),
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

	// Query with Like expression
	{
		/*
			SELECT * FROM `sqlike`.`NormalStruct`
			WHERE `FullText` LIKE "Hal\\%o\\%()#$\\\\\\%^&\\_%"
			LIMIT 1;
		*/
		symbol := "Hal%o%()#$\\%^&_"
		ns = normalStruct{}
		err = table.FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Or(
						expr.Like("FullText", symbol+"%"),
						expr.Like("FullText", expr.Raw(strconv.Quote("%"+symbol+"%"))),
					),
				),
			options.FindOne().SetDebug(true),
		).Decode(&ns)
		require.NoError(t, err)
		require.Equal(t, symbol, ns.FullText)
	}

	// Aggregation
	{
		/*
			SELECT
				`Enum` AS `A`,
				COUNT(`$Key`) AS `B`,
				AVG(`MediumInt`),
				COALESCE(SUM(`SmallInt`),0) AS `C`,
				MAX(`BigInt`),
				MIN(`BigInt`) AS `D`
			FROM `sqlike`.`NormalStruct`
			GROUP BY
				`Enum`,
				`$Key`
			ORDER BY `$Key` DESC;
		*/
		ns = normalStruct{}
		result, err = table.Find(
			ctx,
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
			}, result.Columns(),
		)
	}

	{
		table := db.Table("GeneratedStruct")

		first := newGeneratedStruct()
		cols := []*generatedStruct{
			first,
			newGeneratedStruct(),
		}
		_, err = table.Insert(
			ctx,
			&cols,
			options.Insert().SetDebug(true),
		)
		require.NoError(t, err)
		require.Empty(t, first.ID)

		var result generatedStruct

		err = table.FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("NestedID", first.Nested.ID),
				),
			options.FindOne().SetDebug(true),
		).Decode(&result)
		require.NoError(t, err)
		require.Equal(t, first.Nested.ID, result.ID)
		require.True(t, result.Amount > 0)
	}

}

// FindErrorExamples :
func FindErrorExamples(ctx context.Context, t *testing.T, db *sqlike.Database) {
	var (
		err error
	)

	{
		_, err = db.Table("unknown_table").Find(
			ctx,
			nil,
			options.Find().SetDebug(true),
		)
		require.Error(t, err)

		err = db.Table("NormalStruct").FindOne(
			ctx,
			nil,
			options.FindOne().SetDebug(true),
		).Decode(nil)
		require.Error(t, err)
	}
}
