package casbin

import (
	"context"
	"errors"
	"strings"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/indexes"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/si3nloong/sqlike/x/primitive"
)

// Adapter :
type Adapter struct {
	ctx      context.Context
	table    *sqlike.Table
	filtered bool
}

var _ persist.FilteredAdapter = new(Adapter)

// MustNew :
func MustNew(ctx context.Context, table *sqlike.Table) persist.FilteredAdapter {
	a, err := New(ctx, table)
	if err != nil {
		panic(err)
	}
	return a
}

// New :
func New(ctx context.Context, table *sqlike.Table) (persist.FilteredAdapter, error) {
	if table == nil {
		return nil, errors.New("invalid <nil> table")
	}
	a := &Adapter{
		ctx:   ctx,
		table: table,
	}
	if err := a.createTable(); err != nil {
		return nil, err
	}
	if err := a.table.Indexes().
		CreateOneIfNotExists(
			a.ctx,
			indexes.Index{
				Type: indexes.Primary,
				Columns: indexes.Columns(
					"PType",
					"V0", "V1", "V2",
					"V3", "V4", "V5",
				),
			}); err != nil {
		return nil, err
	}
	return a, nil
}

// LoadPolicy :
func (a *Adapter) LoadPolicy(model model.Model) error {
	result, err := a.table.Find(a.ctx, nil)
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
	x, ok := filter.(primitive.Group)
	if !ok {
		return errors.New("invalid filter data type, expected []interface{}")
	}

	act := new(actions.FindActions)
	act.Conditions = x
	result, err := a.table.Find(
		a.ctx,
		act,
		options.Find().
			SetNoLimit(true).
			SetDebug(true),
	)
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
			policies = append(policies, toPolicy(ptype, r))
		}
	}

	for ptype, ast := range model["g"] {
		for _, r := range ast.Policy {
			policies = append(policies, toPolicy(ptype, r))
		}
	}

	if len(policies) > 0 {
		if _, err := a.table.Insert(
			a.ctx,
			&policies,
			options.Insert().
				SetMode(options.InsertOnDuplicate),
		); err != nil {
			return err
		}
	}
	return nil
}

// AddPolicy : adds a policy policy to the storage. This is part of the Auto-Save feature.
func (a *Adapter) AddPolicy(sec string, ptype string, rules []string) error {
	if _, err := a.table.InsertOne(
		a.ctx,
		toPolicy(ptype, rules),
		options.InsertOne().
			SetMode(options.InsertIgnore),
	); err != nil {
		return err
	}
	return nil
}

// RemovePolicy : removes a policy policy from the storage. This is part of the Auto-Save feature.
func (a *Adapter) RemovePolicy(sec string, ptype string, rules []string) error {
	policy := toPolicy(ptype, rules)
	if _, err := a.table.DeleteOne(
		a.ctx,
		actions.DeleteOne().
			Where(
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

func loadPolicy(policy *Policy, m model.Model) {
	tokens := append([]string{}, policy.PType)

	if len(policy.V0) > 0 {
		tokens = append(tokens, policy.V0)
	}
	if len(policy.V1) > 0 {
		tokens = append(tokens, policy.V1)
	}
	if len(policy.V2) > 0 {
		tokens = append(tokens, policy.V2)
	}
	if len(policy.V3) > 0 {
		tokens = append(tokens, policy.V3)
	}
	if len(policy.V4) > 0 {
		tokens = append(tokens, policy.V4)
	}
	if len(policy.V5) > 0 {
		tokens = append(tokens, policy.V5)
	}

	key := tokens[0]
	sec := key[:1]
	m[sec][key].Policy = append(m[sec][key].Policy, tokens[1:])
	m[sec][key].PolicyMap[strings.Join(tokens[1:], model.DefaultSep)] = len(m[sec][key].Policy) - 1
}

func toPolicy(ptype string, rules []string) *Policy {
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
	return a.table.UnsafeMigrate(a.ctx, Policy{})
}
