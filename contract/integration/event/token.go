package intev

import (
  "github.com/arcorium/rashop/contract/integration/enum"
  "github.com/arcorium/rashop/shared/types"
  "time"
)

const (
  TokenCreatedEvent  = "token.created"
  TokenVerifiedEvent = "token.verified"
)

var _ types.Event = (*TokenCreatedV1)(nil)

func NewTokenCreated(userId, tokenId types.Id, token string, usage, types uint8, expiryTime, createdAt time.Time) *TokenCreatedV1 {
  return &TokenCreatedV1{
    IntegrationV1: NewV1(),
    TokenId:       tokenId.String(),
    Token:         token,
    UserId:        userId.String(),
    Usage:         enum.TokenUsage(usage),
    Type:          enum.TokenType(types),
    ExpiryTime:    expiryTime,
    CreatedAt:     createdAt,
  }
}

type TokenCreatedV1 struct {
  IntegrationV1
  TokenId    string
  Token      string
  UserId     string
  Usage      enum.TokenUsage
  Type       enum.TokenType
  ExpiryTime time.Time
  CreatedAt  time.Time
}

func (c *TokenCreatedV1) EventName() string {
  return TokenCreatedEvent
}

func (c *TokenCreatedV1) Key() (string, bool) {
  return c.TokenId, true
}

var _ types.Event = (*TokenVerifiedV1)(nil)

func NewTokenVerified(userId, tokenId types.Id, token string, usage, types uint8, expiryTime, createdAt time.Time) *TokenVerifiedV1 {
  return &TokenVerifiedV1{
    IntegrationV1: NewV1(),
    TokenId:       tokenId.String(),
    Token:         token,
    UserId:        userId.String(),
    Usage:         enum.TokenUsage(usage),
    Type:          enum.TokenType(types),
    ExpiryTime:    expiryTime,
    CreatedAt:     createdAt,
  }
}

type TokenVerifiedV1 struct {
  IntegrationV1
  TokenId    string
  Token      string
  UserId     string
  Usage      enum.TokenUsage
  Type       enum.TokenType
  ExpiryTime time.Time
  CreatedAt  time.Time
}

func (c *TokenVerifiedV1) EventName() string {
  return TokenVerifiedEvent
}

func (c *TokenVerifiedV1) Key() (string, bool) {
  return c.TokenId, true
}
