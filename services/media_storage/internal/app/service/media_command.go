package service

import (
  "context"
  "github.com/arcorium/rashop/services/media_storage/internal/app/command"
  "github.com/arcorium/rashop/shared/status"
  "github.com/arcorium/rashop/shared/types"
)

type IMediaCommand interface {
  Store(ctx context.Context, cmd *command.StoreMediaCommand) (types.Id, status.Object)
  Delete(ctx context.Context, cmd *command.DeleteMediaCommand) status.Object
}

func NewMediaCommand(config MediaCommandConfig) IMediaCommand {
  return &mediaCommandService{
    i: config,
  }
}

func DefaultMediaCommandConfig(parameter command.CommonHandlerParameter) MediaCommandConfig {
  return MediaCommandConfig{
    Store:  command.NewStoreMediaHandler(parameter),
    Delete: command.NewDeleteMediaHandler(parameter),
  }
}

type MediaCommandConfig struct {
  Store  command.IStoreMediaHandler
  Delete command.IDeleteMediaHandler
}

type mediaCommandService struct {
  i MediaCommandConfig
}

func (m *mediaCommandService) Store(ctx context.Context, cmd *command.StoreMediaCommand) (types.Id, status.Object) {
  return m.i.Store.Handle(ctx, cmd)
}

func (m *mediaCommandService) Delete(ctx context.Context, cmd *command.DeleteMediaCommand) status.Object {
  return m.i.Delete.Handle(ctx, cmd)
}
