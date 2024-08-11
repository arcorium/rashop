package status

import (
  "database/sql"
  "errors"
  "github.com/arcorium/rashop/shared/optional"
  "github.com/arcorium/rashop/shared/util/repo"
)

// FromRepository override only the code when the error is sql.ErrNoRows and use the sql.ErrNoRows as the error
func FromRepository(err error, notFoundStatus ...Object) Object {
  if errors.Is(err, sql.ErrNoRows) {
    return ErrNotFound()
  }
  if len(notFoundStatus) == 1 {
    return notFoundStatus[0]
  }
  return Error(REPOSITORY_ERROR, err)
}

func FromRepository2(err error, notFound optional.Object[Object], exists optional.Object[Object]) Object {
  if errors.Is(err, sql.ErrNoRows) {
    return notFound.ValueOr(ErrNotFound())
  } else if errors.Is(err, repo.ErrAlreadyExists) {
    return exists.ValueOr(ErrAlreadyExist())
  }
  return Error(REPOSITORY_ERROR, err)
}

// FromRepositoryOverride override the code and error when the error is sql.ErrNoRows
func FromRepositoryOverride(err error, notFoundOver ...Object) Object {
  if errors.Is(err, sql.ErrNoRows) && len(notFoundOver) > 0 {
    return notFoundOver[0]
  }
  return Error(REPOSITORY_ERROR, err)
}

// FromRepositoryExist helper function to call FromRepository2 which handle sql.ErrNoRows and integrity violation as
// object already exists error. Used for inserting new data
func FromRepositoryExist(err error) Object {
  return FromRepository2(err, optional.Some(ErrAlreadyExist()), optional.Some(ErrAlreadyExist()))
}
