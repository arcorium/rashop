package config

import sharedConf "github.com/arcorium/rashop/shared/config"

type QueryServer struct {
  sharedConf.Server
  Database
  Broker MessageBroker
}
