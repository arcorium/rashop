package intev

import (
  "github.com/arcorium/rashop/shared/types"
  "time"
)

const (
  EmailVerificationTokenCreatedEvent = "token.email_verification.created"
  ResetPasswordTokenCreatedEvent     = "token.reset_password.created"
  LoginTokenCreatedEvent             = "token.login.created"
)

// TODO: Add userId for each event when user is separated from customer

var _ types.Event = (*EmailVerificationTokenCreatedV1)(nil)

func NewEmailVerificationTokenCreated(tokenId types.Id, token string, recipient types.Email, expiryTime time.Time) *EmailVerificationTokenCreatedV1 {
  return &EmailVerificationTokenCreatedV1{
    IntegrationV1: NewV1(),
    TokenId:       tokenId.String(),
    Token:         token,
    Recipient:     recipient.String(),
    ExpiryTime:    expiryTime,
  }
}

type EmailVerificationTokenCreatedV1 struct {
  IntegrationV1
  TokenId    string
  Token      string
  Recipient  string
  ExpiryTime time.Time
}

func (c *EmailVerificationTokenCreatedV1) EventName() string {
  return EmailVerificationTokenCreatedEvent
}

func (c *EmailVerificationTokenCreatedV1) Key() (string, bool) {
  return c.TokenId, true
}

var _ types.Event = (*ResetPasswordTokenCreatedV1)(nil)

func NewResetPasswordTokenCreated(tokenId types.Id, token string, recipient types.Email, expiryTime time.Time) *ResetPasswordTokenCreatedV1 {
  return &ResetPasswordTokenCreatedV1{
    IntegrationV1: NewV1(),
    TokenId:       tokenId.String(),
    Token:         token,
    Recipient:     recipient.String(),
    ExpiryTime:    expiryTime,
  }
}

type ResetPasswordTokenCreatedV1 struct {
  IntegrationV1
  TokenId    string
  Token      string
  Recipient  string
  ExpiryTime time.Time
}

func (c *ResetPasswordTokenCreatedV1) EventName() string {
  return ResetPasswordTokenCreatedEvent
}

func (c *ResetPasswordTokenCreatedV1) Key() (string, bool) {
  return c.TokenId, true
}

var _ types.Event = (*LoginTokenCreatedV1)(nil)

func NewLoginTokenCreated(tokenId types.Id, token string, recipient types.Email, expiryTime time.Time) *LoginTokenCreatedV1 {
  return &LoginTokenCreatedV1{
    IntegrationV1: NewV1(),
    TokenId:       tokenId.String(),
    Token:         token,
    Recipient:     recipient.String(),
    ExpiryTime:    expiryTime,
  }
}

type LoginTokenCreatedV1 struct {
  IntegrationV1
  TokenId    string
  Token      string
  Recipient  string
  ExpiryTime time.Time
}

func (c *LoginTokenCreatedV1) EventName() string {
  return LoginTokenCreatedEvent
}

func (c *LoginTokenCreatedV1) Key() (string, bool) {
  return c.TokenId, true
}
