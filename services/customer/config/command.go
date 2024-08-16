package config

import sharedConf "github.com/arcorium/rashop/shared/config"

type CommandServer struct {
  sharedConf.Server
  Database
  Broker MessageBroker
}
