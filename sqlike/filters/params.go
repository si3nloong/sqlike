package filters

// Params :
type Params struct {
	Selects []interface{}
	Filters []interface{}
	Sorts   []interface{}
	Limit   uint
}

// FindOne :
// func (p *Params) FindOne() *actions.FindOneActions {
// 	act := &actions.FindOneActions{}
// 	return act
// }
