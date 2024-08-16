package config

import sharedConf "github.com/arcorium/rashop/shared/config"

type CommandServer struct {
  sharedConf.Server
  Database
  Broker MessageBroker
  SMTP   SMTP
}

type SMTP struct {
  Address  string `env:"SMTP_ADDRESS,notEmpty"`
  Port     uint16 `env:"SMTP_PORT,notEmpty"`
  Username string `env:"SMTP_USERNAME,notEmpty"`
  Password string `env:"SMTP_PASSWORD,notEmpty"`
}
