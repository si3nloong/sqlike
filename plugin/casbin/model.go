package casbin

// PermissionRule :
type PermissionRule struct {
	PType string `sqlike:",index"`
	V0    string `sqlike:",index"`
	V1    string `sqlike:",index"`
	V2    string `sqlike:",index"`
	V3    string `sqlike:",index"`
	V4    string `sqlike:",index"`
	V5    string `sqlike:",index"`
}
