package model

import "github.com/uptrace/bun"

type MailRecipient struct {
  bun.BaseModel `bun:"table:mail_recipients,alias:mr"`

  MailId    string `bun:",type:uuid,pk"`
  Recipient string `bun:",pk"`

  Mail *Mail `bun:"rel:belongs-to,join:mail_id=id,on_delete:CASCADE"`
}
