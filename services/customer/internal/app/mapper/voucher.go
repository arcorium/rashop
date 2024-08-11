package mapper

import (
  "mini-shop/services/user/internal/app/dto"
  "mini-shop/services/user/internal/domain/entity"
)

func ToVoucherResponseDTO(voucher *entity.Voucher) dto.VoucherResponseDTO {
  return dto.VoucherResponseDTO{
    Id:      voucher.Id,
    AddedAt: voucher.CreatedAt,
  }
}
