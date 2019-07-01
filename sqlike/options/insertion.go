package options

type insertMode int

// insert modes :
const (
	InsertIgnore insertMode = iota + 1
	InsertOnDuplicate
)

// InsertOptions :
type InsertOptions struct {
	Mode  insertMode
	Omits []string
	Debug bool
}

// SetMode :
func (opt *InsertOptions) SetMode(mode insertMode) *InsertOptions {
	opt.Mode = mode
	return opt
}

// SetDebug :
func (opt *InsertOptions) SetDebug(debug bool) *InsertOptions {
	opt.Debug = debug
	return opt
}

// SetOmitFields :
func (opt *InsertOptions) SetOmitFields(fields ...string) *InsertOptions {
	opt.Omits = fields
	return opt
}

// AppendOmitField :
func (opt *InsertOptions) AppendOmitField(field string) *InsertOptions {
	opt.Omits = append(opt.Omits, field)
	return opt
}

// SetOnConflict :
func (opt *InsertOptions) SetOnConflict(src []interface{}) *InsertOptions {
	return opt
}

// AppendOnConflict :
func (opt *InsertOptions) AppendOnConflict(col string, val interface{}) *InsertOptions {
	return opt
}
