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

	// The struct field name
	FieldName() string

	// Look up option using key
	Option(key string) (val string, exists bool)

	// Look up tag value using key
	LookUp(key string) (val string, exists bool)

	// Get(key string) string
}

// StructTag :
type StructTag struct {
	fieldName string
	name      string
	tag       reflect.StructTag
	opts      map[string]string
}

// Name : returns the name of the struct field
func (st StructTag) Name() string {
	return st.name
}

// FieldName : returns the original field name on struct
func (st StructTag) FieldName() string {
	return st.fieldName
}

func (st StructTag) Option(key string) (val string, exist bool) {
	if st.opts == nil {
		return
	}
	val, exist = st.opts[key]
	return
}

func (st StructTag) Get(key string) string {
	return st.tag.Get(key)
}

func (st StructTag) LookUp(key string) (val string, exist bool) {
	return st.tag.Lookup(key)
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
	// tree representation of the struct
	tree FieldInfo
	// all fields belong to this struct
	fields Fields
	// available properties in sequence
	properties Fields
	// store all field using their index
	indexes map[string]FieldInfo
	// store all field using their name
	names map[string]FieldInfo
}

var _ StructInfo = (*Struct)(nil)

// Fields :
func (s *Struct) Fields() []FieldInfo {
	clone := make(Fields, len(s.fields))
	copy(clone, s.fields)
	return clone
}

// Properties :
func (s *Struct) Properties() []FieldInfo {
	clone := make(Fields, len(s.properties))
	copy(clone, s.properties)
	return clone
}

// LookUpFieldByName : find the field using name
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

			// skip unexported fields (private property)
			if len(f.PkgPath) != 0 && !f.Anonymous {
				continue
			}

			tag := parseTag(f, tagNames, fmtFunc)
			// skip when it's hyphen
			if tag.Name() == "-" {
				continue
			}

			ft := Deref(f.Type)
			sf := &StructField{
				id:       strings.TrimLeft(q.sf.id+"."+strconv.Itoa(i), "."),
				name:     f.Name,
				path:     tag.Name(),
				null:     q.sf.null || IsNullable(f.Type),
				t:        f.Type,
				tag:      tag,
				children: make([]FieldInfo, 0),
				embed:    ft.Kind() == reflect.Struct && f.Anonymous,
			}

			// set parent when it has parent
			if len(q.sf.Index()) > 0 {
				sf.parent = q.sf
			}

			// if tag name is empty, set to field name
			if sf.path == "" && !sf.IsEmbedded() {
				sf.path = f.Name
			}

			if q.pp != "" {
				sf.path = q.pp + "." + sf.path
			}

			q.sf.children = append(q.sf.children, sf)
			sf.idx = appendSlice(q.sf.idx, i)

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
		// if it's an embedded field and name is declare, we should preserve the name
		if sf.Name() != "" && !sf.IsEmbedded() {
			lname = strings.ToLower(sf.Name())
			codec.names[sf.Name()] = sf

			idx := codec.properties.FindIndex(func(fi FieldInfo) bool {
				return strings.ToLower(fi.Name()) == lname
			})
			if idx > -1 {
				// remove item in the slice if the field name is same (overriding embedded struct field)
				codec.properties = append(codec.properties[:idx], codec.properties[idx+1:]...)
			}

			parent := sf.ParentByTraversal(func(f FieldInfo) bool {
				return !f.IsEmbedded()
			})
			if len(sf.Index()) > 1 &&
				sf.Parent() != nil && parent != nil {
				continue
			}

			// not nested embedded struct or embedded struct
			codec.properties = append(codec.properties, sf)
		}
	}

	return codec
}

func appendSlice[T any](s []T, i T) []T {
	x := make([]T, len(s)+1)
	copy(x, s)
	x[len(x)-1] = i
	return x
}

func parseTag(f reflect.StructField, tagNames []string, fmtFunc FormatFunc) (st StructTag) {
	st.fieldName = f.Name
	st.tag = f.Tag
	st.opts = make(map[string]string)

	var (
		name, value string
		parts       []string
		kvs         []string
		ok          bool
	)

	// the latest tag value will override
	for _, tagName := range tagNames {
		value, ok = f.Tag.Lookup(tagName)
		if !ok {
			continue
		}

		parts = strings.Split(value, ",")
		if fname := strings.TrimSpace(parts[0]); fname != "" {
			name = fname
		}

		for _, opt := range parts[1:] {
			opt = strings.TrimSpace(opt)
			kvs = strings.SplitN(opt, "=", 2)
			if len(kvs) >= 2 {
				k := strings.TrimSpace(strings.ToLower(kvs[0]))
				st.opts[k] = strings.TrimSpace(kvs[1])
				continue
			}
			opt = strings.ToLower(opt)
			st.opts[opt] = ""
		}
	}
	if fmtFunc != nil {
		name = fmtFunc(name)
	}
	st.name = name
	return
}
