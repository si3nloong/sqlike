package reflext

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// StructInfo :
type StructInfo interface {
	Fields() []FieldInfo
	Properties() []FieldInfo
	LookUpFieldByName(name string) (FieldInfo, bool)
	GetByTraversal(index []int) FieldInfo
}

// FieldInfo :
type FieldInfo interface {
	// New name of struct field
	Name() string

	// reflect.Type of the field
	Type() reflect.Type

	// index position of the field
	Index() []int

	// get struct tag info
	Tag() FieldTag

	// if the field is struct, parent will not be nil
	// this will be the parent struct of current struct
	Parent() FieldInfo

	ParentByTraversal(cb func(FieldInfo) bool) FieldInfo

	// if the field is struct, children will not be nil
	// this will be the fields of current struct
	Children() []FieldInfo

	// determine the field is nullable
	IsNullable() bool

	// determine the field is embedded struct
	IsEmbedded() bool
}

// FieldTag :
type FieldTag interface {
	Name() string
	FieldName() string

	// look up tag value using key
	LookUp(key string) (val string, exists bool)

	Get(key string) string
}

// StructTag :
type StructTag struct {
	fieldName string
	name      string
	opts      map[string]string
}

// Name :
func (st StructTag) Name() string {
	return st.name
}

// FieldName :
func (st StructTag) FieldName() string {
	return st.fieldName
}

// Get :
func (st StructTag) Get(key string) string {
	if st.opts == nil {
		return ""
	}
	return st.opts[key]
}

// LookUp :
func (st StructTag) LookUp(key string) (val string, exist bool) {
	if st.opts == nil {
		return
	}
	val, exist = st.opts[key]
	return
}

// StructField :
type StructField struct {
	id       string
	idx      []int
	name     string
	path     string
	t        reflect.Type
	null     bool
	tag      StructTag
	embed    bool
	parent   FieldInfo
	children []FieldInfo
}

var _ FieldInfo = (*StructField)(nil)

// Name :
func (sf *StructField) Name() string {
	return sf.path
}

// Type :
func (sf *StructField) Type() reflect.Type {
	return sf.t
}

// Tag :
func (sf *StructField) Tag() FieldTag {
	return sf.tag
}

// Index :
func (sf *StructField) Index() []int {
	return sf.idx
}

// Parent :
func (sf *StructField) Parent() FieldInfo {
	return sf.parent
}

// Children :
func (sf *StructField) Children() []FieldInfo {
	return sf.children
}

// IsNullable :
func (sf *StructField) IsNullable() bool {
	return sf.null
}

// IsEmbedded :
func (sf *StructField) IsEmbedded() bool {
	return sf.embed
}

// ParentByTraversal :
func (sf *StructField) ParentByTraversal(cb func(FieldInfo) bool) FieldInfo {
	prnt := sf.parent
	for prnt != nil {
		if cb(prnt) {
			break
		}
		prnt = prnt.Parent()
	}
	return prnt
}

// Struct :
type Struct struct {
	tree       FieldInfo
	fields     Fields // all fields belong to this struct
	properties Fields // available properties in sequence
	indexes    map[string]FieldInfo
	names      map[string]FieldInfo
}

var _ StructInfo = (*Struct)(nil)

// Fields :
func (s *Struct) Fields() []FieldInfo {
	return append(make(Fields, 0, len(s.fields)), s.fields...)
}

// Properties :
func (s *Struct) Properties() []FieldInfo {
	return append(make(Fields, 0, len(s.properties)), s.properties...)
}

// LookUpFieldByName :
func (s *Struct) LookUpFieldByName(name string) (FieldInfo, bool) {
	x, ok := s.names[name]
	return x, ok
}

// GetByTraversal :
func (s *Struct) GetByTraversal(index []int) FieldInfo {
	if len(index) == 0 {
		return nil
	}

	tree := s.tree
	for _, i := range index {
		children := tree.Children()
		if i >= len(children) || children[i] == nil {
			return nil
		}
		tree = children[i]
	}
	return tree
}

// Fields :
type Fields []FieldInfo

func (x Fields) FindIndex(cb func(f FieldInfo) bool) int {
	for idx, f := range x {
		if cb(f) {
			return idx
		}
	}
	return -1
}

func (x Fields) Len() int { return len(x) }

func (x Fields) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x Fields) Less(i, j int) bool {
	for k, xik := range x[i].Index() {
		if k >= len(x[j].Index()) {
			return false
		}
		if xik != x[j].Index()[k] {
			return xik < x[j].Index()[k]
		}
	}
	return len(x[i].Index()) < len(x[j].Index())
}

type typeQueue struct {
	t  reflect.Type
	sf *StructField
	pp string // parent path
}

func getCodec(t reflect.Type, tagNames []string, fmtFunc FormatFunc) *Struct {
	fields := make([]FieldInfo, 0)

	root := &StructField{}
	queue := []typeQueue{}
	queue = append(queue, typeQueue{Deref(t), root, ""})

	for len(queue) > 0 {
		q := queue[0]
		q.sf.children = make([]FieldInfo, 0)

		for i := 0; i < q.t.NumField(); i++ {
			f := q.t.Field(i)

			// skip unexported fields
			if len(f.PkgPath) != 0 && !f.Anonymous {
				continue
			}

			tag := parseTag(f, tagNames, fmtFunc)
			if tag.name == "-" {
				continue
			}

			sf := &StructField{
				id:       strings.TrimLeft(q.sf.id+"."+strconv.Itoa(i), "."),
				name:     f.Name,
				path:     tag.name,
				null:     q.sf.null || IsNullable(f.Type),
				t:        f.Type,
				tag:      tag,
				children: make([]FieldInfo, 0),
			}

			if len(q.sf.Index()) > 0 {
				sf.parent = q.sf
			}

			if sf.path == "" {
				sf.path = sf.tag.name
			}

			if q.pp != "" {
				sf.path = q.pp + "." + sf.path
			}

			ft := Deref(f.Type)
			q.sf.children = append(q.sf.children, sf)
			sf.idx = appendSlice(q.sf.idx, i)
			sf.embed = ft.Kind() == reflect.Struct && f.Anonymous

			if ft.Kind() == reflect.Struct {
				// check recursive, prevent infinite loop
				if q.t == ft {
					goto nextStep
				}

				// embedded struct
				path := sf.path
				if f.Anonymous {
					if sf.tag.FieldName() == "" {
						path = q.pp
					}
					// queue = append(queue, typeQueue{ft, sf, path})
				}

				queue = append(queue, typeQueue{ft, sf, path})
			}

		nextStep:
			fields = append(fields, sf)
		}

		queue = queue[1:]
	}

	codec := &Struct{
		tree:       root,
		fields:     fields,
		properties: make([]FieldInfo, 0, len(fields)),
		indexes:    make(map[string]FieldInfo),
		names:      make(map[string]FieldInfo),
	}

	lname := ""
	sort.Sort(codec.fields)

	for _, sf := range codec.fields {
		codec.indexes[sf.(*StructField).id] = sf
		if sf.Name() != "" && !sf.IsEmbedded() {
			lname = strings.ToLower(sf.Name())
			codec.names[sf.Name()] = sf

			idx := codec.properties.FindIndex(func(each FieldInfo) bool {
				return strings.ToLower(each.Tag().Name()) == lname
			})
			if idx > -1 {
				// remove item in the slice if the field name is same (overriding embedded struct field)
				codec.properties = append(codec.properties[:idx], codec.properties[idx+1:]...)
			}

			prnt := sf.ParentByTraversal(func(f FieldInfo) bool {
				return !f.IsEmbedded()
			})
			if len(sf.Index()) > 1 &&
				sf.Parent() != nil && prnt != nil {
				continue
			}

			// not nested embedded struct or embedded struct
			codec.properties = append(codec.properties, sf)
		}
	}

	return codec
}

func appendSlice(s []int, i int) []int {
	x := make([]int, len(s)+1)
	copy(x, s)
	x[len(x)-1] = i
	return x
}

func parseTag(f reflect.StructField, tagNames []string, fmtFunc FormatFunc) (st StructTag) {
	parts := strings.Split(f.Tag.Get(tagNames[0]), ",")
	name := strings.TrimSpace(parts[0])
	st.fieldName = name
	if name == "" {
		name = f.Name
		if fmtFunc != nil {
			name = fmtFunc(name)
		}
	}
	st.name = name
	st.opts = make(map[string]string)
	// for _, tagName := range tagNames {
	// 	parts := strings.Split(f.Tag.Get(tagName), ",")
	// 	name := strings.TrimSpace(parts[0])
	// 	if name != "" {
	// 		if fmtFunc != nil {
	// 			name = fmtFunc(name)
	// 		}
	// 		st.name = name
	// 	}
	if len(parts) > 1 {
		for _, opt := range parts[1:] {
			opt = strings.TrimSpace(opt)
			if strings.Contains(opt, "=") {
				kv := strings.SplitN(opt, "=", 2)
				k := strings.TrimSpace(strings.ToLower(kv[0]))
				st.opts[k] = strings.TrimSpace(kv[1])
				continue
			}
			opt = strings.ToLower(opt)
			st.opts[opt] = ""
		}
	}
	// }
	return
}
