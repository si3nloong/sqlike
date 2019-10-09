package rsql

// Params :
type Params struct {
	Selects []interface{}
	Filters []interface{}
	Sorts   []interface{}
	Limit   uint
}
