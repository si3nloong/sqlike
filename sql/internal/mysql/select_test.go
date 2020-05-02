package mysql

import (
	"log"
	"testing"
	"time"

	"github.com/si3nloong/sqlike/sql"
	"github.com/si3nloong/sqlike/sql/expr"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/stretchr/testify/require"
)

func TestSelect(t *testing.T) {
	var (
		now = time.Now()
		err error
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
		stmt := sqlstmt.AcquireStmt(MySQL{})
		defer sqlstmt.ReleaseStmt(stmt)
		err = New().Select(
			stmt,
			actions.Find().From("A", "Test").
				Where(
					expr.And(filters...),
					expr.Or(filters...),
					expr.Equal("E", uint(888)),
					expr.NotBetween("Z", -10, 12933),
				).(*actions.FindActions), 0,
		)
		require.NoError(t, err)
		require.Equal(t, "SELECT * FROM `A`.`Test` WHERE ((`A` = ? AND `B` LIKE ? AND `DateTime` BETWEEN ? AND ?) AND (`A` = ? OR `B` LIKE ? OR `DateTime` BETWEEN ? AND ?) AND `E` = ? AND `Z` NOT BETWEEN ? AND ?);", stmt.String())
	}

	{
		stmt := sql.
			Select(
				expr.As(sql.Select().From("Testing"), "t"),
			).
			From("db", "Hell", expr.Raw("FORCE INDEX")).
			Where(
				expr.Equal("A", 1),
				expr.In("Key", sql.Select().
					From("Inner").
					Where(
						expr.Equal(expr.Column("x", "Z"), true),
					).
					OrderBy(
						expr.Desc(expr.Column("x", "G")),
					).
					Having(
						expr.Equal("L", "ok"),
					).
					Limit(10).
					Offset(1),
				),
				expr.NotNull("B"),
				expr.Or(
					expr.In("C", []string{"1", "2", "3"}),
					expr.Equal("A", 100),
					expr.Between("Time", time.Now(), time.Now().Add(time.Minute*10)),
				),
			).OrderBy(
			expr.Desc("A"),
			expr.Asc("C"),
		).Limit(1)

		x := New()
		stmt2 := sqlstmt.NewStatement(x)
		x.parser.BuildStatement(stmt2, stmt)
		log.Println(stmt2.String())
		// New().Select(stmt)
	}
}
