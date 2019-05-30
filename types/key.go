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

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sqlike/sql/component"
	"github.com/si3nloong/sqlike/util"
	"golang.org/x/xerrors"
)

// Writer :
type writer interface {
	io.Writer
	WriteString(string) (int, error)
	WriteByte(byte) error
}

var (
	latin1    = `latin1`
	latin1Bin = `latin1_bin`
)

// Key :
type Key struct {
	Namespace string
	Kind      string
	ID        int64
	Name      string
	Parent    *Key
}

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

// Scan :
func (k *Key) Scan(it interface{}) error {
	switch vi := it.(type) {
	case []byte:
		if err := k.unmarshal(string(vi)); err != nil {
			return err
		}
	case string:
		if err := k.unmarshal(vi); err != nil {
			return err
		}
	}
	return nil
}

// Incomplete :
func (k *Key) Incomplete() bool {
	return k.Name == "" && k.ID == 0
}

func (k *Key) unmarshal(str string) error {
	var (
		idx   int
		path  string
		paths []string
		value string
		err   error
	)
	for {
		idx = strings.LastIndex(str, "/")
		if idx < 0 {
			break
		}
		path = str[idx:]
		paths = strings.Split(path, ",")
		if len(paths) != 2 {
			return xerrors.New("invalid key path")
		}
		k.Kind = paths[0]
		value = paths[1]
		if len(value) < 1 {
			return xerrors.New("invalid key string")
		}
		if value[0] == '\'' {
			value = strings.Trim(value, "'")
			value, _ = url.PathUnescape(value)
			k.Name = value
		} else {
			k.ID, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
		}
		str = str[:idx]
		if len(str) < 1 {
			break
		}
		if k.Parent == nil {
			k.Parent = new(Key)
		}
		k = k.Parent
	}
	return nil
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
func (k *Key) marshal(w writer, escape bool) {
	if k.Parent != nil {
		k.Parent.marshal(w, escape)
		w.WriteByte('/')
	}
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

// MarshalText :
func (k *Key) MarshalText() ([]byte, error) {
	buf := new(bytes.Buffer)
	k.marshal(buf, true)
	return buf.Bytes(), nil
}

type gobKey struct {
	Kind      string
	StringID  string
	IntID     int64
	Parent    *gobKey
	AppID     string
	Namespace string
}

func keyToGobKey(k *Key) *gobKey {
	if k == nil {
		return nil
	}
	return &gobKey{
		Kind:      k.Kind,
		StringID:  k.Name,
		IntID:     k.ID,
		Parent:    keyToGobKey(k.Parent),
		Namespace: k.Namespace,
	}
}

func gobKeyToKey(gk *gobKey) *Key {
	if gk == nil {
		return nil
	}
	return &Key{
		Kind:      gk.Kind,
		Name:      gk.StringID,
		ID:        gk.IntID,
		Parent:    gobKeyToKey(gk.Parent),
		Namespace: gk.Namespace,
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
	if err := gob.NewEncoder(buf).Encode(keyToGobKey(k)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode unmarshals a sequence of bytes using an encoding/gob.Decoder.
func (k *Key) GobDecode(buf []byte) error {
	gk := new(gobKey)
	if err := gob.NewDecoder(bytes.NewBuffer(buf)).Decode(gk); err != nil {
		return err
	}
	*k = *gobKeyToKey(gk)
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
