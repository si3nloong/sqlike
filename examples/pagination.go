package examples

import (
	"sort"
	"testing"
	"time"

	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/stretchr/testify/require"
)

type Status string

const (
	StatusActive  = "ACTIVE"
	StatusSuspend = "SUSPEND"
)

// User :
type User struct {
	ID        int64
	Name      string
	Age       int
	Status    Status `sqlike:",enum:ACTIVE|SUSPEND"`
	CreatedAt time.Time
}

type Users []User

// Len is part of sort.Interface.
func (usrs Users) Len() int {
	return len(usrs)
}

// Swap is part of sort.Interface.
func (usrs Users) Swap(i, j int) {
	usrs[i], usrs[j] = usrs[j], usrs[i]
}

// PaginationExamples :
func PaginationExamples(t *testing.T, c *sqlike.Client) {
	var (
		// result *sqlike.Result
		err error
	)

	db := c.SetPrimaryKey("ID").Database("sqlike")
	table := db.Table("User")

	{
		err = table.DropIfExits()
		require.NoError(t, err)
	}

	{
		err = table.Migrate(User{})
		require.NoError(t, err)
	}

	data := []User{
		User{10, "User A", 18, StatusActive, time.Now().UTC()},
		User{88, "User B", 12, StatusActive, time.Now().UTC()},
		User{8, "User F", 20, StatusActive, time.Now().UTC()},
		User{27, "User C", 16, StatusSuspend, time.Now().UTC()},
		User{20, "User C", 16, StatusActive, time.Now().UTC()},
		User{100, "User G", 10, StatusSuspend, time.Now().UTC()},
		User{21, "User C", 16, StatusActive, time.Now().UTC()},
		User{50, "User D", 23, StatusActive, time.Now().UTC()},
		User{5, "User E", 30, StatusSuspend, time.Now().UTC()},
	}

	{
		_, err = table.InsertMany(data, options.
			InsertMany().SetDebug(true))
		require.NoError(t, err)
	}

	var (
		users  []User
		cursor interface{}
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

	{
		pg, err := table.Paginate(actions.Paginate().
			Where().
			OrderBy(
				expr.Desc("Age"),
			).Limit(1),
			options.Paginate().
				SetDebug(true))
		require.NoError(t, err)

		for i := 0; i < len(data); i++ {
			if pg.NextPage(cursor) != nil {
				break
			}
			users = []User{}
			err = pg.All(&users)
			require.NoError(t, err)
			if len(users) == 0 {
				break
			}

			require.Equal(t, data[i], users[0])
			cursor = users[len(users)-1].ID
		}
	}

	{
		actuals := [][]User{
			data[:5],
			data[5:],
		}

		cursor = nil
		pg, err := table.Paginate(actions.Paginate().
			Where().
			OrderBy(
				expr.Desc("Age"),
			).Limit(5),
			options.Paginate().
				SetDebug(true))
		require.NoError(t, err)

		for i := 0; i < len(actuals); i++ {
			if pg.NextPage(cursor) != nil {
				break
			}
			users = []User{}
			err = pg.All(&users)
			require.NoError(t, err)
			if len(users) == 0 {
				break
			}
			require.ElementsMatch(t, actuals[i], users)
			cursor = users[len(users)-1].ID
		}
	}
}
