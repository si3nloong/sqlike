package reflext

import (
	"reflect"
	"runtime"
	"sync"
)

// MapFunc :
type MapFunc func(*StructField) (skip bool)

// Mapper :
type Mapper struct {
	mutex   sync.Mutex
	tag     string
	cache   map[reflect.Type]*Struct
	mapFunc MapFunc
}

// NewMapperFunc :
func NewMapperFunc(tag string, mapFunc MapFunc) *Mapper {
	return &Mapper{
		cache:   make(map[reflect.Type]*Struct),
		tag:     tag,
		mapFunc: mapFunc,
	}
}

// CodecByType :
func (m *Mapper) CodecByType(t reflect.Type) *Struct {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	mapping, isOk := m.cache[t]
	if !isOk {
		mapping = getCodec(t, m.tag, m.mapFunc)
		m.cache[t] = mapping
	}
	return mapping
}

// FieldByName :
func (m *Mapper) FieldByName(v reflect.Value, name string) reflect.Value {
	v = Indirect(v)
	mustBe(v, reflect.Struct)

	tm := m.CodecByType(v.Type())
	fi, isOk := tm.Names[name]
	if !isOk {
		return v
	}
	return FieldByIndexes(v, fi.Index)
}

// LookUpFieldByName :
func (m *Mapper) LookUpFieldByName(v reflect.Value, name string) (reflect.Value, bool) {
	v = Indirect(v)
	mustBe(v, reflect.Struct)

	tm := m.CodecByType(v.Type())
	fi, isOk := tm.Names[name]
	if !isOk {
		return v, false
	}
	return FieldByIndexes(v, fi.Index), true
}

// FieldByIndexes :
func (m *Mapper) FieldByIndexes(v reflect.Value, idxs []int) reflect.Value {
	return FieldByIndexes(v, idxs)
}

// FieldByIndexesReadOnly :
func (m *Mapper) FieldByIndexesReadOnly(v reflect.Value, idxs []int) reflect.Value {
	return FieldByIndexesReadOnly(v, idxs)
}

// TraversalsByName :
func (m *Mapper) TraversalsByName(t reflect.Type, names []string) (idxs [][]int) {
	t = Deref(t)
	mustBe(t, reflect.Struct)

	idxs = make([][]int, 0, len(names))
	cdc := m.CodecByType(t)
	for _, name := range names {
		sf, exist := cdc.Names[name]
		if exist {
			idxs = append(idxs, sf.Index)
		} else {
			// idxs = append(idxs, []int{})
		}
	}
	return idxs
}

// Init :
func Init(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	if v.Kind() == reflect.Map && v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}
	return v
}

// FieldByIndexes :
func FieldByIndexes(v reflect.Value, indexes []int) reflect.Value {
	for _, i := range indexes {
		v = Indirect(v).Field(i)
		v = Init(v)
	}
	return v
}

// FieldByIndexesReadOnly :
func FieldByIndexesReadOnly(v reflect.Value, indexes []int) reflect.Value {
	for _, i := range indexes {
		if v.Kind() == reflect.Ptr && v.IsNil() {
			v = reflect.Zero(v.Type())
			break
		}
		v = Indirect(v).Field(i)
	}
	return v
}

type kinder interface {
	Kind() reflect.Kind
}

func mustBe(v kinder, k reflect.Kind) {
	if v.Kind() != k {
		panic(&reflect.ValueError{Method: methodName(), Kind: k})
	}
}

func methodName() string {
	pc, _, _, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown method"
	}
	return f.Name()
}
