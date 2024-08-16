package query

import (
  sharedDto "github.com/arcorium/rashop/shared/dto"
  "github.com/arcorium/rashop/shared/types"
)

type GetMailsQuery struct {
  sharedDto.PagedElementDTO
}

type FindMailByIdsQuery struct {
  Ids []types.Id
}
