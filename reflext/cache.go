package reflext

import (
	"reflect"
	"runtime"
	"sync"
)

// DefaultMapper :
var DefaultMapper = NewMapperFunc("sqlike", nil)

// StructMapper :
type StructMapper interface {
	CodecByType(t reflect.Type) Structer
	FieldByName(v reflect.Value, name string) reflect.Value
	FieldByIndexes(v reflect.Value, idxs []int) reflect.Value
	FieldByIndexesReadOnly(v reflect.Value, idxs []int) reflect.Value
	LookUpFieldByName(v reflect.Value, name string) (reflect.Value, bool)
	TraversalsByName(t reflect.Type, names []string) (idxs [][]int)
	TraversalsByNameFunc(t reflect.Type, names []string, fn func(int, []int)) (idxs [][]int)
}

// MapFunc :
type MapFunc func(StructFielder) (skip bool)

// FormatFunc :
type FormatFunc func(string) string

// Mapper :
type Mapper struct {
	mutex   sync.Mutex
	tag     string
	cache   map[reflect.Type]*Struct
	fmtFunc FormatFunc
}

var _ StructMapper = (*Mapper)(nil)

// NewMapperFunc :
func NewMapperFunc(tag string, fmtFunc FormatFunc) *Mapper {
	return &Mapper{
		cache:   make(map[reflect.Type]*Struct),
		tag:     tag,
		fmtFunc: fmtFunc,
	}
}

// CodecByType :
func (m *Mapper) CodecByType(t reflect.Type) Structer {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	mapping, ok := m.cache[t]
	if !ok {
		mapping = getCodec(t, m.tag, m.fmtFunc)
		m.cache[t] = mapping
	}
	return mapping
}

// FieldByName : get reflect.Value from struct by field name
func (m *Mapper) FieldByName(v reflect.Value, name string) reflect.Value {
	v = Indirect(v)
	mustBe(v, reflect.Struct)

	tm := m.CodecByType(v.Type())
	fi, ok := tm.LookUpFieldByName(name)
	if !ok {
		return v
	}
	return FieldByIndexes(v, fi.Index())
}

// FieldByIndexes : get reflect.Value from struct by indexes. If the reflect.Value is nil, it will get initialized
func (m *Mapper) FieldByIndexes(v reflect.Value, idxs []int) reflect.Value {
	return FieldByIndexes(v, idxs)
}

// FieldByIndexesReadOnly : get reflect.Value from struct by indexes without initialized
func (m *Mapper) FieldByIndexesReadOnly(v reflect.Value, idxs []int) reflect.Value {
	return FieldByIndexesReadOnly(v, idxs)
}

// LookUpFieldByName : lookup reflect.Value from struct by field name
func (m *Mapper) LookUpFieldByName(v reflect.Value, name string) (reflect.Value, bool) {
	v = Indirect(v)
	mustBe(v, reflect.Struct)

	tm := m.CodecByType(v.Type())
	fi, ok := tm.LookUpFieldByName(name)
	if !ok {
		return v, false
	}
	return FieldByIndexes(v, fi.Index()), true
}

// TraversalsByName :
func (m *Mapper) TraversalsByName(t reflect.Type, names []string) (idxs [][]int) {
	idxs = make([][]int, 0, len(names))
	m.TraversalsByNameFunc(t, names, func(i int, idx []int) {
		if idxs != nil {
			idxs = append(idxs, idx)
		} else {
			idxs = append(idxs, nil)
		}
	})
	return idxs
}

// TraversalsByNameFunc :
func (m *Mapper) TraversalsByNameFunc(t reflect.Type, names []string, fn func(int, []int)) (idxs [][]int) {
	t = Deref(t)
	mustBe(t, reflect.Struct)

	idxs = make([][]int, 0, len(names))
	cdc := m.CodecByType(t)
	for i, name := range names {
		sf, ok := cdc.LookUpFieldByName(name)
		if ok {
			fn(i, sf.Index())
		} else {
			fn(i, nil)
		}
	}
	return idxs
}

// FieldByIndexes : get reflect.Value from struct by indexes. If the reflect.Value is nil, it will get initialized
func FieldByIndexes(v reflect.Value, indexes []int) reflect.Value {
	for _, i := range indexes {
		v = Indirect(v).Field(i)
		v = Init(v)
	}
	return v
}

// FieldByIndexesReadOnly : get reflect.Value from struct by indexes without initialized
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
