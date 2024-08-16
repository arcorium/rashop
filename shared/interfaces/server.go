package interfaces

import (
  "github.com/google/uuid"
  "time"
)

// IServer will provide each instance (server) with the different id and track when the server is up and running
type IServer interface {
  Identity() string
  StartAt() time.Time
  RunAt() time.Time
  StartupTime() time.Duration
}

func NewServer() ServerBase {
  return ServerBase{id: uuid.NewString(), startedAt: time.Now()}
}

type ServerBase struct {
  id        string
  startedAt time.Time
  runAt     time.Time
}

func (b *ServerBase) MarkRunning() {
  b.runAt = time.Now()
}

func (b *ServerBase) StartAt() time.Time {
  return b.startedAt
}

func (b *ServerBase) RunAt() time.Time {
  return b.runAt
}

func (b *ServerBase) StartupTime() time.Duration {
  return b.runAt.Sub(b.startedAt)
}

func (b *ServerBase) Identity() string {
  return b.id
}
