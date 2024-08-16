package repository

import (
  "context"
  "github.com/arcorium/rashop/services/mailer/internal/domain/entity"
  "github.com/arcorium/rashop/shared/types"
  "github.com/arcorium/rashop/shared/util/repo"
  "time"
)

type IMailPersistent interface {
  Get(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Mail], error)
  FindByIds(ctx context.Context, mailIds ...types.Id) ([]entity.Mail, error)
  Create(ctx context.Context, mail *entity.Mail) error
  Update(ctx context.Context, mail *entity.Mail) error
  Delete(ctx context.Context, startTime, endTime time.Time) (uint64, error)
}
