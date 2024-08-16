package cqrs

import (
  "rashop/services/customer/internal/domain/repository"
)

type CommonHandlerParameter struct {
  Repo      repository.ICustomer
  Publisher repository.IMessagePublisher
}
