package examples

import (
	"log"
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
		User{10, "User A", 18, StatusActive, time.Now()},
		User{88, "User B", 12, StatusActive, time.Now()},
		User{8, "User F", 20, StatusActive, time.Now()},
		User{27, "User C", 16, StatusSuspend, time.Now()},
		User{20, "User C", 16, StatusActive, time.Now()},
		User{100, "User G", 10, StatusSuspend, time.Now()},
		User{21, "User C", 16, StatusActive, time.Now()},
		User{50, "User D", 23, StatusActive, time.Now()},
		User{5, "User E", 30, StatusSuspend, time.Now()},
	}

	{
		_, err = table.InsertMany(data, options.
			InsertMany().SetDebug(true))
		require.NoError(t, err)
	}

	{
		var (
		// users []User
		// cursor interface{}
		)

		paginator, err := table.Paginate(actions.Paginate().
			Where().
			OrderBy(
				expr.Desc("Age"),
			).Limit(100),
			options.Paginate().
				SetDebug(true))
		require.NoError(t, err)

		log.Println(paginator)
		// for {

		// 	err = result.All(&users)
		// 	require.NoError(t, err)
		// 	if len(users) == 0 {
		// 		break
		// 	}
		// 	cursor = users[len(users)-1].ID
		// }

	}

	// limits := uint(2)

	// SELECT * FROM `NormalStruct` ORDER BY `Float64`, `$Key` LIMIT 3;
	//

	// result, err = db.Table("NormalStruct").
	// 	Find(actions.Find().
	// 		OrderBy(
	// 			expr.Asc("Float64"),
	// 			expr.Asc("$Key"),
	// 		).
	// 		Limit(limits+1),
	// 		options.Find().SetDebug(true))
	// require.NoError(t, err)

	// err = result.All(&nss)
	// require.NoError(t, err)
	// length := len(nss)
	// if uint(length) > limits {
	// 	// csr := nss[length-1].ID.String()
	// 	nss = nss[:length-1]
	// }
}
