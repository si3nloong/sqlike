package primitive

type LockType int

const (
	LockForUpdate LockType = iota + 1
	LockForShare
)

type LockOption int

const (
	NoWait LockOption = iota + 1
	SkipLocked
)

// Lock :
type Lock struct {
	Type   LockType
	Option LockOption
	Of     *Pair
}
