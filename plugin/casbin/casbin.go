package casbin

import (
	"errors"
	"strings"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/indexes"
	"github.com/si3nloong/sqlike/sqlike/options"
)

// Adapter :
type Adapter struct {
	table    *sqlike.Table
	filtered bool
}

// MustNew :
func MustNew(table *sqlike.Table) persist.FilteredAdapter {
	a, err := New(table)
	if err != nil {
		panic(err)
	}
	return a
}

// New :
func New(table *sqlike.Table) (persist.FilteredAdapter, error) {
	if table == nil {
		return nil, errors.New("invalid <nil> table")
	}
	a := &Adapter{
		table: table,
	}
	if err := a.createTable(); err != nil {
		return nil, err
	}
	if err := a.table.Indexes().
		CreateOneIfNotExists(indexes.Index{
			Name: "PrimaryKey",
			Type: indexes.Unique,
			Columns: []indexes.Column{
				indexes.Column{Name: "PType"},
				indexes.Column{Name: "V0"},
				indexes.Column{Name: "V1"},
				indexes.Column{Name: "V2"},
				indexes.Column{Name: "V3"},
				indexes.Column{Name: "V4"},
				indexes.Column{Name: "V5"},
			},
		}); err != nil {
		return nil, err
	}
	return a, nil
}

// LoadPolicy :
func (a *Adapter) LoadPolicy(model model.Model) error {
	result, err := a.table.Find(nil)
	if err != nil {
		return err
	}

	policies := []*Policy{}
	if err := result.All(&policies); err != nil {
		return err
	}

	for _, r := range policies {
		loadPolicy(r, model)
	}
	return nil
}

// LoadFilteredPolicy :
func (a *Adapter) LoadFilteredPolicy(model model.Model, filter interface{}) error {
	var policies []*Policy
	x, ok := filter.([]interface{})
	if !ok {
		return errors.New("invalid filter data type, expected []interface{}")
	}

	result, err := a.table.Find(actions.Find().Where(x...))
	if err != nil {
		return err
	}

	if err := result.All(&policies); err != nil {
		return err
	}

	for _, policy := range policies {
		loadPolicy(policy, model)
	}
	a.filtered = true
	return nil
}

// SavePolicy : saves all policy rules to the storage.
func (a *Adapter) SavePolicy(model model.Model) error {
	var policies []*Policy
	for ptype, ast := range model["p"] {
		for _, r := range ast.Policy {
			policies = append(policies, toPermissionRule(ptype, r))
		}
	}

	for ptype, ast := range model["g"] {
		for _, r := range ast.Policy {
			policies = append(policies, toPermissionRule(ptype, r))
		}
	}

	if len(policies) > 0 {
		if _, err := a.table.
			Insert(&policies, options.Insert().
				SetMode(options.InsertOnDuplicate)); err != nil {
			return err
		}
	}
	return nil
}

// AddPolicy : adds a policy policy to the storage. This is part of the Auto-Save feature.
func (a *Adapter) AddPolicy(sec string, ptype string, rules []string) error {
	if _, err := a.table.InsertOne(
		toPermissionRule(ptype, rules),
		options.InsertOne().
			SetMode(options.InsertIgnore),
	); err != nil {
		return err
	}
	return nil
}

// RemovePolicy : removes a policy policy from the storage. This is part of the Auto-Save feature.
func (a *Adapter) RemovePolicy(sec string, ptype string, rules []string) error {
	policy := toPermissionRule(ptype, rules)
	if _, err := a.table.DeleteOne(
		actions.DeleteOne().Where(
			expr.Equal("PType", policy.PType),
			expr.Equal("V0", policy.V0),
			expr.Equal("V1", policy.V1),
			expr.Equal("V2", policy.V2),
			expr.Equal("V3", policy.V3),
			expr.Equal("V4", policy.V4),
			expr.Equal("V5", policy.V5),
		),
	); err != nil {
		return err
	}
	return nil
}

// RemoveFilteredPolicy : removes policy rules that match the filter from the storage. This is part of the Auto-Save feature.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, idx int, values ...string) error {
	policy := new(Policy)
	policy.PType = ptype
	length := len(values)
	if idx <= 0 && 0 < idx+length {
		policy.V0 = values[0-idx]
	}
	if idx <= 1 && 1 < idx+length {
		policy.V1 = values[1-idx]
	}
	if idx <= 2 && 2 < idx+length {
		policy.V2 = values[2-idx]
	}
	if idx <= 3 && 3 < idx+length {
		policy.V3 = values[3-idx]
	}
	if idx <= 4 && 4 < idx+length {
		policy.V4 = values[4-idx]
	}
	if idx <= 5 && 5 < idx+length {
		policy.V5 = values[5-idx]
	}
	return nil
}

// IsFiltered :
func (a *Adapter) IsFiltered() bool {
	return a.filtered
}

func loadPolicy(policy *Policy, model model.Model) {
	const prefixLine = ", "
	var sb strings.Builder

	sb.WriteString(policy.PType)
	if len(policy.V0) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(policy.V0)
	}
	if len(policy.V1) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(policy.V1)
	}
	if len(policy.V2) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(policy.V2)
	}
	if len(policy.V3) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(policy.V3)
	}
	if len(policy.V4) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(policy.V4)
	}
	if len(policy.V5) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(policy.V5)
	}

	persist.LoadPolicyLine(sb.String(), model)
}

func toPermissionRule(ptype string, rules []string) *Policy {
	policy := new(Policy)
	policy.PType = ptype
	length := len(rules)
	if length > 0 {
		policy.V0 = rules[0]
	}
	if length > 1 {
		policy.V1 = rules[1]
	}
	if length > 2 {
		policy.V2 = rules[2]
	}
	if length > 3 {
		policy.V3 = rules[3]
	}
	if length > 4 {
		policy.V4 = rules[4]
	}
	if length > 5 {
		policy.V5 = rules[5]
	}
	return policy
}

func (a *Adapter) createTable() error {
	return a.table.UnsafeMigrate(Policy{})
}
