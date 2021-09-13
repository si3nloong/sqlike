package examples

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/si3nloong/sqlike/types"
	"github.com/stretchr/testify/require"
)

// InsertExamples :
func InsertExamples(ctx context.Context, t *testing.T, db *sqlike.Database) {
	var (
		err      error
		result   sql.Result
		affected int64
	)

	table := db.Table("NormalStruct")
	ns := newNormalStruct()

	// Single insert
	{
		result, err = table.InsertOne(
			ctx,
			&ns,
			options.InsertOne().
				SetOmitFields("Int").
				SetDebug(true))
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)
	}

	// Single upsert
	// - https://dev.mysql.com/doc/refman/8.0/en/insert-on-duplicate.html
	{
		ns.Emoji = `🤕`
		m := make(map[string]int)
		m["one"] = 1
		m["two"] = 2
		ns.Map = m
		result, err = table.InsertOne(
			ctx,
			&ns,
			options.InsertOne().
				SetDebug(true).
				SetMode(options.InsertOnDuplicate))
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		// upsert affected is 2 instead of 1
		require.Equal(t, int64(2), affected)

	}

	// upsert with omitfield
	{
		ns2 := newNormalStruct()
		result, err = table.InsertOne(
			ctx,
			&ns2,
			options.InsertOne().
				SetDebug(true).
				SetMode(options.InsertOnDuplicate))
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)

		temp := ns2
		temp.Date = types.Date{Year: 2020, Month: 12, Day: 7}
		temp.BigUint = 188
		temp.BigInt = 188
		temp.Byte = []byte("testing 123")
		temp.Emoji = "🥶 😱 😨 😰"
		result, err = table.InsertOne(
			ctx,
			&temp,
			options.InsertOne().
				SetOmitFields("Date", "BigUint", "EmptyByte", "Byte").
				SetDebug(true).
				SetMode(options.InsertOnDuplicate))
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		// upsert affected is 2 instead of 1
		require.Equal(t, int64(2), affected)

		var o normalStruct
		err = table.FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("$Key", ns2.ID),
				),
			options.FindOne().SetDebug(true),
		).Decode(&o)
		require.NoError(t, err)

		// ensure the value didn't get modified on duplicate
		require.Equal(t, o.BigUint, ns2.BigUint)
		require.Equal(t, o.Byte, ns2.Byte)
		require.NotEqual(t, o.BigInt, ns2.BigInt)
		require.NotEqual(t, o.Emoji, ns2.Emoji)
		require.Equal(t, o.Date.String(), ns2.Date.String())
	}

	// upsert with unordered primary key
	{
		type unorderedStruct struct {
			Text   string
			Number int8
			ID     string `sqlike:",primary_key"`
			Date   civil.Date
		}

		table := db.Table("unorderedStruct")
		table.DropIfExists(ctx)
		table.MustUnsafeMigrate(ctx, unorderedStruct{})
		_, err = table.InsertOne(
			ctx,
			&unorderedStruct{
				ID:   "224",
				Date: civil.DateOf(time.Now()),
			},
			options.InsertOne().SetMode(options.InsertOnDuplicate),
		)
		require.NoError(t, err)
	}

	// Multiple insert
	{
		nss := [...]normalStruct{
			newNormalStruct(),
			newNormalStruct(),
			newNormalStruct(),
		}
		result, err = table.Insert(
			ctx,
			&nss,
			options.Insert().SetDebug(true),
		)
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(3), affected)
	}

	table2 := db.Table("PtrStruct")

	// Pointer insertion
	{
		_, err = table2.InsertOne(
			ctx,
			&ptrStruct{},
			options.InsertOne().SetDebug(true),
		)
		require.NoError(t, err)
	}

	// Pointer insertion
	{
		ps := []ptrStruct{
			newPtrStruct(),
			newPtrStruct(),
			newPtrStruct(),
			newPtrStruct(),
			newPtrStruct(),
		}
		_, err = table2.Insert(
			ctx,
			&ps,
			options.Insert().SetDebug(true),
		)
		require.NoError(t, err)
	}

	// Error insertion
	{
		_, err = table.InsertOne(
			ctx,
			&struct {
				Interface interface{}
			}{},
		)
		require.Error(t, err)
		_, err = table.InsertOne(ctx, struct{}{})
		require.Error(t, err)
		var empty *struct{}
		_, err = table.InsertOne(ctx, empty)
		require.Error(t, err)

		_, err = table.Insert(ctx, []interface{}{})
		require.Error(t, err)
	}

	table3 := db.Table("GeneratedStruct")

	// generated column insertion
	{
		cols := []*generatedStruct{
			newGeneratedStruct(),
			newGeneratedStruct(),
			newGeneratedStruct(),
			newGeneratedStruct(),
			newGeneratedStruct(),
		}
		_, err = table3.Insert(
			ctx,
			&cols,
			options.Insert().SetDebug(true),
		)
		require.NoError(t, err)
	}

	{
		table := db.Table("Override")
		err := table.Truncate(ctx)
		require.NoError(t, err)
		table.MustUnsafeMigrate(ctx, overrideStruct{})

		now, err := time.Parse("2006-01-02 15:04:05", "2018-01-01 00:00:00")
		require.NoError(t, err)

		os := overrideStruct{}
		os.Amount = 9000
		os.Amount = 100
		os.ID = 1007
		os.generatedStruct.ID = ""
		os.CivilDate = civil.DateOf(now)
		os.CreatedAt = now
		os.UpdatedAt = now
		_, err = table.InsertOne(ctx, &os)
		require.NoError(t, err)

		var out overrideStruct
		err = table.FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("ID", os.ID),
				),
		).Decode(&out)
		require.NoError(t, err)

		require.Equal(t, os.ID, out.ID)
		require.Equal(t, now, out.CreatedAt)
		require.Equal(t, now, out.UpdatedAt)
		require.Equal(t, os.Amount, out.Amount)
		require.Equal(t, os.CivilDate, out.CivilDate)
	}
}

// InsertErrorExamples :
func InsertErrorExamples(ctx context.Context, t *testing.T, db *sqlike.Database) {
	var (
		ns  normalStruct
		err error
	)

	{
		_, err = db.Table("NormalStruct").InsertOne(ctx, nil)
		require.Error(t, err)

		var uninitialized *normalStruct
		_, err = db.Table("NormalStruct").InsertOne(ctx, uninitialized)
		require.Error(t, err)

		ns = normalStruct{}
		_, err = db.Table("NormalStruct").InsertOne(ctx, ns)
		require.Error(t, err)
	}

	{
		_, err = db.Table("NormalStruct").Insert(
			ctx,
			[]normalStruct{},
		)
		require.Error(t, err)
	}
}
