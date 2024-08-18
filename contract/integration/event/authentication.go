package intev

import "github.com/arcorium/rashop/shared/types"

const (
  LoginTokenRequestedEvent = "authN.login-token.requested"
)

var _ types.Event = (*LoginTokenRequestedV1)(nil)

func NewLoginTokenRequestedV1(userId types.Id, email types.Email, username string) *LoginTokenRequestedV1 {
  return &LoginTokenRequestedV1{
    IntegrationV1: NewV1(),
    UserId:        userId.String(),
    Email:         email.String(),
    Username:      username,
  }
}

type LoginTokenRequestedV1 struct {
  IntegrationV1
  UserId   string
  Email    string
  Username string
}

func (c *LoginTokenRequestedV1) EventName() string {
  return LoginTokenRequestedEvent
}
