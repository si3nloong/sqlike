package filters

// Params :
type Params struct {
	Selects []interface{}
	Filters []interface{}
	Sorts   []interface{}
	Limit   uint
}
