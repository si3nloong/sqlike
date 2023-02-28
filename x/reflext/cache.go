package reflext

import (
	"reflect"
	"runtime"

	"github.com/si3nloong/sqlike/v2/internal/lrucache"
)

var defaultMapper = NewMapperFunc(500, []string{"sqlike", "db", "sql"})

// DefaultMapper : return a default struct mapper
func DefaultMapper() StructMapper {
	return defaultMapper
}

// StructMapper :
type StructMapper interface {
	CodecByType(t reflect.Type) StructInfo
	FieldByName(v reflect.Value, name string) reflect.Value
	FieldByIndexes(v reflect.Value, idxs []int) reflect.Value
	FieldByIndexesReadOnly(v reflect.Value, idxs []int) reflect.Value
	LookUpFieldByName(v reflect.Value, name string) (reflect.Value, bool)
	TraversalsByName(t reflect.Type, names []string) (idxs [][]int)
	TraversalsByNameFunc(t reflect.Type, names []string, fn func(int, []int)) (idxs [][]int)
}

// FormatFunc :
type FormatFunc func(string) string

type mapper struct {
	tags    []string
	cache   lrucache.Cache[reflect.Type, *Struct]
	fmtFunc FormatFunc
}

var _ StructMapper = (*mapper)(nil)

// NewMapperFunc : return an object which complied to `StructMapper` interface
func NewMapperFunc(size int, tags []string, formatter ...FormatFunc) StructMapper {
	fmtFunc := func(v string) string {
		return v
	}
	if len(formatter) > 0 {
		fmtFunc = formatter[0]
	}
	return &mapper{
		cache:   lrucache.New[reflect.Type, *Struct](size),
		tags:    tags,
		fmtFunc: fmtFunc,
	}
}

// CodecByType :
func (m *mapper) CodecByType(t reflect.Type) StructInfo {
	mapping, ok := m.cache.Get(t)
	if !ok {
		// m.mutex.Lock()
		// mapping = getCodec(t, m.tags, m.fmtFunc)
		// _, ok = m.cache[t]
		// if !ok {
		// 	m.cache[t] = mapping
		// }
		// m.mutex.Unlock()
		mapping = getCodec(t, m.tags, m.fmtFunc)
		m.cache.Set(t, mapping)
	}
	return mapping
}

// FieldByName : get reflect.Value from struct by field name
func (m *mapper) FieldByName(v reflect.Value, name string) reflect.Value {
	v = Indirect(v)
	mustBe(v, reflect.Struct)

	tm := m.CodecByType(v.Type())
	fi, ok := tm.LookUpFieldByName(name)
	if !ok {
		panic("field not exists")
	}
	return FieldByIndexes(v, fi.Index())
}

// FieldByIndexes : get reflect.Value from struct by indexes. If the reflect.Value is nil, it will get initialized
func (m *mapper) FieldByIndexes(v reflect.Value, idxs []int) reflect.Value {
	return FieldByIndexes(v, idxs)
}

// FieldByIndexesReadOnly : get reflect.Value from struct by indexes without initialized
func (m *mapper) FieldByIndexesReadOnly(v reflect.Value, idxs []int) reflect.Value {
	return FieldByIndexesReadOnly(v, idxs)
}

// LookUpFieldByName : lookup reflect.Value from struct by field name
func (m *mapper) LookUpFieldByName(v reflect.Value, name string) (reflect.Value, bool) {
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
func (m *mapper) TraversalsByName(t reflect.Type, names []string) (idxs [][]int) {
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
func (m *mapper) TraversalsByNameFunc(t reflect.Type, names []string, fn func(int, []int)) (idxs [][]int) {
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
