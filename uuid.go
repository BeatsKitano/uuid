// Package uuid wraps github.com/google/uuid as a protobuf type with
// implementations of various serialization interfaces (JSON, BSON, GraphQL).
package uuid

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/google/uuid"
	"github.com/prometheus/common/log"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"google.golang.org/protobuf/proto"
)

// Ensure UUID implements proto.Message.
var _ proto.Message = (*UUID)(nil)

// Parse parses a UUID from its string representation.
func Parse(s string) (*UUID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return nil, err
	}
	return &UUID{Value: parsed[:]}, nil
}

// New generates a new random (v4) UUID.
func New() *UUID {
	id, err := uuid.NewRandomFromReader(rand.Reader)
	if err != nil {
		panic(fmt.Sprintf("uuid: failed to generate random UUID: %v", err))
	}
	return &UUID{Value: id[:]}
}

// ToUUID converts the protobuf UUID to a google/uuid.UUID.
func (x *UUID) ToUUID() (uuid.UUID, error) {
	if x == nil {
		return uuid.Nil, nil
	}
	return uuid.FromBytes(x.Value)
}

// MarshalJSON implements the json.Marshaler interface.
func (x *UUID) MarshalJSON() ([]byte, error) {
	if x == nil {
		return json.Marshal(nil)
	}
	id, err := x.ToUUID()
	if err != nil {
		return nil, err
	}
	return json.Marshal(id.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (x *UUID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := uuid.Parse(s)
	if err != nil {
		return err
	}
	x.Value = parsed[:]
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface.
func (x *UUID) MarshalGQL(w io.Writer) {
	data, err := x.MarshalJSON()
	if err != nil {
		log.Errorf("Error marshalling %v to GraphQL: %s", x, err)
		return
	}
	if _, err := w.Write(data); err != nil {
		log.Errorf("Error writing %v to GraphQL writer: %s", x, err)
	}
}

// UnmarshalGQL implements the graphql.Unmarshaler interface.
func (x *UUID) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("value for unmarshalling was not a string: %v", v)
	}
	return x.UnmarshalJSON([]byte(strconv.Quote(s)))
}

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (x *UUID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if x == nil {
		return bsontype.Null, nil, nil
	}
	val := bsonx.Binary(bsontype.BinaryUUID, x.Value)
	return val.MarshalBSONValue()
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface.
func (x *UUID) UnmarshalBSONValue(bsonType bsontype.Type, data []byte) error {
	if bsonType != bsontype.Binary || len(data) < 5 || data[0] != 0x10 || data[4] != bsontype.BinaryUUID {
		return fmt.Errorf("could not unmarshal %v as a UUID", bsonType)
	}
	x.Value = make([]byte, 16)
	copy(x.Value, data[5:21])
	return nil
}
