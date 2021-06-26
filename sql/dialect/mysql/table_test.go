package mysql

import (
	"testing"

	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
	"github.com/stretchr/testify/require"
)

func TestHasPrimaryKey(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)

	ms.HasPrimaryKey(stmt, "db", "table")
	require.Equal(t, "SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND CONSTRAINT_TYPE = 'PRIMARY KEY';", stmt.String())
	require.ElementsMatch(t, []interface{}{"db", "table"}, stmt.Args())
}

func TestDropTable(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)

	{
		ms.DropTable(stmt, "db", "table", true, false)
		require.Equal(t, "DROP TABLE IF EXISTS `db`.`table`;", stmt.String())
		require.ElementsMatch(t, []interface{}{}, stmt.Args())
	}

	stmt.Reset()

	{
		ms.DropTable(stmt, "db", "table", false, false)
		require.Equal(t, "DROP TABLE `db`.`table`;", stmt.String())
		require.ElementsMatch(t, []interface{}{}, stmt.Args())
	}
}

func TestRenameTable(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)
	ms.RenameTable(stmt, "db", "oldName", "newName")
	require.Equal(t, "RENAME TABLE `db`.`oldName` TO `db`.`newName`;", stmt.String())
	require.ElementsMatch(t, []interface{}{}, stmt.Args())
}

func TestTruncateTable(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)

	ms.TruncateTable(stmt, "db", "table")
	require.Equal(t, "TRUNCATE TABLE `db`.`table`;", stmt.String())
	require.ElementsMatch(t, []interface{}{}, stmt.Args())
}

func TestHasTable(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)

	ms.HasTable(stmt, "db", "table")
	require.Equal(t, "SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;", stmt.String())
	require.ElementsMatch(t, []interface{}{"db", "table"}, stmt.Args())

}
