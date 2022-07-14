package examples

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/si3nloong/sqlike/v2"
	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/options"
	"github.com/si3nloong/sqlike/v2/sql/expr"
	"github.com/stretchr/testify/require"
)

type status string

const (
	statusActive  status = "ACTIVE"
	statusSuspend status = "SUSPEND"
)

// User :
type User struct {
	ID        int64 `sqlike:",auto_increment"`
	Name      string
	Age       int
	Status    status `sqlike:",enum=ACTIVE|SUSPEND"`
	CreatedAt time.Time
}

// Users :
type Users []User

// Len is part of sort.Interface.
func (usrs Users) Len() int {
	return len(usrs)
}

// Swap is part of sort.Interface.
func (usrs Users) Swap(i, j int) {
	usrs[i], usrs[j] = usrs[j], usrs[i]
}

type UserAddress struct {
	ID     int64 `sqlike:",auto_increment"`
	UserID int64 `sqlike:",foreign_key=User:ID"`
}

// PaginationExamples :
func PaginationExamples(ctx context.Context, t *testing.T, c *sqlike.Client) {
	var (
		// result *sqlike.Result
		err error
	)

	db := c.SetPrimaryKey("ID").Database("sqlike")
	table := db.Table("User")

	{
		err = db.Table("UserAddress").DropIfExists(ctx)
		require.NoError(t, err)

		err = table.DropIfExists(ctx)
		require.NoError(t, err)
	}

	{
		err = table.Migrate(ctx, User{})
		require.NoError(t, err)
	}

	{
		err = db.Table("UserAddress").Migrate(ctx, UserAddress{})
		require.NoError(t, err)
	}

	data := []User{
		{1, "User A", 18, statusActive, time.Now().UTC()},
		{2, "User B", 12, statusActive, time.Now().UTC()},
		{3, "User F", 20, statusActive, time.Now().UTC()},
		{4, "User C", 16, statusSuspend, time.Now().UTC()},
		{5, "User C", 16, statusActive, time.Now().UTC()},
		{6, "User G", 10, statusSuspend, time.Now().UTC()},
		{7, "User C", 16, statusActive, time.Now().UTC()},
		{8, "User D", 23, statusActive, time.Now().UTC()},
		{9, "User E", 30, statusSuspend, time.Now().UTC()},
	}

	{
		_, err = table.Insert(
			ctx,
			data, options.Insert().SetDebug(true),
		)
		require.NoError(t, err)
	}

	var (
		users  []User
		cursor any
	)

	sort.SliceStable(data, func(i, j int) bool {
		if data[i].Age > data[j].Age {
			return true
		}
		if data[i].Age < data[j].Age {
			return false
		}
		return data[i].ID > data[j].ID
	})

	// Paginate with simple query
	{
		pg, err := table.Paginate(
			ctx,
			actions.Paginate().
				Where(
					expr.GreaterOrEqual("Age", 0),
				).
				OrderBy(
					expr.Desc("Age"),
				).
				Limit(2),
			options.Paginate().
				SetDebug(true))
		require.NoError(t, err)

		for i := 0; i < len(data); i++ {
			users = []User{}
			err = pg.All(&users)
			require.NoError(t, err)
			if len(users) == 0 {
				break
			}
			require.Equal(t, data[i].ID, users[0].ID)
			cursor = users[len(users)-1].ID
			if pg.Cursor(ctx, cursor) != nil {
				break
			}
		}
	}

	// Paginate with complex query
	{
		users = []User{} // reset
		var result *sqlike.Result
		result, err := table.Find(
			ctx,
			actions.Find().
				Where(
					expr.GreaterOrEqual("Age", 16),
				).
				OrderBy(
					expr.Desc("Age"),
					expr.Desc("ID"),
				).
				Limit(100),
			options.Find().
				SetDebug(true))
		require.NoError(t, err)
		err = result.All(&users)
		require.NoError(t, err)

		results := []User{}
		limit := 2
		pg, err := table.Paginate(
			ctx,
			actions.Paginate().
				Where(
					expr.GreaterOrEqual("Age", 16),
				).
				OrderBy(
					expr.Desc("Age"),
				).
				Limit(uint(limit)),
			options.Paginate().SetDebug(true),
		)
		require.NoError(t, err)

		var (
			cursor int64
			i      int
		)

		for {
			err = pg.All(&results)
			if err != nil {
				require.NoError(t, err)
			}

			if len(results) == 0 || len(results) < limit {
				break
			}

			cursor = results[len(results)-1].ID
			require.True(t, len(users) > i)

			require.Equal(t, results[0], users[i])
			if err := pg.Cursor(ctx, cursor); err != nil {
				require.NoError(t, err)
			}

			i++
		}
	}

	length := 4
	actuals := [][]User{
		data[:length],
		data[length:(length * 2)],
		data[length*2:],
	}

	pg, err := table.Paginate(
		ctx,
		actions.Paginate().
			OrderBy(
				expr.Desc("Age"),
			).
			Limit(uint(length)),
		options.Paginate().
			SetDebug(true))
	require.NoError(t, err)

	// Expected paginate with error
	{
		err = pg.Cursor(ctx, nil)
		require.Error(t, err)
		err = pg.Cursor(ctx, []string{})
		require.Error(t, err)
		var nilslice []string
		err = pg.Cursor(ctx, nilslice)
		require.Error(t, err)
		var nilmap map[string]any
		err = pg.Cursor(ctx, nilmap)
		require.Error(t, err)
		err = pg.Cursor(ctx, "")
		require.Error(t, err)
		err = pg.Cursor(ctx, 0)
		require.Error(t, err)
		err = pg.Cursor(ctx, false)
		require.Error(t, err)
		err = pg.Cursor(ctx, float64(0))
		require.Error(t, err)
		err = pg.Cursor(ctx, []byte(nil))
		require.Error(t, err)
	}

	// pagination required more than 1 record
	{
		pg, err := table.Paginate(
			ctx,
			actions.Paginate().
				OrderBy(
					expr.Desc("Age"),
				).
				Limit(1),
			options.Paginate().
				SetDebug(true))
		require.Error(t, err)
		require.Nil(t, pg)
	}

	// Loop and get result set
	{
		users = []User{} // reset
		for i := 0; i < len(actuals); i++ {
			users = []User{}
			err = pg.All(&users)
			require.NoError(t, err)
			if len(users) == 0 {
				break
			}
			// require.ElementsMatch(t, actuals[i], users)
			cursor = users[len(users)-1].ID
			if pg.Cursor(ctx, cursor) != nil {
				break
			}
		}
	}
}
