package repo

type OrderedIdInput struct {
  Id string `bun:"input_id,type:uuid"`
}
