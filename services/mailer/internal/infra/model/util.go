package model

import (
  "github.com/arcorium/rashop/shared/database"
  "github.com/arcorium/rashop/shared/types"
  "github.com/uptrace/bun"
)

var models = []any{
  types.Nil[Mail](),
  types.Nil[MailRecipient](),
}

func RegisterBunModels(db *bun.DB) {
  database.RegisterBunModels(db, models...)
}

func CreateTables(db bun.IDB) error {
  return database.CreateTables(db, models...)
}
