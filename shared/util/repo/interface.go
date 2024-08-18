package repo

import (
  "context"
  "github.com/arcorium/rashop/shared/types"
)

type IGetPersistent[D IDataAccessModelWithError[T], T any] interface {
  Get(ctx context.Context, parameter QueryParameter) (PaginatedResult[D], error)
}

type IFindByIdsPersistent[D IDataAccessModelWithError[T], T any] interface {
  FindByIds(ctx context.Context, ids ...types.Id) ([]D, error)
}

type ICreatePersistent[D any] interface {
  Create(ctx context.Context, object D) error
}

type IUpdatePersistent[D any] interface {
  Update(ctx context.Context, object D) error
}

type IDeletePersistent interface {
  Delete(ctx context.Context, id types.Id) error
}

type IDeletesPersistent interface {
  Deletes(ctx context.Context, ids ...types.Id) error
}

type UOWBlock[T any] func(context.Context, T) error

type IUnitOfWork[V any] interface {
  // AsUnit create transaction and run function f. Repository storage on f should be in transaction
  // error returned is error forwarded from Block
  AsUnit(ctx context.Context, f UOWBlock[V]) error
}
