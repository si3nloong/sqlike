package mysql

import (
	"testing"
	"time"

	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/internal/primitive"
	"github.com/si3nloong/sqlike/v2/sql"
	"github.com/si3nloong/sqlike/v2/sql/expr"
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
	"github.com/stretchr/testify/require"
)

func TestSelect(t *testing.T) {
	var (
		now = time.Now()
		err error
		ms  = New()
	)

	invalids := []any{
		expr.And(),
		nil,
		struct{}{},
		expr.Or(),
		make([]any, 0),
		[]any{},
		[]any(nil),
		map[string]string(nil),
	}

	filters := []any{
		expr.Equal("A", 1),
		expr.Like("B", "abc%"),
		expr.Between("DateTime", now, now.Add(5*time.Minute)),
	}
	filters = append(filters, invalids...)

	// Complex select statement
	{

		act := actions.Find().
			Where(
				expr.And(filters...),
				expr.Or(filters...),
				expr.Equal("E", uint(888)),
				expr.NotBetween("Z", -10, 12933),
			)

		v, ok := act.(*actions.FindActions)
		require.True(t, ok)
		v.Database = "A"
		v.Table = "Test"

		stmt := sqlstmt.AcquireStmt(ms)
		defer sqlstmt.ReleaseStmt(stmt)
		err = New().Select(
			stmt,
			v,
			primitive.Lock{},
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

		stmt2 := sqlstmt.AcquireStmt(ms)
		defer sqlstmt.ReleaseStmt(stmt2)
		err := ms.parser.BuildStatement(stmt2, stmt)
		require.NoError(t, err)
	}
}
