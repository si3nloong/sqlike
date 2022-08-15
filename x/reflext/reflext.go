package reflext

import (
	"log"
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
	// name based on struct tag value
	Name() string

	// original field name on struct
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

// FieldName : returns the original field name on struct
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
	// struct type to inspect
	t  reflect.Type
	sf *StructField
	// parent path
	pp string
}

func getCodec(t reflect.Type, tagNames []string, fmtFunc FormatFunc) *Struct {
	fields := make([]FieldInfo, 0)

	root := &StructField{}
	queue := []typeQueue{}
	queue = append(queue, typeQueue{Deref(t), root, ""})

	for len(queue) > 0 {
		// Pop the first item
		q := queue[0]
		q.sf.children = make([]FieldInfo, 0)

		noOfField := q.t.NumField()
		for i := 0; i < noOfField; i++ {
			f := q.t.Field(i)

			// Skip unexported fields
			if len(f.PkgPath) != 0 && !f.Anonymous {
				continue
			}

			tag := parseTag(f, tagNames, fmtFunc)
			// If the tag value is "-", we skip the field
			if tag.name == "-" {
				continue
			}

			sf := &StructField{
				id:   strings.TrimLeft(q.sf.id+"."+strconv.Itoa(i), "."),
				path: tag.name,
				null: q.sf.null || IsNullable(f.Type),
				t:    f.Type,
				tag:  tag,
			}

			if len(q.sf.Index()) > 0 {
				sf.parent = q.sf
			}

			if q.pp != "" {
				sf.path = q.pp + "." + sf.path
			}

			ft := Deref(f.Type)
			q.sf.children = append(q.sf.children, sf)
			sf.idx = appendSlice(q.sf.idx, i)
			sf.embed = ft.Kind() == reflect.Struct && f.Anonymous

			// If the struct property is a `struct`
			if ft.Kind() == reflect.Struct {
				// Check recursive, prevent infinite loop
				if q.t == ft {
					goto nextField
				}

				// If the struct is embedded struct
				path := sf.path
				if f.Anonymous {
					// If the struct poperty is embedded and the tag value is empty,
					// mean it will respect to parent path name
					if sf.tag.name == "" {
						log.Println("embedded is empty")
						path = q.pp
						// queue = append(queue, typeQueue{ft, sf, q.pp})
					}
					log.Println("FieldName =>", q.pp, "Name =>", sf.tag.Name(), "Path =>", path)
					// queue = append(queue, typeQueue{ft, sf, q.pp})
					// goto nextStep
				}

				queue = append(queue, typeQueue{ft, sf, path})
			}

		nextField:
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

func appendSlice[T any](s []T, i T) []T {
	x := make([]T, len(s)+1)
	copy(x, s)
	x[len(x)-1] = i
	return x
}

func parseTag(f reflect.StructField, tagNames []string, fmtFunc FormatFunc) (tag StructTag) {
	tag.fieldName = f.Name
	tag.opts = make(map[string]string)

	for _, tagName := range tagNames {
		name, exists := f.Tag.Lookup(tagName)
		if !exists && f.Anonymous {
			continue
		}

		values := strings.Split(name, ",")
		name = strings.TrimSpace(values[0])
		if name == "" {
			name = f.Name
			if fmtFunc != nil {
				name = fmtFunc(name)
			}
		}

		tag.name = name
		if len(values) > 1 {
			for _, opt := range values[1:] {
				opt = strings.TrimSpace(opt)
				if strings.Contains(opt, "=") {
					kv := strings.SplitN(opt, "=", 2)
					k := strings.TrimSpace(strings.ToLower(kv[0]))
					tag.opts[k] = strings.TrimSpace(kv[1])
					continue
				}
				opt = strings.ToLower(opt)
				tag.opts[opt] = ""
			}
		}
	}
	return
}
