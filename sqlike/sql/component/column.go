package component

// Column :
type Column struct {
	Name         string
	DataType     string
	Type         string
	Nullable     bool
	DefaultValue *string
	CharSet      *string
	Collation    *string
	Extra        string
}
