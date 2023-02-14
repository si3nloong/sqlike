package mysql

import (
	"testing"

	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
	"github.com/stretchr/testify/require"
)

func TestGetColumns(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)
	ms.GetColumns(stmt, "db", "table")
	require.Equal(t, `SELECT ORDINAL_POSITION, COLUMN_NAME, COLUMN_TYPE, COLUMN_DEFAULT, IS_NULLABLE,
	DATA_TYPE, CHARACTER_SET_NAME, COLLATION_NAME, COLUMN_COMMENT, EXTRA FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION;`, stmt.String())
	require.ElementsMatch(t, []any{"db", "table"}, stmt.Args())
}

func TestRenameColumn(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)
	ms.RenameColumn(stmt, "db", "table", "c1", "_column_")
	require.Equal(t, "ALTER TABLE `db`.`table` RENAME COLUMN `c1` TO `_column_`;", stmt.String())
	require.ElementsMatch(t, []any{}, stmt.Args())
}

func TestDropColumn(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)

	t.Cleanup(func() {
		sqlstmt.ReleaseStmt(stmt)
	})
	ms.DropColumn(stmt, "db", "table", "c1")
	require.Equal(t, "ALTER TABLE `db`.`table` DROP COLUMN `c1`;", stmt.String())
	require.ElementsMatch(t, []any{}, stmt.Args())
}
