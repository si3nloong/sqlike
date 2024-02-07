package examples

func filter[I any, O any](v []I, cb func(I) O) []O {
	o := make([]O, len(v))
	for idx := range v {
		o[idx] = cb(v[idx])
	}
	return o
}
