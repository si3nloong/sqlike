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
		_, err = table.Insert(data, options.
			Insert().SetDebug(true))
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

	// Paginate with simple query
	{
		pg, err := table.Paginate(actions.Paginate().
			Where(
				expr.GreaterOrEqual("Age", 0),
			).
			OrderBy(
				expr.Desc("Age"),
			).Limit(1),
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
			require.Equal(t, data[i], users[0])
			cursor = users[len(users)-1].ID
			if pg.NextPage(cursor) != nil {
				break
			}
		}
	}

	length := 4
	actuals := [][]User{
		data[:length],
		data[length:(length * 2)],
		data[length*2:],
	}

	cursor = nil
	pg, err := table.Paginate(actions.Paginate().
		OrderBy(
			expr.Desc("Age"),
		).Limit(uint(length)),
		options.Paginate().
			SetDebug(true))
	require.NoError(t, err)

	// Expected paginate with error
	{
		err = pg.NextPage(nil)
		require.Error(t, err)
		err = pg.NextPage([]string{})
		require.Error(t, err)
		var nilslice []string
		err = pg.NextPage(nilslice)
		require.Error(t, err)
		var nilmap map[string]interface{}
		err = pg.NextPage(nilmap)
		require.Error(t, err)
		err = pg.NextPage("")
		require.Error(t, err)
		err = pg.NextPage(0)
		require.Error(t, err)
		err = pg.NextPage(false)
		require.Error(t, err)
		err = pg.NextPage(float64(0))
		require.Error(t, err)
		err = pg.NextPage([]byte(nil))
		require.Error(t, err)
	}

	{
		for i := 0; i < len(actuals); i++ {
			users = []User{}
			err = pg.All(&users)
			require.NoError(t, err)
			if len(users) == 0 {
				break
			}
			require.ElementsMatch(t, actuals[i], users)
			cursor = users[len(users)-1].ID
			if pg.NextPage(cursor) != nil {
				break
			}
		}
	}
}
