package options

// LockMode :
type LockMode int

// Locking :
const (
	NoLock LockMode = iota
	LockForUpdate
	LockForShare
)
