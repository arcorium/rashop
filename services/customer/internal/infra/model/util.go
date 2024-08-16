package model

import (
  "context"
  "github.com/arcorium/rashop/shared/database"
  "github.com/arcorium/rashop/shared/types"
  "github.com/uptrace/bun"
)

var modelsPair = []types.Pair[any, bool]{
  {types.Nil[User](), false},
  {types.Nil[Customer](), true},
  {types.Nil[ShippingAddress](), true},
  {types.Nil[Voucher](), true},
}

var models = []any{
  types.Nil[Customer](),
  types.Nil[User](),
  types.Nil[ShippingAddress](),
  types.Nil[Voucher](),
}

func RegisterBunModels(db *bun.DB) {
  database.RegisterBunModels(db, models...)
}

func CreateTables(db bun.IDB) error {
  // Need custom, because the user should not have foreign keys to profiles
  ctx := context.Background()
  for _, pair := range modelsPair {
    q := db.NewCreateTable().
      Model(pair.First).
      IfNotExists()

    if pair.Second {
      q = q.WithForeignKeys()
    }

    _, err := q.Exec(ctx)

    if err != nil {
      return err
    }
  }
  return nil
}
