package repository

import (
  "context"
  "github.com/arcorium/rashop/services/token/internal/domain/entity"
  "github.com/arcorium/rashop/shared/types"
)

type ITokenPersistent interface {
  FindByToken(ctx context.Context, token string) (entity.Token, error)
  FindByUserId(ctx context.Context, userId types.Id) ([]entity.Token, error)
  Create(ctx context.Context, token *entity.Token) error
  Delete(ctx context.Context, token *entity.Token) error
}
