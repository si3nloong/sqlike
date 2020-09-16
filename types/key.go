package types

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	pb "github.com/si3nloong/sqlike/proto"
	"github.com/si3nloong/sqlike/reflext"
	sqldriver "github.com/si3nloong/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sqlike/columns"
	"github.com/si3nloong/sqlike/util"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// Writer :
type writer interface {
	io.Writer
	WriteString(string) (int, error)
	WriteByte(byte) error
}

var (
	latin1    = "latin1"
	latin1Bin = "latin1_bin"
)

// Key :
type Key struct {
	Namespace string
	Kind      string
	IntID     int64
	NameID    string
	Parent    *Key
}

var (
	_ driver.Valuer          = (*Key)(nil)
	_ sql.Scanner            = (*Key)(nil)
	_ fmt.Stringer           = (*Key)(nil)
	_ encoding.TextMarshaler = (*Key)(nil)
	_ json.Marshaler         = (*Key)(nil)
)

// DataType :
func (k Key) DataType(t sqldriver.Info, sf reflext.StructFielder) columns.Column {
	return columns.Column{
		Name:      sf.Name(),
		DataType:  "VARCHAR",
		Type:      "VARCHAR(512)",
		Nullable:  reflext.IsNullable(sf.Type()),
		Charset:   &latin1,
		Collation: &latin1Bin,
	}
}

// ID :
func (k Key) ID() string {
	if k.NameID != "" {
		return k.NameID
	}
	return strconv.FormatInt(k.IntID, 10)
}

// Root :
func (k *Key) Root() *Key {
	for {
		if k.Parent == nil {
			return k
		}
		k = k.Parent
	}
}

// Clone :
func (k *Key) Clone() *Key {
	return copyKey(k)
}

func copyKey(k *Key) *Key {
	if k == nil {
		return nil
	}
	nk := new(Key)
	nk.Namespace = k.Namespace
	nk.Kind = k.Kind
	nk.NameID = k.NameID
	nk.IntID = k.IntID
	nk.Parent = copyKey(k.Parent)
	return nk
}

// Value :
func (k Key) Value() (driver.Value, error) {
	blr := util.AcquireString()
	defer util.ReleaseString(blr)
	marshal(&k, blr, true)
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

// Incomplete : is a safe method to check key is nil or empty
func (k *Key) Incomplete() bool {
	if k == nil {
		return true
	}
	return k.NameID == "" && k.IntID == 0
}

// // valid returns whether the key is valid.
// func (k *Key) IsZero() bool {
// 	if k == nil {
// 		return false
// 	}
// 	for ; k != nil; k = k.Parent {
// 		if k.Kind == "" {
// 			return false
// 		}
// 		if k.NameID != "" && k.IntID != 0 {
// 			return false
// 		}
// 		if k.Parent != nil {
// 			if k.Parent.Incomplete() {
// 				return false
// 			}
// 			if k.Parent.Namespace != k.Namespace {
// 				return false
// 			}
// 		}
// 	}
// 	return true
// }

// Equal reports whether two keys are equal. Two keys are equal if they are
// both nil, or if their kinds, IDs, names, namespaces and parents are equal.
func (k *Key) Equal(o *Key) bool {
	for {
		if k == nil || o == nil {
			return k == o // if either is nil, both must be nil
		}
		if k.Namespace != o.Namespace || k.NameID != o.NameID || k.IntID != o.IntID || k.Kind != o.Kind {
			return false
		}
		if k.Parent == nil && o.Parent == nil {
			return true
		}
		k = k.Parent
		o = o.Parent
	}
}

// MarshalText :
func (k Key) MarshalText() ([]byte, error) {
	buf := new(bytes.Buffer)
	marshal(&k, buf, true)
	return buf.Bytes(), nil
}

// MarshalBinary :
func (k Key) MarshalBinary() ([]byte, error) {
	if k.Incomplete() {
		return []byte(`null`), nil
	}
	return []byte(`"` + k.Encode() + `"`), nil
}

// MarshalJSON :
func (k Key) MarshalJSON() ([]byte, error) {
	if k.Incomplete() {
		return []byte(`null`), nil
	}
	return []byte(`"` + k.Encode() + `"`), nil
}

// MarshalJSONB :
func (k Key) MarshalJSONB() ([]byte, error) {
	if k.Incomplete() {
		return []byte(`null`), nil
	}
	buf := new(bytes.Buffer)
	buf.WriteRune('"')
	marshal(&k, buf, true)
	buf.WriteRune('"')
	return buf.Bytes(), nil
}

// UnmarshalBinary :
func (k *Key) UnmarshalBinary(b []byte) error {
	str := string(b)
	if str == "null" {
		return nil
	}
	key, err := DecodeKey(str)
	if err != nil {
		return err
	}
	*k = *key
	return nil
}

// UnmarshalText :
func (k *Key) UnmarshalText(b []byte) error {
	str := string(b)
	if str == "null" {
		return nil
	}
	key, err := DecodeKey(str)
	if err != nil {
		return err
	}
	*k = *key
	return nil
}

// UnmarshalJSON :
func (k *Key) UnmarshalJSON(b []byte) error {
	length := len(b)
	if length < 2 {
		return errors.New("types: invalid key json value")
	}
	str := string(b)
	if str == "null" {
		return nil
	}
	str = string(b[1 : length-1])
	key, err := DecodeKey(str)
	if err != nil {
		return err
	}
	*k = *key
	return nil
}

// UnmarshalJSONB :
func (k *Key) UnmarshalJSONB(b []byte) error {
	length := len(b)
	if length < 2 {
		return errors.New("types: invalid key json value")
	}
	str := string(b)
	if str == "null" {
		return nil
	}
	str = string(b[1 : length-1])
	return k.unmarshal(str)
}

// MarshalBSONValue :
func (k Key) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if k.Incomplete() {
		return bsontype.Null, nil, nil
	}
	return bsontype.String, bsoncore.AppendString(nil, k.String()), nil
}

// UnmarshalBSONValue :
func (k *Key) UnmarshalBSONValue(t bsontype.Type, b []byte) error {
	if k == nil {
		return errors.New("types: invalid key value <nil>")
	}
	v, _, ok := bsoncore.ReadString(b)
	if !ok {
		return errors.New("types: invalid bson string value")
	}
	return k.unmarshal(v)
}

// String returns a string representation of the key.
func (k Key) String() string {
	b := bytes.NewBuffer(make([]byte, 0, 512))
	marshal(&k, b, false)
	return b.String()
}

// marshal marshals the key's string representation to the buffer.
func marshal(k *Key, w writer, escape bool) {
	if k.Parent != nil {
		marshal(k.Parent, w, escape)
		w.WriteByte('/')
	}
	w.WriteString(k.Kind)
	w.WriteByte(',')
	if k.NameID != "" {
		w.WriteByte('\'')
		w.WriteString(url.PathEscape(k.NameID))
		w.WriteByte('\'')
	} else {
		w.WriteString(strconv.FormatInt(k.IntID, 10))
	}
}

func (k *Key) unmarshal(str string) error {
	if str == "null" {
		k = nil
		return nil
	}

	var (
		idx    int
		path   string
		length int
		paths  []string
		value  string
		err    error
	)

	for {
		idx = strings.LastIndex(str, "/")
		path = str
		if idx > -1 {
			path = str[idx+1:]
		}
		paths = strings.Split(path, ",")
		if len(paths) != 2 {
			return errors.New("invalid key path")
		}
		k.Kind = paths[0]
		value = paths[1]
		length = len(value)
		if length < 1 {
			return errors.New("invalid key string")
		}
		if length > 2 && value[0] == '\'' && value[length-1] == '\'' {
			value = value[1 : length-1]
			value, _ = url.PathUnescape(value)
			k.NameID = value
		} else {
			k.IntID, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
		}

		if idx > -1 {
			str = str[:idx]
			if len(str) < 1 {
				return nil
			}
		} else {
			return nil
		}

		if k.Parent == nil {
			k.Parent = new(Key)
		}
		k = k.Parent
	}
}

type gobKey struct {
	Namespace string
	AppID     string
	Kind      string
	StringID  string
	IntID     int64
	Parent    *gobKey
}

func keyToGobKey(k *Key) *gobKey {
	if k == nil {
		return nil
	}
	return &gobKey{
		Kind:      k.Kind,
		StringID:  k.NameID,
		IntID:     k.IntID,
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
		IntID:     gk.IntID,
		NameID:    gk.StringID,
		Parent:    gobKeyToKey(gk.Parent),
		Namespace: gk.Namespace,
	}
}

// Encode returns an opaque representation of the key
// suitable for use in HTML and URLs.
// This is compatible with the Python and Java runtimes.
func (k Key) Encode() string {
	pk := keyToProto(&k)
	b, err := proto.Marshal(pk)
	if err != nil {
		panic(err)
	}
	// Trailing padding is stripped.
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}

// ParseKey :
func ParseKey(value string) (*Key, error) {
	k := new(Key)
	if err := k.unmarshal(value); err != nil {
		return nil, err
	}
	return k, nil
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

	k := new(pb.Key)
	if err := proto.Unmarshal(b, k); err != nil {
		return nil, err
	}

	return protoToKey(k), nil
}

// GobEncode marshals the key into a sequence of bytes
// using an encoding/gob.Encoder.
func (k Key) GobEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(keyToGobKey(&k)); err != nil {
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

func keyToProto(k *Key) *pb.Key {
	if k == nil {
		return nil
	}

	return &pb.Key{
		Namespace: k.Namespace,
		Kind:      k.Kind,
		NameID:    k.NameID,
		IntID:     k.IntID,
		Parent:    keyToProto(k.Parent),
	}
}

func protoToKey(pk *pb.Key) *Key {
	if pk == nil {
		return nil
	}

	return &Key{
		Namespace: pk.Namespace,
		Kind:      pk.Kind,
		NameID:    pk.NameID,
		IntID:     pk.IntID,
		Parent:    protoToKey(pk.Parent),
	}
}

// NameKey creates a new key with a name.
// The supplied kind cannot be empty.
// The supplied parent must either be a complete key or nil.
// The namespace of the new key is empty.
func NameKey(kind, name string, parent *Key) *Key {
	return &Key{
		Kind:   kind,
		NameID: name,
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
		IntID:  id,
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
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		panic(err)
	}

	return &Key{
		Kind:   kind,
		IntID:  id,
		Parent: parent,
	}
}

// NewNameKey :
func NewNameKey(kind string, parent *Key) *Key {
	return &Key{
		Kind:   kind,
		NameID: uuid.New().String(),
		Parent: parent,
	}
}
