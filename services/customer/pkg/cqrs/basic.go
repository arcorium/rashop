package cqrs

import (
  "mini-shop/services/user/internal/domain/repository"
)

type CommonHandlerParameter struct {
  Repo      repository.ICustomer
  Publisher repository.IMessagePublisher
}
