package examples

import (
	"context"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/si3nloong/sqlike/v2"
	plugin "github.com/si3nloong/sqlike/v2/plugin/casbin"
	"github.com/si3nloong/sqlike/v2/sql/expr"
	"github.com/stretchr/testify/require"
)

// CasbinExamples :
func CasbinExamples(ctx context.Context, t *testing.T, db *sqlike.Database) {
	var (
		a   persist.FilteredAdapter
		e   *casbin.Enforcer
		err error
		ok  bool
	)

	table := db.Table("AccessPolicy")
	// Init policy
	{
		err = table.DropIfExists(ctx)
		require.NoError(t, err)
		a = plugin.MustNew(ctx, table)
		e, err = casbin.NewEnforcer("./rbac_model.conf", a)
		require.NoError(t, err)
		err = e.LoadModel()
		require.NoError(t, err)
		err = e.LoadPolicy()
		require.NoError(t, err)
	}

	adminRules := [...][]string{
		{"admin", "/login", "POST"},
		{"admin", "/logout", "POST"},
		{"admin", "/dashboard", "GET"},
	}

	marketingRules := [...][]string{
		{"marketing", "/dashboard", "GET"},
	}

	// Create policy
	{
		ok, err = e.AddNamedPolicy("p", "casbin", "/*", "GET")
		require.True(t, ok)
		require.NoError(t, err)
		_, err = e.AddNamedPolicy("p", "username", "/*", "*")
		require.True(t, ok)
		require.NoError(t, err)
		_, err = e.AddGroupingPolicy("admin", "tester", "/*")
		require.NoError(t, err)
		_, err = e.AddPolicy(adminRules[0])
		require.NoError(t, err)
		_, err = e.AddPolicy(adminRules[1])
		require.NoError(t, err)
		_, err = e.AddPolicy(adminRules[2])
		require.NoError(t, err)
		_, err = e.AddPolicy(marketingRules[0])
		require.NoError(t, err)
		_, err = e.AddNamedPolicy("p", "admin", "/login", "POST")
		require.NoError(t, err)
		_, err = e.AddNamedPolicy("p", "admin", "/login", "POST")
		require.NoError(t, err)

		adminPolicies := e.GetFilteredPolicy(0, "admin")
		require.ElementsMatch(t, adminRules, adminPolicies)

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

	// check permission
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

	// Remove Policy
	{

	}

	// Query Policy with where conditions
	{
		e.ClearPolicy()
		err = e.LoadFilteredPolicy(
			plugin.Filter(
				expr.Equal("V0", "admin"),
			),
		)
		require.NoError(t, err)
		require.ElementsMatch(t, adminRules, e.GetPolicy())

		err = e.LoadFilteredPolicy(
			plugin.Filter(
				expr.Equal("V0", "marketing"),
			),
		)
		require.NoError(t, err)
		require.ElementsMatch(t, marketingRules, e.GetPolicy())
	}

}
