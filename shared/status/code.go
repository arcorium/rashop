package status

import (
  "github.com/arcorium/rashop/shared/optional"
)

type Code uint

const (
  // Succeed Code
  SUCCESS Code = iota
  CREATED
  UPDATED
  DELETED

  // Error Code
  // General
  INTERNAL_SERVER_ERROR
  EXTERNAL_SERVICE_ERROR
  REPOSITORY_ERROR
  BAD_REQUEST_ERROR
  FIELD_VALIDATION_ERROR
  NOT_AUTHORIZED_ERROR
  NOT_AUTHENTICATED_ERROR
  SERVICE_UNAVAILABLE_ERROR
  // Other
  OBJECT_NOT_FOUND
  OBJECT_ALREADY_EXIST
)

// Succeed Code Helper
func Succeed() Object {
  return New(SUCCESS, nil)
}

func SomeSuccess() optional.Object[Object] {
  return optional.Some(Succeed())
}

func Created() Object {
  return New(CREATED, nil)
}

func SomeCreated() optional.Object[Object] {
  return optional.Some(Created())
}

func Updated() Object {
  return New(UPDATED, nil)
}

func SomeUpdated() optional.Object[Object] {
  return optional.Some(Updated())
}

func Deleted() Object {
  return New(DELETED, nil)
}

func SomeDeleted() optional.Object[Object] {
  return optional.Some(Deleted())
}

// Error Code Helper
func ErrInternal(err error) Object {
  return New(INTERNAL_SERVER_ERROR, err)
}

func ErrExternal(err error) Object {
  return New(EXTERNAL_SERVICE_ERROR, err)
}

func ErrUnAuthorized(err error) Object {
  return New(NOT_AUTHORIZED_ERROR, err)
}

func ErrUnAuthenticated(err error) Object {
  return New(NOT_AUTHENTICATED_ERROR, err)
}

func ErrBadRequest(err error) Object {
  return New(BAD_REQUEST_ERROR, err)
}

func ErrFieldValidation(err error) Object {
  return New(FIELD_VALIDATION_ERROR, err)
}

func ErrNotFound() Object {
  return NewWithMessage(OBJECT_NOT_FOUND, "Object not found")
}

func ErrAlreadyExist() Object {
  return NewWithMessage(OBJECT_ALREADY_EXIST, "Object already exists")
}
