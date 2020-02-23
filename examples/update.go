package examples

import (
	"context"
	"database/sql"
	"testing"
	"time"

	uuid "github.com/google/uuid"
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/stretchr/testify/require"
)

// UpdateExamples :
func UpdateExamples(t *testing.T, db *sqlike.Database) {
	var (
		ns       normalStruct
		err      error
		result   sql.Result
		affected int64
		ctx      = context.Background()
	)

	table := db.Table("NormalStruct")
	uid, _ := uuid.Parse(`be72fc34-917b-11e9-af91-6c96cfd87b17`)

	{
		ns = normalStruct{}
		ns.ID = uid
		ns.Timestamp = time.Now()
		result, err = table.InsertOne(
			ctx,
			&ns,
			options.InsertOne().
				SetMode(options.InsertIgnore))
		affected, _ = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)
	}

	cnp := "✂️ Copy and 📋 Paste"
	now := time.Now().UTC()
	bigInt := int64(23127381237222)

	{
		// ModifyOne
		ns.Emoji = cnp
		ns.Enum = Success
		ns.BigInt = bigInt
		ns.DateTime = now
		ns.CreatedAt = now
		ns.UpdatedAt = now
		err = table.ModifyOne(
			ctx,
			&ns,
			options.ModifyOne().SetDebug(true))
		require.NoError(t, err)

		err = table.ModifyOne(
			ctx,
			&ns,
			options.ModifyOne().SetDebug(true))
		require.Error(t, err)

		err = table.ModifyOne(
			ctx,
			&ns,
			options.ModifyOne().
				SetStrict(false).
				SetDebug(true))
		require.NoError(t, err)

		ns2 := normalStruct{}
		err = table.FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("$Key", ns.ID),
				),
		).Decode(&ns2)
		require.NoError(t, err)
		require.Equal(t, cnp, ns2.Emoji)
		require.Equal(t, Success, ns2.Enum)
		require.Equal(t, bigInt, ns2.BigInt)
		require.NotZero(t, ns2.DateTime)
		require.NotZero(t, ns2.CreatedAt)
		require.NotZero(t, ns2.UpdatedAt)
	}

	{
		// ModifyOne with custom primary key
		type newStruct struct {
			Key     int64 `sqlike:"$Key"`
			No      int64
			Message string
			Flag    bool
			ID      int64 `sqlike:",primary_key"`
		}

		tbl := db.Table("NewStruct")

		tbl.MustMigrate(ctx, newStruct{})
		err = tbl.Truncate(ctx)
		require.NoError(t, err)

		ns := newStruct{}
		ns.Key = 8888
		ns.No = 1500
		ns.ID = 1
		result, err = tbl.InsertOne(
			ctx,
			&ns,
			options.InsertOne().
				SetMode(options.InsertIgnore))
		affected, _ = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)

		ns.Key = 6767
		ns.Message = "hello world"
		ns.Flag = true
		ns.No = 1800
		err = tbl.ModifyOne(
			ctx,
			&ns,
			options.ModifyOne().SetDebug(true),
		)
		require.NoError(t, err)

		ns2 := newStruct{}
		err = tbl.FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("ID", 1),
				),
		).Decode(&ns2)
		require.NoError(t, err)
		require.Equal(t, int64(6767), ns2.Key)
		require.Equal(t, int64(1800), ns2.No)
		require.Equal(t, "hello world", ns2.Message)
		require.Equal(t, true, ns2.Flag)
	}

	// Update single record
	{
		affected, err = table.UpdateOne(
			ctx,
			actions.UpdateOne().
				Where(expr.Equal("$Key", uid)).
				Set(
					expr.ColumnValue("LongStr", "1234abcd"),
					expr.ColumnValue("Emoji", "<😗>"),
				),
			options.UpdateOne().SetDebug(true),
		)

		require.NoError(t, err)
		require.Equal(t, int64(1), affected)
	}

	// Advance update query
	{
		affected, err = table.UpdateOne(
			ctx,
			actions.UpdateOne().
				Where(expr.Equal("$Key", uid)).
				Set(
					expr.ColumnValue("Emoji", "<😗>"),
					expr.ColumnValue("SID", expr.Column("LongStr")),
					expr.ColumnValue("Int", expr.Increment("Int", 100)),
					expr.ColumnValue("Tinyint", expr.Raw("80")),
				),
			options.UpdateOne().SetDebug(true),
		)

		require.NoError(t, err)
		require.Equal(t, int64(1), affected)
	}

}

// UpdateErrorExamples :
func UpdateErrorExamples(t *testing.T, db *sqlike.Database) {
	var (
		ns  normalStruct
		err error
		ctx = context.Background()
	)

	table := db.Table("NormalStruct")

	{
		err = table.ModifyOne(ctx, nil)
		require.Error(t, err)

		err = table.ModifyOne(ctx, &struct{}{})
		require.Error(t, err)

		err = table.ModifyOne(ctx, &ns)
		require.Equal(t, sqlike.ErrNoRecordAffected, err)
	}
}
