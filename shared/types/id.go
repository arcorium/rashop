package types

import (
  "encoding/hex"
  "errors"
  "github.com/google/uuid"
  "hash"
)

func IdFromString(id string) (Id, error) {
  uid, err := uuid.Parse(id)
  if err != nil {
    return NullId(), ErrMalformedUUID
  }
  return Id(uid), nil
}

func NewId() (Id, error) {
  uid, err := uuid.NewRandom()
  return Id(uid), err
}

func MustCreateId() Id {
  return DropError(NewId()) // TODO: Panic on error?
  //return types.Must(NewId())
}

func NullId() Id { return Id(uuid.UUID{}) }

type Id uuid.UUID

func (i Id) Underlying() uuid.UUID {
  return uuid.UUID(i)
}

func (i Id) String() string {
  if i == NullId() {
    return "" // Prevent bad uuid, which is "000000-..."
  }
  return i.Underlying().String()
}

func (i Id) Hash(hash hash.Hash) string {
  hash.Write([]byte(i.String()))
  return hex.EncodeToString(hash.Sum(nil))
}

func (i Id) EqWithString(uuid string) bool {
  return i.Underlying().String() == uuid
}

func (i Id) Eq(other Id) bool {
  return i.Underlying().String() == other.Underlying().String()
}

var ErrMalformedUUID = errors.New("value has malformed format for an UUID")
