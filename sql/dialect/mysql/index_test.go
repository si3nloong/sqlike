package mysql

import (
	"testing"

	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
	"github.com/si3nloong/sqlike/v2/sqlike/indexes"
	"github.com/stretchr/testify/require"
)

func TestHasIndexByName(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)
	ms.HasIndexByName(stmt, "db", "table", "idx1")
	require.Equal(t, `SELECT COUNT(1) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?;`, stmt.String())
	require.ElementsMatch(t, []interface{}{"db", "table", "idx1"}, stmt.Args())
}

func TestGetIndexes(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)

	ms.GetIndexes(stmt, "db", "table")
	require.Equal(t, "SELECT DISTINCT INDEX_NAME, INDEX_TYPE, NON_UNIQUE FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;", stmt.String())
	require.ElementsMatch(t, []interface{}{"db", "table"}, stmt.Args())
}

func TestGetIndexByType(t *testing.T) {
	ms := New()
	require.Equal(t, "FULLTEXT INDEX", ms.getIndexByType(indexes.FullText))
	require.Equal(t, "SPATIAL INDEX", ms.getIndexByType(indexes.Spatial))
	require.Equal(t, "UNIQUE INDEX", ms.getIndexByType(indexes.Unique))
	require.Equal(t, "PRIMARY KEY", ms.getIndexByType(indexes.Primary))
	require.Equal(t, "INDEX", ms.getIndexByType(0))
}
