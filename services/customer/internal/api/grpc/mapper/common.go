package mapper

import (
  "github.com/arcorium/rashop/proto/gen/go/common"
  "github.com/arcorium/rashop/shared/dto"
)

type IPagedRequest interface {
  GetDetail() *common.PagedElementRequest
}

func ToPagedParameter(request IPagedRequest) dto.PagedElementDTO {
  if request == nil || request.GetDetail() == nil {
    return dto.PagedElementDTO{}
  }

  detail := request.GetDetail()
  return dto.PagedElementDTO{
    Element: detail.Element,
    Page:    detail.Page,
  }
}

func ToProtoPagedElementDetails[T any](result *dto.PagedElementResult[T]) *common.PagedElementResponseDetail {
  return &common.PagedElementResponseDetail{
    Element:       result.Element,
    Page:          result.Page,
    TotalElements: result.TotalElements,
    TotalPages:    result.TotalPages,
  }
}
