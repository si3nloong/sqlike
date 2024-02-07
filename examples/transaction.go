package examples

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"

	"github.com/si3nloong/sqlike/v2"
	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/options"
	"github.com/si3nloong/sqlike/v2/sql/expr"
	"github.com/stretchr/testify/require"
)

type user struct {
	ID   int `sqlike:"$Key"`
	Name string
}

// TransactionExamples :
func TransactionExamples(ctx context.Context, t *testing.T, db *sqlike.Database) {
	var (
		uid      uuid.UUID
		ns       normalStruct
		result   sql.Result
		affected int64
		err      error
		tx       *sqlike.Transaction
	)

	// Commit Transaction
	t.Run("Commit Transaction", func(t *testing.T) {
		uid, _ = uuid.Parse(`be72fc34-917b-11e9-af91-6c96cfd87a51`)
		now := time.Now()

		ns = normalStruct{}
		ns.ID = uid
		ns.Date = civil.DateOf(now)
		ns.DateTime = now
		ns.Timestamp = now
		ns.CreatedAt = now
		ns.UpdatedAt = now
		tx, err = db.BeginTransaction(ctx)
		require.NoError(t, err)
		result, err = db.Table("NormalStruct").InsertOne(tx, &ns, options.InsertOne().SetDebug(true))
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)

		err = tx.Commit()
		require.NoError(t, err)
	})

	// Abort Transaction
	t.Run("Abort Transaction", func(t *testing.T) {
		uid, _ = uuid.Parse(`be7191c8-917b-11e9-af91-6c96cfd87a51`)
		now := time.Now()

		ns = normalStruct{}
		ns.ID = uid
		ns.Date = civil.DateOf(now)
		ns.DateTime = now
		ns.Timestamp = now
		ns.CreatedAt = now
		ns.UpdatedAt = now
		tx, err = db.BeginTransaction(ctx)
		require.NoError(t, err)
		result, err = db.Table("NormalStruct").InsertOne(tx, &ns)
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)

		err = tx.Rollback()
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
	})

	// RunInTransaction
	t.Run("RunInTransaction", func(t *testing.T) {
		err = db.RunInTransaction(ctx, func(sess context.Context) error {
			table := db.Table("NormalStruct")
			uid, _ = uuid.Parse(`4ab3898c-9192-11e9-b500-6c96cfd87a51`)
			now := time.Now()

			ns = normalStruct{}
			ns.ID = uid
			ns.Date = civil.DateOf(now)
			ns.DateTime = now
			ns.Timestamp = now
			ns.CreatedAt = now
			ns.UpdatedAt = now

			result, err := table.InsertOne(sess, &ns)
			if err != nil {
				return err
			}

			ns.Int = 888
			if _, err := table.UpdateOne(
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
	})

	// Timeout transaction
	t.Run("RunInTransaction with timeout", func(t *testing.T) {
		uid, _ = uuid.Parse(`5eb3f5c6-bfdb-11e9-88c7-6c96cfd87a51`)
		now := time.Now()
		err = db.RunInTransaction(ctx, func(sess context.Context) error {
			ns = normalStruct{}
			ns.ID = uid
			ns.Date = civil.DateOf(now)
			ns.DateTime = now
			ns.Timestamp = now
			ns.CreatedAt = now
			ns.UpdatedAt = now
			_, err := db.Table("NormalStruct").
				InsertOne(
					sess,
					&ns,
					options.InsertOne().SetDebug(true),
				)
			if err != nil {
				return err
			}
			time.Sleep(3 * time.Second)
			return nil
		}, options.Transaction().SetTimeOut(2*time.Second))
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
		require.True(t, errors.Is(err, sql.ErrNoRows))
		require.Equal(t, normalStruct{}, rslt)
	})

	// Lock record using transaction
	t.Run("RunInTransaction with lock", func(t *testing.T) {
		err = db.RunInTransaction(
			ctx, func(sess context.Context) error {
				nss := []normalStruct{}
				result, err := db.Table("NormalStruct").
					Find(
						sess,
						nil, options.Find().
							SetLockMode(options.LockForUpdate()).
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
	})

	err = db.Table("UserAddress").DropIfExists(ctx)
	require.NoError(t, err)

	table := db.Table("user")
	err = db.Table("UserAddress").DropIfExists(ctx)
	require.NoError(t, err)
	err = table.DropIfExists(ctx)
	require.NoError(t, err)
	table.MustMigrate(ctx, new(user))

	// Commit Transaction
	t.Run("BeginTransaction with different context", func(t *testing.T) {
		data := &user{ID: rand.Intn(10000), Name: "Oska"}
		tx, err := db.BeginTransaction(ctx)
		require.NoError(t, err)
		defer tx.Rollback()

		_, err = db.Table("user").InsertOne(tx, data)
		require.NoError(t, err)

		err = db.Table("user").FindOne(
			tx,
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

		err = tx.Commit()
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
	})

	// Rollback Transaction
	t.Run("Rollback Transaction", func(t *testing.T) {
		data := &user{ID: 1234, Name: "Oska"}
		tx, _ := db.BeginTransaction(ctx)
		_, err = db.Table("user").InsertOne(tx, data)
		require.NoError(t, err)

		err = db.Table("user").FindOne(
			tx,
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

		err = tx.Rollback()
		require.NoError(t, err)

		err = db.Table("user").FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("$Key", data.ID),
				),
		).Decode(&user{})
		require.Error(t, err)
	})

	// Nested transaction
	t.Run("BeginTransaction with nested transactions", func(t *testing.T) {
		// valid transaction
		tx1, err := db.BeginTransaction(context.TODO())
		require.NoError(t, err)
		require.NotNil(t, tx1)

		// shouldn't allow beginTransaction with nested transaction
		tx2, err := db.BeginTransaction(tx1)
		require.Error(t, err)
		require.Nil(t, tx2)

		// no matter how many level the context wrap with other context, it will still consider as nested transaction
		ctx := context.WithValue(tx1, "key", "lv2")
		ctx = context.WithValue(ctx, "key", "lv1")
		tx2, err = db.BeginTransaction(ctx)
		require.Error(t, err)
		require.Nil(t, tx2)
	})
}
