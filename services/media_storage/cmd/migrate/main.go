package main

import (
  "github.com/arcorium/rashop/services/media_storage/internal/infra/model"
  sharedConf "github.com/arcorium/rashop/shared/config"
  "github.com/arcorium/rashop/shared/database"
  "github.com/arcorium/rashop/shared/env"
  "log"
)

func main() {
  envName := ".env"
  if sharedConf.IsDebug() {
    envName = "dev.env"
  }
  _ = env.LoadEnvs(envName)

  dbConfig, err := sharedConf.Load[sharedConf.PostgresDatabase]()
  if err != nil {
    log.Fatalln(err)
  }

  db, err := database.OpenPostgresWithConfig(dbConfig, true)
  if err != nil {
    log.Fatalln(err)
  }
  defer db.Close()

  model.RegisterBunModels(db)

  if err = model.CreateTables(db); err != nil {
    log.Fatalln(err)
  }

  log.Println("Succeed migrate database: ", dbConfig.DSN())
}
