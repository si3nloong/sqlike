package columns

// DataType :
type DataType int

// Column :
type Column struct {
	Name         string
	DataType     string
	Type         string
	Size         int
	Nullable     bool
	DefaultValue *string
	CharSet      *string
	Collation    *string
	Extra        string
}
