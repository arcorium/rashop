package errors

import "fmt"

func NewFieldError(field string, desc error) FieldError {
  return FieldError{
    Name:        field,
    Description: desc,
  }
}

type FieldError struct {
  Name        string
  Description error
}

func (e FieldError) Error() string {
  return fmt.Sprintf("%s: %s", e.Name, e.Description)
}

func (e FieldError) ToGrpcError() error {
  return GrpcFieldErrors(e)
}
