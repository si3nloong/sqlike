package examples

import (
	"testing"
	"time"

	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/stretchr/testify/require"
)

// User :
type User struct {
	ID        int64
	Name      string
	Age       int
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

	{
		_, err = table.InsertMany(
			[]User{
				User{10, "User A", 18, time.Now()},
				User{88, "User B", 12, time.Now()},
				User{8, "User F", 20, time.Now()},
				User{27, "User C", 16, time.Now()},
				User{20, "User C", 16, time.Now()},
				User{21, "User C", 16, time.Now()},
				User{50, "User D", 23, time.Now()},
				User{5, "User E", 30, time.Now()},
			},
			options.InsertMany().SetDebug(true))
		require.NoError(t, err)
	}

	{
		table.Paginate(actions.Paginate().
			Where().Limit(100),
			options.Paginate().
				SetCursor(20).
				SetDebug(true))
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
