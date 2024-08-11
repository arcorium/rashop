package types

import (
  "errors"
  "golang.org/x/exp/maps"
)

var ErrDuplicatedValue = errors.New("set: value duplicated")
var ErrValueNotFound = errors.New("set: value not found")

type Set[T comparable] map[T]struct{}

func (s Set[T]) IsExist(val T) bool {
  _, ok := s[val]
  return ok
}

// Add will append value into set and when the value already exists it will be ignored
func (s Set[T]) Add(val T) {
  s[val] = struct{}{}
}

// TryAdd works like Add but when the value already exist it will return ErrDuplicatedValue
func (s Set[T]) TryAdd(val T) error {
  if s.IsExist(val) {
    return ErrDuplicatedValue
  }
  s.Add(val)
  return nil
}

// Delete will remove value from set and if the value doesn't exist it will be ignored
func (s Set[T]) Delete(val T) {
  delete(s, val)
}

// TryDelete works like Delete, but when the value doesn't exist it will return ErrValueNotFound
func (s Set[T]) TryDelete(val T) error {
  if !s.IsExist(val) {
    return ErrValueNotFound
  }
  s.Delete(val)
  return nil
}

// Values get all set values as slice
func (s Set[T]) Values() []T {
  return maps.Keys(s)
}

type StringSet = Set[string]
