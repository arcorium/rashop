package errors

import (
  "errors"
  "fmt"
  "strings"
)

func NewIndices(indices ...IndexedError) IndicesError {
  return IndicesError{errs: indices}
}

type IndicesError struct {
  errs []IndexedError
}

func (i IndicesError) Append(indexedError IndexedError) {
  i.errs = append(i.errs, indexedError)
}

func (i IndicesError) Error() string {
  var str strings.Builder
  for _, err := range i.errs {
    str.WriteString(err.Error() + "\n")
  }
  return str.String()
}

func (i IndicesError) IsNil() bool {
  return len(i.errs) == 0
}

func (i IndicesError) IsError() bool {
  return len(i.errs) > 0
}

func (i IndicesError) Err() []IndexedError {
  return i.errs
}

func (i IndicesError) ToGRPCError(fieldName string) error {
  return GrpcFieldIndexedErrors(fieldName, i)
}

func (i IndicesError) ToFieldError(fieldName string) FieldError {
  return NewFieldError(fieldName, i)
}

// IsEmptySlice check if the error is due to empty slice. Always check if it has error, otherwise
// it will panic due to indexing nil slice
func (i IndicesError) IsEmptySlice() bool {
  return i.errs[0].index == -1 || errors.Is(i.errs[0].err, ErrEmptySlice)
}

func NewIndex(i int, err error) IndexedError {
  return IndexedError{
    index: i,
    err:   err,
  }
}

type IndexedError struct {
  index int
  err   error
}

func (i IndexedError) Error() string {
  return fmt.Sprintf("%d: %v", i.index, i.err)
}

func (i IndexedError) IsNil() bool {
  return i.err == nil
}

func (i IndexedError) Err() (int, error) {
  return i.index, i.err
}
