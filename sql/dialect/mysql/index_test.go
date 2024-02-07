package mysql

import (
	"testing"

	"github.com/si3nloong/sqlike/v2/sql"
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
	"github.com/stretchr/testify/require"
)

func TestHasIndexByName(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)
	ms.HasIndexByName(stmt, "db", "table", "idx1")
	require.Equal(t, `SELECT COUNT(1) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?;`, stmt.String())
	require.ElementsMatch(t, []any{"db", "table", "idx1"}, stmt.Args())
}

func TestGetIndexes(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)

	ms.GetIndexes(stmt, "db", "table")
	require.Equal(t, "SELECT DISTINCT INDEX_NAME, INDEX_TYPE, NON_UNIQUE FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;", stmt.String())
	require.ElementsMatch(t, []any{"db", "table"}, stmt.Args())
}

func TestGetIndexByType(t *testing.T) {
	ms := New()
	require.Equal(t, "FULLTEXT INDEX", ms.getIndexByType(sql.FullText))
	require.Equal(t, "SPATIAL INDEX", ms.getIndexByType(sql.Spatial))
	require.Equal(t, "UNIQUE INDEX", ms.getIndexByType(sql.Unique))
	require.Equal(t, "PRIMARY KEY", ms.getIndexByType(sql.Primary))
	require.Equal(t, "INDEX", ms.getIndexByType(0))
}
