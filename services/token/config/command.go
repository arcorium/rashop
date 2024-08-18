package config

import (
  sharedConf "github.com/arcorium/rashop/shared/config"
  "time"
)

type CommandServer struct {
  sharedConf.Server
  Duration Duration
  Database
  Broker MessageBroker
}

type Duration struct {
  VerificationTokenExpiryTime time.Duration `env:"VERIFICATION_TOKEN_EXPIRY_TIME" envDefault:"1h"`
  ResetTokenExpiryTime        time.Duration `env:"RESET_TOKEN_EXPIRY_TIME" envDefault:"1h"`
  LoginTokenExpiryTime        time.Duration `env:"LOGIN_TOKEN_EXPIRY_TIME" envDefault:"1h"`
  GeneralTokenExpiryTime      time.Duration `env:"GENERAL_TOKEN_EXPIRY_TIME" envDefault:"1h"`

  SingleExpirationTime time.Duration `env:"SINGLE_EXPIRATION_TIME"`
}

func (d *Duration) IsSingle() bool {
  return d.SingleExpirationTime != 0
}
