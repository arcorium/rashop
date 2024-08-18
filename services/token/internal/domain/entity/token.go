package entity

import (
  "errors"
  "github.com/arcorium/rashop/services/token/internal/domain/event"
  vob "github.com/arcorium/rashop/services/token/internal/domain/valueobject"
  "github.com/arcorium/rashop/shared/types"
  "time"
)

var (
  ErrTokenIsExpired         = errors.New("token is expired")
  ErrTokenHasDifferentUsage = errors.New("token is not expected to be used here")
)

func CreateToken(token *Token) (Token, error) {
  result := Token{}

  ev := &event.TokenCreatedV1{
    DomainV1:  event.NewV1(),
    TokenId:   token.Id.String(),
    Token:     token.Token,
    UserId:    token.UserId.String(),
    Usage:     token.Usage,
    Type:      token.Type,
    ExpiredAt: token.ExpiredAt,
    CreatedAt: token.GeneratedAt,
  }

  result.AddEvents(ev)
  return result, nil
}

var _ types.Aggregate = (*Token)(nil)

type Token struct {
  types.AggregateBase
  types.AggregateHelper

  Id          types.Id
  Token       string
  UserId      types.Id
  Usage       vob.TokenUsage
  Type        vob.TokenType
  ExpiredAt   time.Time
  GeneratedAt time.Time
}

func (t *Token) Verify(expectedUsage vob.TokenUsage) (types.Event, error) {
  if t.Usage != expectedUsage {
    return nil, ErrTokenHasDifferentUsage
  }

  if t.IsExpired() {
    return nil, ErrTokenIsExpired
  }

  ev := &event.TokenVerifiedV1{
    DomainV1: event.NewV1(),
    TokenId:  t.Id.String(),
    Token:    t.Token,
    UserId:   t.UserId.String(),
  }
  t.AddEvents(ev)
  return ev, nil
}

func (t *Token) IsExpired() bool {
  return t.ExpiredAt.After(time.Now())
}

func (t *Token) Identity() string {
  return t.Id.String()
}

func (t *Token) ApplyEvent(ev types.Event) error {
  switch cur := ev.(type) {
  case *event.TokenCreatedV1:
    t.Id = types.Must(types.IdFromString(cur.TokenId))
    t.Token = cur.Token
    t.Usage = cur.Usage
    t.ExpiredAt = cur.ExpiredAt
    t.Type = cur.Type
    t.ExpiredAt = cur.ExpiredAt
    t.GeneratedAt = cur.CreatedAt
    t.MarkCreated()
  case *event.TokenVerifiedV1:
    t.MarkDeleted()
  default:
    return errors.New("unknown event")
  }
  return nil
}
