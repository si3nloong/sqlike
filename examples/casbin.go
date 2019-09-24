package examples

import (
	"testing"

	"github.com/casbin/casbin/v2"
	plugin "github.com/si3nloong/sqlike/plugin/casbin"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/stretchr/testify/require"
)

// CasbinExamples :
func CasbinExamples(t *testing.T, db *sqlike.Database) {
	var (
		ok  bool
		e   *casbin.Enforcer
		err error
	)
	a := plugin.New(db)
	e, err = casbin.NewEnforcer("./rbac_model.conf", a)
	require.NoError(t, err)
	e.LoadModel()
	err = e.LoadPolicy()

	{
		ok, err = e.Enforce("username", "/*", "*")
		require.True(t, ok)
		require.NoError(t, err)
	}

	{
		ok, err = e.Enforce("s1", "/*", "*")
		require.False(t, ok)
		require.NoError(t, err)
	}

	{
		ok, err = e.Enforce("admin", "/login", "POST")
		require.True(t, ok)
		require.NoError(t, err)

		ok, err = e.Enforce("admin", "/login", "GET")
		require.False(t, ok)
		require.NoError(t, err)

		ok, err = e.Enforce("admin", "/logout", "*")
		require.False(t, ok)
		require.NoError(t, err)
	}

	e.SavePolicy()
}
