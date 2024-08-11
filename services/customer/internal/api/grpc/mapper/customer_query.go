package mapper

import (
  customerv1 "github.com/arcorium/rashop/proto/gen/go/customer/v1"
  "github.com/arcorium/rashop/shared/errors"
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  "google.golang.org/protobuf/types/known/timestamppb"
  "mini-shop/services/user/internal/app/dto"
  "mini-shop/services/user/internal/app/query"
)

func ToGetCustomerByIdsQuery(request *customerv1.FindCustomerByIdsRequest) (query.GetCustomerByIdsQuery, error) {
  ids, ierr := sharedUtil.CastSliceErrs(request.CustomerIds, types.IdFromString)
  if !ierr.IsNil() {
    return query.GetCustomerByIdsQuery{}, ierr.ToGRPCError("customer_ids")
  }

  return query.GetCustomerByIdsQuery{
    CustomerIds: ids,
  }, nil
}

func ToGetCustomerAddressesQuery(request *customerv1.FindCustomerAddressesRequest) (query.GetCustomerAddressesQuery, error) {
  id, err := types.IdFromString(request.CustomerId)
  if err != nil {
    return query.GetCustomerAddressesQuery{}, errors.NewFieldError("customer_id", err).ToGrpcError()
  }

  return query.GetCustomerAddressesQuery{
    CustomerId: id,
  }, nil
}

func ToGetCustomerVouchersQuery(request *customerv1.FindCustomerVouchersRequest) (query.GetCustomerVouchersQuery, error) {
  id, err := types.IdFromString(request.CustomerId)
  if err != nil {
    return query.GetCustomerVouchersQuery{}, errors.NewFieldError("customer_id", err).ToGrpcError()
  }

  return query.GetCustomerVouchersQuery{
    CustomerId: id,
  }, nil
}

func ToProtoCustomer(dto *dto.CustomerResponseDTO) *customerv1.Customer {
  var lastModified *timestamppb.Timestamp
  if !dto.LastModifiedAt.IsZero() {
    lastModified = timestamppb.New(dto.LastModifiedAt)
  }

  return &customerv1.Customer{
    Id:                dto.Id.String(),
    Username:          dto.Username,
    FirstName:         dto.FirstName,
    LastName:          dto.LastName,
    Email:             dto.Email.String(),
    Balance:           dto.Balance,
    Point:             dto.Point,
    IsVerified:        dto.IsVerified,
    IsDisabled:        dto.IsDisabled,
    ShippingAddresses: sharedUtil.CastSliceP(dto.Addresses, ToProtoAddress),
    Vouchers:          sharedUtil.CastSliceP(dto.Vouchers, ToProtoVoucher),
    LastModifiedTime:  lastModified,
    CreatedTime:       timestamppb.New(dto.CreatedAt),
  }
}

func ToProtoAddress(dto *dto.AddressResponseDTO) *customerv1.Address {
  var lastModified *timestamppb.Timestamp
  if !dto.LastModifiedAt.IsZero() {
    lastModified = timestamppb.New(dto.LastModifiedAt)
  }
  return &customerv1.Address{
    Id:               dto.Id.String(),
    StreetAddress_1:  dto.StreetAddress1,
    StreetAddress_2:  dto.StreetAddress2,
    City:             dto.City,
    State:            dto.State,
    PostalCode:       dto.PostalCode,
    LastModifiedTime: lastModified,
    CreatedTime:      timestamppb.New(dto.CreatedAt),
  }
}

func ToProtoVoucher(dto *dto.VoucherResponseDTO) *customerv1.Voucher {
  return &customerv1.Voucher{
    Id:        dto.Id.String(),
    AddedTime: timestamppb.New(dto.AddedAt),
  }
}
