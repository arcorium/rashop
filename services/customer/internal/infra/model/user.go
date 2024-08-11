package model

import (
  "github.com/uptrace/bun"
  "time"
)

type User struct {
  bun.BaseModel `bun:"table:users,alias:u"`

  Id         string    `bun:",type:uuid,nullzero,pk"`
  Username   string    `bun:",nullzero,notnull"`
  Email      string    `bun:",nullzero,notnull"`
  Password   string    `bun:",nullzero,notnull"`
  IsVerified bool      `bun:",default:false"`
  UpdatedAt  time.Time `bun:",nullzero"`
  CreatedAt  time.Time `bun:",notnull,nullzero"`

  Customer *Customer `bun:"rel:has-one,join:id=user_id"`
}
