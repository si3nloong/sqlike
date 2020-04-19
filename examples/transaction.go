package examples

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"testing"
	"time"

	uuid "github.com/google/uuid"

	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/stretchr/testify/require"
)

type user struct {
	ID   int `sqlike:"$Key"`
	Name string
}

// TransactionExamples :
func TransactionExamples(t *testing.T, ctx context.Context, db *sqlike.Database) {
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
		uid, _ = uuid.Parse(`be72fc34-917b-11e9-af91-6c96cfd87a51`)
		now := time.Now()

		ns = normalStruct{}
		ns.ID = uid
		ns.DateTime = now
		ns.Timestamp = now
		ns.CreatedAt = now
		ns.UpdatedAt = now
		tx, err = db.BeginTransaction(ctx)
		require.NoError(t, err)
		result, err = tx.Table("NormalStruct").InsertOne(ctx, &ns)
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)

		err = tx.CommitTransaction()
		require.NoError(t, err)
	}

	// Abort Transaction
	{
		uid, _ = uuid.Parse(`be7191c8-917b-11e9-af91-6c96cfd87a51`)
		now := time.Now()

		ns = normalStruct{}
		ns.ID = uid
		ns.DateTime = now
		ns.Timestamp = now
		ns.CreatedAt = now
		ns.UpdatedAt = now
		tx, err = db.BeginTransaction(ctx)
		require.NoError(t, err)
		result, err = tx.Table("NormalStruct").InsertOne(ctx, &ns)
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)

		err = tx.RollbackTransaction()
		require.NoError(t, err)

		ns = normalStruct{}
		err = db.Table("NormalStruct").FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("$Key", uid),
				),
		).Decode(&ns)
		require.Equal(t, sql.ErrNoRows, err)
	}

	// RunInTransaction
	{
		err = db.RunInTransaction(ctx,
			func(sess sqlike.SessionContext) error {
				uid, _ = uuid.Parse(`4ab3898c-9192-11e9-b500-6c96cfd87a51`)
				now := time.Now()

				ns = normalStruct{}
				ns.ID = uid
				ns.DateTime = now
				ns.Timestamp = now
				ns.CreatedAt = now
				ns.UpdatedAt = now
				result, err := sess.Table("NormalStruct").InsertOne(sess, &ns)
				if err != nil {
					return err
				}

				ns.Int = 888
				if _, err := sess.Table("NormalStruct").
					UpdateOne(
						sess,
						actions.UpdateOne().
							Where(
								expr.Equal("$Key", ns.ID),
							).
							Set(
								expr.ColumnValue("Int", ns.Int),
							),
					); err != nil {
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

	// Timeout transaction
	{
		uid, _ = uuid.Parse(`5eb3f5c6-bfdb-11e9-88c7-6c96cfd87a51`)
		now := time.Now()
		err = db.RunInTransaction(
			ctx, func(sess sqlike.SessionContext) error {
				ns = normalStruct{}
				ns.ID = uid
				ns.DateTime = now
				ns.Timestamp = now
				ns.CreatedAt = now
				ns.UpdatedAt = now
				_, err := sess.Table("NormalStruct").
					InsertOne(
						sess,
						&ns, options.InsertOne().SetDebug(true),
					)
				if err != nil {
					return err
				}
				time.Sleep(5 * time.Second)
				return nil
			}, options.Transaction().SetTimeOut(3*time.Second))
		require.Equal(t, sql.ErrTxDone, err)

		rslt := normalStruct{}
		err = db.Table("NormalStruct").FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("$Key", uid),
				),
			options.FindOne().SetDebug(true),
		).Decode(&rslt)
		require.Error(t, err)
		require.Equal(t, normalStruct{}, rslt)
	}

	// Lock record using transaction
	{
		err = db.RunInTransaction(
			ctx, func(sess sqlike.SessionContext) error {
				nss := []normalStruct{}
				result, err := sess.Table("NormalStruct").
					Find(
						sess,
						nil, options.Find().
							SetLockMode(options.LockForUpdate).
							SetDebug(true),
					)
				if err != nil {
					return err
				}
				if err := result.All(&nss); err != nil {
					return err
				}
				time.Sleep(1 * time.Second)
				return nil
			})
		require.NoError(t, err)
	}

	table := db.Table("user")
	err = table.DropIfExists(ctx)
	require.NoError(t, err)
	table.MustMigrate(ctx, new(user))

	// Commit Transaction
	{
		data := &user{ID: rand.Intn(10000), Name: "Oska"}
		trx, _ := db.BeginTransaction(ctx)
		_, err = trx.Table("user").InsertOne(ctx, data)
		require.NoError(t, err)

		err = trx.Table("user").FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("$Key", data.ID),
				),
		).Decode(&user{})
		require.NoError(t, err)

		err = db.Table("user").FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("$Key", data.ID),
				),
		).Decode(&user{})
		require.Error(t, err)

		err = trx.CommitTransaction()
		require.NoError(t, err)

		err = db.Table("user").FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("$Key", data.ID),
				),
		).Decode(&user{})
		require.NoError(t, err)

		// Remove tested data
		err = db.Table("user").DestroyOne(ctx, data)
		require.NoError(t, err)
	}

	// Rollback Transaction
	{
		data := &user{ID: 1234, Name: "Oska"}
		trx, _ := db.BeginTransaction(ctx)
		_, err = trx.Table("user").InsertOne(ctx, data)
		require.NoError(t, err)

		err = trx.Table("user").FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("$Key", data.ID),
				),
		).Decode(&user{})
		require.NoError(t, err)

		err = db.Table("user").FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("$Key", data.ID),
				),
		).Decode(&user{})
		require.Error(t, err)

		err = trx.RollbackTransaction()
		require.NoError(t, err)

		err = db.Table("user").FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("$Key", data.ID),
				),
		).Decode(&user{})
		require.Error(t, err)
	}
}
