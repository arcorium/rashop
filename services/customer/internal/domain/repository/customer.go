package repository

import (
  "context"
  "github.com/arcorium/rashop/shared/types"
  "github.com/arcorium/rashop/shared/util/repo"
  "rashop/services/customer/internal/domain/entity"
)

type ICustomer interface {
  Create(ctx context.Context, customer *entity.Customer) error
  Get(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Customer], error)
  FindByIds(ctx context.Context, ids ...types.Id) ([]entity.Customer, error)
  FindByEmails(ctx context.Context, emails ...types.Email) ([]entity.Customer, error)
  Update(ctx context.Context, customer *entity.Customer) error
}
