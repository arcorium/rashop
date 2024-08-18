package service

import (
  "context"
  "github.com/arcorium/rashop/services/token/internal/app/command"
  "github.com/arcorium/rashop/services/token/internal/app/dto"
  "github.com/arcorium/rashop/shared/status"
)

type ITokenCommand interface {
  Generate(ctx context.Context, cmd *command.GenerateTokenCommand) (dto.TokenResponseDTO, status.Object)
  Verify(ctx context.Context, cmd *command.VerifyTokenCommand) status.Object
}

func NewTokenCommand(config TokenCommandFactory) ITokenCommand {
  return &tokenCommandService{
    i: config,
  }
}

func DefaultTokenCommandFactory(parameter command.CommonHandlerParameter, config command.TokenGenerationConfig) TokenCommandFactory {
  return TokenCommandFactory{
    Generate: command.NewGenerateTokenHandler(parameter, config),
    Verify:   command.NewVerifyTokenHandler(parameter),
  }
}

type TokenCommandFactory struct {
  Generate command.IGenerateTokenHandler
  Verify   command.IVerifyTokenHandler
}

type tokenCommandService struct {
  i TokenCommandFactory
}

func (t *tokenCommandService) Generate(ctx context.Context, cmd *command.GenerateTokenCommand) (dto.TokenResponseDTO, status.Object) {
  return t.i.Generate.Handle(ctx, cmd)
}

func (t *tokenCommandService) Verify(ctx context.Context, cmd *command.VerifyTokenCommand) status.Object {
  return t.i.Verify.Handle(ctx, cmd)
}
