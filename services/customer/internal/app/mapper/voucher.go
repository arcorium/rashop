package mapper

import (
  "rashop/services/customer/internal/app/dto"
  "rashop/services/customer/internal/domain/entity"
)

func ToVoucherResponseDTO(voucher *entity.Voucher) dto.VoucherResponseDTO {
  return dto.VoucherResponseDTO{
    Id:      voucher.Id,
    AddedAt: voucher.CreatedAt,
  }
}
