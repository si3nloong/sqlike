package casbin

// Policy :
type Policy struct {
	PType string `sqlike:",size:3"`
	V0    string
	V1    string
	V2    string `sqlike:",size:50"`
	V3    string `sqlike:",size:50"`
	V4    string `sqlike:",size:50"`
	V5    string `sqlike:",size:50"`
}
