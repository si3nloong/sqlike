package sql

// Column :
type Column struct {
	Name         string
	DataType     string
	Type         string
	Size         int
	Nullable     bool
	DefaultValue *string
	Charset      *string
	Collation    *string
	Extra        string
}
