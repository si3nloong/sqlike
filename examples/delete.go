package examples

import (
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
func DeleteExamples(t *testing.T, db *sqlike.Database) {
	var (
		affected int64
		result   sql.Result
		ns       normalStruct
		err      error
	)

	table := db.Table("NormalStruct")

	{
		err = table.FindOne(
			actions.FindOne().
				OrderBy(expr.Desc("$Key")),
		).Decode(&ns)
		require.NoError(t, err)
		err = table.DestroyOne(&ns)
		require.NoError(t, err)
	}

	// Single delete
	{
		ns := newNormalStruct()
		result, err = table.InsertOne(&ns,
			options.InsertOne().
				SetDebug(true),
		)
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)
		affected, err = table.DeleteOne(
			actions.DeleteOne().
				Where(
					expr.Equal("$Key", ns.ID),
				), options.DeleteOne().
				SetDebug(true),
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
		result, err = table.Insert(&nss,
			options.Insert().
				SetDebug(true),
		)
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(3), affected)
		affected, err = table.Delete(
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
