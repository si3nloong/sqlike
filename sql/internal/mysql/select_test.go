package mysql

import (
	"testing"
	"time"

	"github.com/si3nloong/sqlike/sql/expr"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/stretchr/testify/require"
)

func TestSelect(t *testing.T) {
	var (
		now  = time.Now()
		err  error
		stmt *sqlstmt.Statement
	)

	invalids := []interface{}{
		expr.And(),
		nil,
		struct{}{},
		expr.Or(),
		make([]interface{}, 0),
		[]interface{}{},
		[]interface{}(nil),
		map[string]string(nil),
	}

	filters := []interface{}{
		expr.Equal("A", 1),
		expr.Like("B", "abc%"),
		expr.Between("DateTime", now, now.Add(5*time.Minute)),
	}
	filters = append(filters, invalids...)

	// Complex select statement
	{
		stmt, err = New().Select(
			actions.Find().From("Test").
				Where(
					expr.And(filters...),
					expr.Or(filters...),
					expr.Equal("E", uint(888)),
					expr.NotBetween("Z", -10, 12933),
				).(*actions.FindActions), 0,
		)
		require.NoError(t, err)
		require.Equal(t, "SELECT * FROM `Test` WHERE ((`A` = ? AND `B` LIKE ? AND `DateTime` BETWEEN ? AND ?) AND (`A` = ? OR `B` LIKE ? OR `DateTime` BETWEEN ? AND ?) AND `E` = ? AND `Z` NOT BETWEEN ? AND ?);", stmt.String())
	}
}
