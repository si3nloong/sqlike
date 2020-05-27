package mysql

import (
	"testing"

	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/stretchr/testify/require"
)

func TestGetColumns(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)
	ms.GetColumns(stmt, "db", "table")
	require.Equal(t, `SELECT ORDINAL_POSITION, COLUMN_NAME, COLUMN_TYPE, COLUMN_DEFAULT, IS_NULLABLE,
	DATA_TYPE, CHARACTER_SET_NAME, COLLATION_NAME, EXTRA FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;`, stmt.String())
	require.ElementsMatch(t, []interface{}{"db", "table"}, stmt.Args())
}

func TestRenameColumn(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)
	ms.RenameTable(stmt, "db", "oldName", "newName")
	require.Equal(t, "RENAME TABLE `db`.`oldName` TO `db`.`newName`;", stmt.String())
	require.ElementsMatch(t, []interface{}{}, stmt.Args())
}

func TestDropColumn(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)
	ms.DropColumn(stmt, "db", "table", "c1")
	require.Equal(t, "ALTER TABLE `db`.`table` DROP COLUMN `c1`;", stmt.String())
	require.ElementsMatch(t, []interface{}{}, stmt.Args())
}
