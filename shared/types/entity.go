package types

type Entity interface {
  Identity() string
}

func NewEntity(id string) EntityBase {
  return EntityBase{id: id}
}

type EntityBase struct {
  id string
}

func (e EntityBase) Identity() string {
  return e.id
}
