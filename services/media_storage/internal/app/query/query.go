package query

import "github.com/arcorium/rashop/shared/types"

type GetMediaQuery struct {
  MediaId types.Id
}

type GetMediaMetadataQuery struct {
  MediaIds []types.Id
}
