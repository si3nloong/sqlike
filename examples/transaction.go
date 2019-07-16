package examples

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
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

	// RunInTransaction
	{
		err = db.RunInTransaction(func(ctx sqlike.SessionContext) error {
			uid, _ = uuid.FromString(`4ab3898c-9192-11e9-b500-6c96cfd87a51`)

			ns = normalStruct{}
			ns.ID = uid
			ns.Timestamp = time.Now()
			result, err := ctx.Table("NormalStruct").InsertOne(&ns)
			if err != nil {
				return err
			}

			affected, err := result.RowsAffected()
			if err != nil {
				return err
			}
			if affected < 1 {
				return errors.New("no result affected")
			}
			return nil
		})
		require.NoError(t, err)
	}

	{
		err = db.RunInTransaction(func(sessCtx sqlike.SessionContext) error {
			nss := []normalStruct{}
			result, err := sessCtx.Table("NormalStruct").
				Find(nil, options.LockForUpdate,
					options.Find().SetDebug(true))
			if err != nil {
				return err
			}
			result.All(&nss)
			time.Sleep(1 * time.Second)
			return nil
		})
		require.NoError(t, err)
	}
}