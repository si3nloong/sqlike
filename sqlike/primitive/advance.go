package primitive

// Case :
type Case struct {
	WhenClauses [][2]interface{}
	ElseClause  interface{}
}

// When :
func (c *Case) When(cond, result interface{}) *Case {
	c.WhenClauses = append(c.WhenClauses, [2]interface{}{cond, result})
	return c
}

// Else :
func (c *Case) Else(result interface{}) *Case {
	c.ElseClause = result
	return c
}
