package options

type insertMode int

// insert modes :
const (
	InsertIgnore insertMode = iota + 1
	InsertOnDuplicate
)

// InsertOptions :
type InsertOptions struct {
	Mode    insertMode
	IsDebug bool
}
