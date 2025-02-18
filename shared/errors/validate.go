package errors

import (
  "fmt"
  "google.golang.org/genproto/googleapis/rpc/errdetails"
  "strings"
)

type EmptyFieldError struct {
  Errs []FieldError
}

func (f EmptyFieldError) IsNil() bool {
  return f.Errs == nil
}

func (f EmptyFieldError) Error() string {
  var str strings.Builder

  for _, err := range f.Errs {
    str.WriteString(err.Error() + "\n")
  }

  return str.String()
}

func (f EmptyFieldError) ToGRPCError() error {
  var result []*errdetails.BadRequest_FieldViolation
  for key, val := range f.Errs {
    result = append(result, &errdetails.BadRequest_FieldViolation{
      Field:       fmt.Sprintf("%s[%d]", val.Name, key),
      Description: val.Description.Error(),
    })
  }

  return grpcFieldErrors(result...)
}
