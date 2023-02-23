package primitive

type LockType int

const (
	NoLock LockType = iota
	LockForUpdate
	LockForShare
)

type LockOption int

const (
	NoLockOption LockOption = iota
	NoWait
	SkipLocked
)

// FIXME:
type Lock struct {
	Type   LockType
	Option LockOption
}

func (l *Lock) NoWait() *Lock {
	l.Option = NoWait
	return l
}

func (l *Lock) SkipLocked() *Lock {
	l.Option = SkipLocked
	return l
}
