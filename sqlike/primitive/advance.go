package primitive

// Case :
type Case struct {
	Whens []*WhenThen
}

// WhenThen :
type WhenThen struct {
	Conds  []interface{}
	Result interface{}
}

// Then :
func (wt *WhenThen) Then(result interface{}) *WhenThen {
	wt.Result = result
	return wt
}
