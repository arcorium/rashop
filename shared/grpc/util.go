package grpc

import (
  "context"
  "errors"
  "github.com/arcorium/rashop/shared/types"
  "google.golang.org/grpc/metadata"
)

var (
  ErrMetadataNotFound = errors.New("metadata not found")
)

func ExtractUserId(ctx context.Context) (types.Id, error) {
  md, ok := metadata.FromIncomingContext(ctx)
  if !ok {
    return types.NullId(), ErrMetadataNotFound
  }
  users := md.Get(UserIdMetadataKey)
  if len(users) == 0 {
    return types.NullId(), ErrMetadataNotFound
  }

  id, err := types.IdFromString(users[0])
  return id, err
}
