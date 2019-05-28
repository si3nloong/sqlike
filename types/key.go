package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/base64"
	"encoding/gob"
	"io"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/SianLoong/sqlike/reflext"
	"bitbucket.org/SianLoong/sqlike/sqlike/sql/component"
	"bitbucket.org/SianLoong/sqlike/util"
)

// Writer :
type Writer interface {
	io.Writer
	WriteString(string) (int, error)
	WriteByte(byte) error
}

// Key :
type Key struct {
	Namespace string
	Kind      string
	ID        int64
	Name      string
	Parent    *Key
}

var (
	latin1    = `latin1`
	latin1Bin = `latin1_bin`
)

// DataType :
func (k *Key) DataType(driver string, sf *reflext.StructField) component.Column {
	return component.Column{
		Name:      sf.Path,
		DataType:  `VARCHAR`,
		Type:      `VARCHAR(512)`,
		Nullable:  false,
		CharSet:   &latin1,
		Collation: &latin1Bin,
	}
}

// Value :
func (k *Key) Value() (driver.Value, error) {
	blr := util.AcquireString()
	defer util.ReleaseString(blr)
	k.marshal(blr, true)
	return blr.String(), nil
}

// Incomplete :
func (k *Key) Incomplete() bool {
	return k.Name == "" && k.ID == 0
}

// valid returns whether the key is valid.
func (k *Key) valid() bool {
	if k == nil {
		return false
	}
	for ; k != nil; k = k.Parent {
		if k.Kind == "" {
			return false
		}
		if k.Name != "" && k.ID != 0 {
			return false
		}
		if k.Parent != nil {
			if k.Parent.Incomplete() {
				return false
			}
			if k.Parent.Namespace != k.Namespace {
				return false
			}
		}
	}
	return true
}

// Equal reports whether two keys are equal. Two keys are equal if they are
// both nil, or if their kinds, IDs, names, namespaces and parents are equal.
func (k *Key) Equal(o *Key) bool {
	for {
		if k == nil || o == nil {
			return k == o // if either is nil, both must be nil
		}
		if k.Namespace != o.Namespace || k.Name != o.Name || k.ID != o.ID || k.Kind != o.Kind {
			return false
		}
		if k.Parent == nil && o.Parent == nil {
			return true
		}
		k = k.Parent
		o = o.Parent
	}
}

// marshal marshals the key's string representation to the buffer.
func (k *Key) marshal(w Writer, escape bool) {
	if k.Parent != nil {
		k.Parent.marshal(w, escape)
	}
	w.WriteByte('/')
	w.WriteString(k.Kind)
	w.WriteByte(',')
	if k.Name != "" {
		w.WriteByte('\'')
		w.WriteString(url.PathEscape(k.Name))
		w.WriteByte('\'')
	} else {
		w.WriteString(strconv.FormatInt(k.ID, 10))
	}
}

// Encode returns an opaque representation of the key
// suitable for use in HTML and URLs.
// This is compatible with the Python and Java runtimes.
func (k *Key) Encode() string {
	b, err := k.GobEncode()
	if err != nil {
		panic(err)
	}
	// Trailing padding is stripped.
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}

// DecodeKey decodes a key from the opaque representation returned by Encode.
func DecodeKey(encoded string) (*Key, error) {
	// Re-add padding.
	if m := len(encoded) % 4; m != 0 {
		encoded += strings.Repeat("=", 4-m)
	}

	b, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	k := new(Key)
	if err := k.GobDecode(b); err != nil {
		return nil, err
	}
	return k, nil
}

// String returns a string representation of the key.
func (k *Key) String() string {
	if k == nil {
		return ""
	}
	b := bytes.NewBuffer(make([]byte, 0, 512))
	k.marshal(b, false)
	return b.String()
}

// GobEncode marshals the key into a sequence of bytes
// using an encoding/gob.Encoder.
func (k *Key) GobEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(k); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode unmarshals a sequence of bytes using an encoding/gob.Decoder.
func (k *Key) GobDecode(buf []byte) error {
	if err := gob.NewDecoder(bytes.NewBuffer(buf)).Decode(k); err != nil {
		return err
	}
	return nil
}

// NameKey creates a new key with a name.
// The supplied kind cannot be empty.
// The supplied parent must either be a complete key or nil.
// The namespace of the new key is empty.
func NameKey(kind, name string, parent *Key) *Key {
	return &Key{
		Kind:   kind,
		Name:   name,
		Parent: parent,
	}
}

// IDKey creates a new key with an ID.
// The supplied kind cannot be empty.
// The supplied parent must either be a complete key or nil.
// The namespace of the new key is empty.
func IDKey(kind string, id int64, parent *Key) *Key {
	return &Key{
		Kind:   kind,
		ID:     id,
		Parent: parent,
	}
}

const (
	minSeed = int64(100000000)
	maxSeed = int64(999999999)
)

// NewIDKey :
func NewIDKey(kind string, parent *Key) *Key {
	rand.Seed(time.Now().UnixNano())
	strID := strconv.FormatInt(time.Now().Unix(), 10) + strconv.FormatInt(rand.Int63n(maxSeed-minSeed)+minSeed, 10)
	id, _ := strconv.ParseInt(strID, 10, 64)

	return &Key{
		Kind:   kind,
		ID:     id,
		Parent: parent,
	}
}
