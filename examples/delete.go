package examples

import (
	"context"
	"database/sql"
	"testing"

	uuid "github.com/google/uuid"
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/stretchr/testify/require"
)

// DeleteExamples :
func DeleteExamples(ctx context.Context, t *testing.T, db *sqlike.Database) {
	var (
		affected int64
		result   sql.Result
		ns       normalStruct
		err      error
	)

	table := db.Table("NormalStruct")

	// Delete record with primary key
	{
		err = table.FindOne(
			ctx,
			actions.FindOne().
				OrderBy(expr.Desc("$Key")),
		).Decode(&ns)
		require.NoError(t, err)
		err = table.DestroyOne(ctx, &ns)
		require.NoError(t, err)
	}

	// Delete record with primary key tag (DestroyOne)
	{
		type dummy struct {
			UUID uuid.UUID `sqlike:",primary_key"`
		}

		table := db.Table("testDB")
		err = table.DropIfExists(ctx)
		require.NoError(t, err)

		table.MustUnsafeMigrate(ctx, dummy{})
		records := []dummy{
			{UUID: uuid.New()},
			{UUID: uuid.New()},
		}

		_, err = table.Insert(ctx, &records)
		require.NoError(t, err)

		// destroy with empty value should error
		{
			var nilDummy *dummy
			err = table.DestroyOne(ctx, nilDummy)
			require.Error(t, err)

			err = table.DestroyOne(ctx, nil)
			require.Error(t, err)
		}

		err = table.DestroyOne(
			ctx,
			records[0],
			options.DestroyOne().SetDebug(true),
		)
		require.NoError(t, err)

		err = table.DestroyOne(
			ctx,
			&records[1],
			options.DestroyOne().SetDebug(true),
		)
		require.NoError(t, err)

		var count uint
		if err := table.FindOne(
			ctx,
			actions.FindOne().Select(expr.Count("UUID")),
			options.FindOne().SetDebug(true),
		).Scan(&count); err != nil {
			require.NoError(t, err)
		}
		require.Equal(t, uint(0), count)
	}

	// Single delete
	{
		ns := newNormalStruct()
		result, err = table.InsertOne(
			ctx,
			&ns,
			options.InsertOne().SetDebug(true),
		)
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)
		affected, err = table.DeleteOne(
			ctx,
			actions.DeleteOne().
				Where(
					expr.Equal("$Key", ns.ID),
				),
			options.DeleteOne().SetDebug(true),
		)
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)
	}

	// Multiple delete
	{
		nss := [...]normalStruct{
			newNormalStruct(),
			newNormalStruct(),
			newNormalStruct(),
		}
		result, err = table.Insert(
			ctx,
			&nss,
			options.Insert().
				SetDebug(true),
		)
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(3), affected)
		affected, err = table.Delete(
			ctx,
			actions.Delete().
				Where(
					expr.In("$Key", []uuid.UUID{
						nss[0].ID,
						nss[1].ID,
						nss[2].ID,
					}),
				), options.Delete().
				SetDebug(true),
		)
		require.NoError(t, err)
		require.Equal(t, int64(3), affected)
	}
}
