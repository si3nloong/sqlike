package mysql

import (
	"testing"

	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/stretchr/testify/require"
)

func TestDropTable(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)

	{
		ms.DropTable(stmt, "db", "table", true)
		require.Equal(t, "DROP TABLE IF EXISTS `db`.`table`;", stmt.String())
		require.ElementsMatch(t, []interface{}{}, stmt.Args())
	}

	stmt.Reset()

	{
		ms.DropTable(stmt, "db", "table", false)
		require.Equal(t, "DROP TABLE `db`.`table`;", stmt.String())
		require.ElementsMatch(t, []interface{}{}, stmt.Args())
	}
}

func TestTruncateTable(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)

	ms.TruncateTable(stmt, "db", "table")
	require.Equal(t, "TRUNCATE TABLE `db`.`table`;", stmt.String())
}
