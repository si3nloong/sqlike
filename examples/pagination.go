package examples

import (
	"context"
	"sort"
	"testing"

	"github.com/si3nloong/sqlike/v2"
	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/options"
	"github.com/si3nloong/sqlike/v2/sql/expr"
	"github.com/stretchr/testify/require"
)

// PaginationExamples :
func PaginationExamples(ctx context.Context, t *testing.T, c *sqlike.Client) {
	var (
		// result *sqlike.Result
		err error
	)

	db := c.SetPrimaryKey("ID").Database("sqlike")
	table := db.Table("User")
	addressTable := db.Table("UserAddress")

	{
		err = addressTable.DropIfExists(ctx)
		require.NoError(t, err)

		err = table.DropIfExists(ctx)
		require.NoError(t, err)
	}

	{
		err = table.Migrate(ctx, User{})
		require.NoError(t, err)
	}

	{
		err = addressTable.Migrate(ctx, UserAddress{})
		require.NoError(t, err)
	}

	data := []User{
		{ID: 1, Name: "User A", Age: 18, Status: userStatusActive},
		{ID: 2, Name: "User B", Age: 12, Status: userStatusActive},
		{ID: 3, Name: "User F", Age: 20, Status: userStatusActive},
		{ID: 4, Name: "User C", Age: 16, Status: userStatusSuspend},
		{ID: 5, Name: "User C", Age: 16, Status: userStatusActive},
		{ID: 6, Name: "User G", Age: 10, Status: userStatusSuspend},
		{ID: 7, Name: "User C", Age: 16, Status: userStatusActive},
		{ID: 8, Name: "User D", Age: 23, Status: userStatusActive},
		{ID: 9, Name: "User E", Age: 30, Status: userStatusSuspend},
	}

	{
		_, err = table.Insert(
			ctx,
			data, options.Insert().SetDebug(true),
		)
		require.NoError(t, err)

		_, err = addressTable.Insert(
			ctx,
			&[]UserAddress{
				{UserID: data[0].ID},
				{UserID: data[3].ID},
				{UserID: data[2].ID},
				{UserID: data[6].ID},
			}, options.Insert().SetDebug(true),
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
			if pg.After(ctx, cursor) != nil {
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
			if err := pg.After(ctx, cursor); err != nil {
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
		err = pg.After(ctx, nil)
		require.Error(t, err)
		err = pg.After(ctx, []string{})
		require.Error(t, err)
		var nilslice []string
		err = pg.After(ctx, nilslice)
		require.Error(t, err)
		var nilmap map[string]any
		err = pg.After(ctx, nilmap)
		require.Error(t, err)
		err = pg.After(ctx, "")
		require.Error(t, err)
		err = pg.After(ctx, 0)
		require.Error(t, err)
		err = pg.After(ctx, false)
		require.Error(t, err)
		err = pg.After(ctx, float64(0))
		require.Error(t, err)
		err = pg.After(ctx, []byte(nil))
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
			if pg.After(ctx, cursor) != nil {
				break
			}
		}
	}
}
