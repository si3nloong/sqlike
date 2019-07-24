package reflext

import "reflect"

// MapKeys :
type MapKeys []reflect.Value

func (x MapKeys) Len() int { return len(x) }

func (x MapKeys) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x MapKeys) Less(i, j int) bool {
	return x[i].String() < x[j].String()
}
