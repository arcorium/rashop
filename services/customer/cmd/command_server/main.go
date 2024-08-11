package main

import (
  sharedConf "github.com/arcorium/rashop/shared/config"
  "github.com/arcorium/rashop/shared/env"
  "github.com/arcorium/rashop/shared/logger"
  "log"
  "mini-shop/services/user/config"
)

func main() {
  var filename = ".env"
  if sharedConf.IsDebug() {
    filename = "dev.env"
  }
  err := env.LoadEnvs(filename)

  // Config
  serverConfig, err := sharedConf.Load[config.CommandServer]()
  if err != nil {
    env.LogError(err, -1)
  }

  // Init global logger
  logg, err := logger.NewZapLogger(sharedConf.IsDebug())
  if err != nil {
    log.Fatalln(err)
  }
  logger.SetGlobal(logg)

  server, err := NewServer(serverConfig)
  if err != nil {
    log.Fatalln(err)
  }
  log.Fatalln(server.Run())
}
