package primitive

// Case :
type Case struct {
	WhenClauses [][2]any
	ElseClause  any
}

// When :
func (c *Case) When(cond, result any) *Case {
	c.WhenClauses = append(c.WhenClauses, [2]any{cond, result})
	return c
}

// Else :
func (c *Case) Else(result any) *Case {
	c.ElseClause = result
	return c
}
