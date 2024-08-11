package errors

import (
  "errors"
  "fmt"
  "github.com/go-playground/validator/v10"
  "google.golang.org/genproto/googleapis/rpc/errdetails"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
  "log"
)

func GrpcFieldErrors(details ...FieldError) error {
  var casted []*errdetails.BadRequest_FieldViolation
  for _, fe := range details {
    casted = append(casted, &errdetails.BadRequest_FieldViolation{
      Field:       fe.Name,
      Description: fe.Description.Error(),
    })
  }
  return grpcFieldErrors(casted...)
}

func grpcFieldErrors(details ...*errdetails.BadRequest_FieldViolation) error {
  badreq := &errdetails.BadRequest{}
  badreq.FieldViolations = append(badreq.FieldViolations, details...)

  st := status.New(codes.InvalidArgument, "Invalid Argument")
  stats, err := st.WithDetails(badreq)
  if err != nil {
    log.Fatalln("Fatal error: ", err)
  }

  return stats.Err()
}

func GrpcFieldValidationErrors(verr validator.ValidationErrors) error {
  var result []*errdetails.BadRequest_FieldViolation
  for _, val := range verr {
    result = append(result, &errdetails.BadRequest_FieldViolation{
      Field:       val.StructField(),
      Description: val.Error(),
    })
  }

  return grpcFieldErrors(result...)
}

func GrpcFieldIndexedErrors(field string, ierr IndicesError) error {
  if errors.Is(ierr.Err()[0], ErrEmptySlice) {
    return grpcFieldErrors(&errdetails.BadRequest_FieldViolation{
      Field:       field,
      Description: fmt.Sprintf("%s should not be empty", field),
    })
  }

  var result []*errdetails.BadRequest_FieldViolation
  for key, val := range ierr.Err() {
    result = append(result, &errdetails.BadRequest_FieldViolation{
      Field:       fmt.Sprintf("%s[%d]", field, key),
      Description: val.Error(),
    })
  }

  return grpcFieldErrors(result...)
}
