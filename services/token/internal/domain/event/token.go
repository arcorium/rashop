package event

import (
  vob "github.com/arcorium/rashop/services/token/internal/domain/valueobject"
  "github.com/arcorium/rashop/shared/types"
  "time"
)

const (
  TokenCreatedEvent  = "token.created"
  TokenVerifiedEvent = "token.verified"
)

var _ types.Event = (*TokenCreatedV1)(nil)

type TokenCreatedV1 struct {
  DomainV1
  TokenId   string
  Token     string
  UserId    string
  Usage     vob.TokenUsage
  Type      vob.TokenType
  ExpiredAt time.Time
  CreatedAt time.Time
}

func (c *TokenCreatedV1) EventName() string {
  return TokenCreatedEvent
}

func (c *TokenCreatedV1) Key() (string, bool) {
  return c.TokenId, true
}

var _ types.Event = (*TokenVerifiedV1)(nil)

type TokenVerifiedV1 struct {
  DomainV1
  TokenId string
  Token   string
  UserId  string
}

func (c *TokenVerifiedV1) EventName() string {
  return TokenVerifiedEvent
}

func (c *TokenVerifiedV1) Key() (string, bool) {
  return c.TokenId, true
}
