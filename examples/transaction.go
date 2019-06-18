package examples

import (
	"database/sql"
	"log"
	"strings"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/sql/expr"
	"github.com/stretchr/testify/require"
)

// TransactionExamples :
func TransactionExamples(t *testing.T, db *sqlike.Database) {
	var (
		uid      uuid.UUID
		ns       normalStruct
		result   sql.Result
		affected int64
		err      error
		tx       *sqlike.Transaction
	)

	log.Println("BeginTransaction ", strings.Repeat("=", 20))

	// Commit Transaction
	{
		uid, _ = uuid.FromString(`be72fc34-917b-11e9-af91-6c96cfd87a51`)

		ns = normalStruct{}
		ns.ID = uid
		ns.Timestamp = time.Now()
		tx, err = db.BeginTransaction()
		require.NoError(t, err)
		result, err = tx.Table("NormalStruct").InsertOne(&ns)
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)

		err = tx.CommitTransaction()
		require.NoError(t, err)
	}

	// Abort Transaction
	{
		uid, _ = uuid.FromString(`be7191c8-917b-11e9-af91-6c96cfd87a51`)

		ns = normalStruct{}
		ns.ID = uid
		ns.Timestamp = time.Now()
		tx, err = db.BeginTransaction()
		log.Println("Error :", err)
		require.NoError(t, err)
		result, err = tx.Table("NormalStruct").InsertOne(&ns)
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)

		err = tx.RollbackTransaction()
		require.NoError(t, err)

		ns = normalStruct{}
		err = db.Table("NormalStruct").FindOne(
			actions.FindOne().Where(
				expr.Equal("$Key", uid),
			),
		).Decode(&ns)
		require.Equal(t, sql.ErrNoRows, err)
	}
}
