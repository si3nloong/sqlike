package examples

import (
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	plugin "github.com/si3nloong/sqlike/plugin/casbin"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/stretchr/testify/require"
)

// CasbinExamples :
func CasbinExamples(t *testing.T, db *sqlike.Database) {
	var (
		a   persist.FilteredAdapter
		e   *casbin.Enforcer
		err error
		ok  bool
	)

	// Init policy
	{
		a = plugin.MustNew(db.Table("AccessPolicy"))
		a = plugin.MustNew(db.Table("AccessPolicy"))
		e, err = casbin.NewEnforcer("./rbac_model.conf", a)
		require.NoError(t, err)
		err = e.LoadModel()
		require.NoError(t, err)
		err = e.LoadPolicy()
		require.NoError(t, err)
	}

	// Create policy
	{
		ok, err = e.AddNamedPolicy("p", "casbin", "/*", "GET")
		require.True(t, ok)
		require.NoError(t, err)
		e.AddNamedPolicy("p", "username", "/*", "*")
		require.True(t, ok)
		require.NoError(t, err)
		e.AddGroupingPolicy("admin", "tester", "/*")
		e.AddPolicy("admin", "/login", "POST")
		e.AddPolicy("admin", "/logout", "POST")
		e.AddNamedPolicy("p", "admin", "/login", "POST")
		e.AddNamedPolicy("p", "admin", "/login", "POST")

		policies := e.GetFilteredPolicy(0, "admin")
		require.ElementsMatch(t, [][]string{
			[]string{"admin", "/login", "POST"},
			[]string{"admin", "/logout", "POST"},
		}, policies)

		err = e.SavePolicy()
		require.NoError(t, err)
	}

	// Check success access
	{
		ok, err = e.Enforce("username", "/*", "*")
		require.True(t, ok)
		require.NoError(t, err)
	}

	// Check failed access
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

}
