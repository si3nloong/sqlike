package casbin

import (
	"strings"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/si3nloong/sqlike/sqlike"
)

// Adapter :
type Adapter struct {
	db *sqlike.Database
}

// New :
func New(db *sqlike.Database) *Adapter {
	a := &Adapter{
		db: db,
	}

	a.createTable()
	return a
}

var _ persist.Adapter = new(Adapter)

// LoadPolicy :
func (a *Adapter) LoadPolicy(model model.Model) error {
	result, err := a.db.Table("PermissionRule").Find(nil)
	if err != nil {
		return err
	}

	rules := []*PermissionRule{}
	if err := result.All(&rules); err != nil {
		return err
	}

	for _, r := range rules {
		loadPolicyLine(r, model)
	}
	return nil
}

// SavePolicy : saves all policy rules to the storage.
func (a *Adapter) SavePolicy(model model.Model) error {
	// TODO : save casbin policy
	return nil
}

// AddPolicy : adds a policy rule to the storage. This is part of the Auto-Save feature.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	// TODO : add casbin policy
	return nil
}

// RemovePolicy : removes a policy rule from the storage. This is part of the Auto-Save feature.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	// TODO : remove casbin policy
	return nil
}

// RemoveFilteredPolicy : removes policy rules that match the filter from the storage. This is part of the Auto-Save feature.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	// TODO : remove casbin filter policy
	return nil
}

func loadPolicyLine(line *PermissionRule, model model.Model) {
	const prefixLine = ", "
	var sb strings.Builder

	sb.WriteString(line.PType)
	if len(line.V0) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(line.V0)
	}
	if len(line.V1) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(line.V1)
	}
	if len(line.V2) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(line.V2)
	}
	if len(line.V3) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(line.V3)
	}
	if len(line.V4) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(line.V4)
	}
	if len(line.V5) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(line.V5)
	}

	persist.LoadPolicyLine(sb.String(), model)
}

func (a *Adapter) createTable() {
	a.db.Table("PermissionRule").MustUnsafeMigrate(PermissionRule{})
}
