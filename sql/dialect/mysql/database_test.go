package mysql

import (
	"testing"

	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
	"github.com/stretchr/testify/require"
)

func TestUseDatabase(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)
	ms.UseDatabase(stmt, "db")
	require.Equal(t, "USE `db`;", stmt.String())
	require.ElementsMatch(t, []any{}, stmt.Args())
}

func TestCreateDatabase(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)

	{
		ms.CreateDatabase(stmt, "db", false)
		require.Equal(t, "CREATE DATABASE `db`;", stmt.String())
		require.ElementsMatch(t, []any{}, stmt.Args())
	}

	stmt.Reset()

	{
		ms.CreateDatabase(stmt, "db", true)
		require.Equal(t, "CREATE DATABASE IF NOT EXISTS `db`;", stmt.String())
		require.ElementsMatch(t, []any{}, stmt.Args())
	}
}

func TestDropDatabase(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)

	{
		ms.DropDatabase(stmt, "db", false)
		require.Equal(t, "DROP SCHEMA `db`;", stmt.String())
		require.ElementsMatch(t, []any{}, stmt.Args())
	}

	stmt.Reset()

	{
		ms.DropDatabase(stmt, "db", true)
		require.Equal(t, "DROP SCHEMA IF EXISTS `db`;", stmt.String())
		require.ElementsMatch(t, []any{}, stmt.Args())
	}
}

func TestGetDatabases(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)
	ms.GetDatabases(stmt)
	require.Equal(t, "SHOW DATABASES;", stmt.String())
	require.ElementsMatch(t, []any{}, stmt.Args())
}
