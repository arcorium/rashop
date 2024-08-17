package repo

import (
  "context"
  "github.com/arcorium/rashop/shared/types"
)

type IGetRepository[D IDataAccessModelWithError[T], T any] interface {
  Get(ctx context.Context, parameter QueryParameter) (PaginatedResult[D], error)
}

type IFindByIdsRepository[D IDataAccessModelWithError[T], T any] interface {
  FindByIds(ctx context.Context, ids ...types.Id) ([]D, error)
}

type ICreateRepository[D IDataAccessModelWithError[T], T any] interface {
  Create(ctx context.Context, object D) error
}

type IUpdateRepository[D IDataAccessModelWithError[T], T any] interface {
  Update(ctx context.Context, object D) error
}

type IDeleteRepository[D IDataAccessModelWithError[T], T any] interface {
  Delete(ctx context.Context, id types.Id) error
}

type IDeletesRepository[D IDataAccessModelWithError[T], T any] interface {
  Deletes(ctx context.Context, ids ...types.Id) error
}

type UOWBlock[T any] func(context.Context, T) error

type IUnitOfWork[V any] interface {
  // AsUnit create transaction and run function f. Repository storage on f should be in transaction
  // error returned is error forwarded from Block
  AsUnit(ctx context.Context, f UOWBlock[V]) error
}
