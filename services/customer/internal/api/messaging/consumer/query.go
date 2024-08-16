package consumer

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "rashop/services/customer/internal/app/service"
  "rashop/services/customer/internal/domain/event"
  "rashop/services/customer/pkg/tracer"
)

func NewCustomerQueryHandler(svc service.ICustomerQueryConsumer) ICustomerQueryHandler {
  return &customerQueryConsumerHandler{
    svc:    svc,
    tracer: tracer.Get(),
  }
}

type ICustomerQueryHandler interface {
  OnAddressAddedV1(ctx context.Context, ev *event.CustomerAddressAddedV1) error
  OnAddressDeletedV1(ctx context.Context, ev *event.CustomerAddressDeletedV1) error
  OnAddressUpdatedV1(ctx context.Context, ev *event.CustomerAddressUpdatedV1) error
  OnBalanceUpdatedV1(ctx context.Context, ev *event.CustomerBalanceUpdatedV1) error
  OnCreatedV1(ctx context.Context, ev *event.CustomerCreatedV1) error
  OnDefaultAddressUpdatedV1(ctx context.Context, ev *event.CustomerDefaultAddressUpdatedV1) error
  OnEmailVerifiedV1(ctx context.Context, ev *event.CustomerEmailVerifiedV1) error
  OnPasswordUpdatedV1(ctx context.Context, ev *event.CustomerPasswordUpdatedV1) error
  OnPhotoUpdatedV1(ctx context.Context, ev *event.CustomerPhotoUpdatedV1) error
  OnStatusUpdatedV1(ctx context.Context, ev *event.CustomerStatusUpdatedV1) error
  OnUpdatedV1(ctx context.Context, ev *event.CustomerUpdatedV1) error
  OnVoucherAddedV1(ctx context.Context, ev *event.CustomerVoucherAddedV1) error
  OnVoucherDeletedV1(ctx context.Context, ev *event.CustomerVoucherDeletedV1) error
  OnVoucherUpdatedV1(ctx context.Context, ev *event.CustomerVoucherUpdatedV1) error
}

type customerQueryConsumerHandler struct {
  svc    service.ICustomerQueryConsumer
  tracer trace.Tracer
}

func (c *customerQueryConsumerHandler) OnAddressAddedV1(ctx context.Context, ev *event.CustomerAddressAddedV1) error {
  ctx, span := c.tracer.Start(ctx, "customerQueryConsumerHandler.OnAddressAddedV1")
  defer span.End()

  stat := c.svc.AddressAddedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *customerQueryConsumerHandler) OnAddressDeletedV1(ctx context.Context, ev *event.CustomerAddressDeletedV1) error {
  ctx, span := c.tracer.Start(ctx, "customerQueryConsumerHandler.OnAddressDeletedV1")
  defer span.End()

  stat := c.svc.AddressDeletedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *customerQueryConsumerHandler) OnAddressUpdatedV1(ctx context.Context, ev *event.CustomerAddressUpdatedV1) error {
  ctx, span := c.tracer.Start(ctx, "customerQueryConsumerHandler.OnAddressUpdatedV1")
  defer span.End()

  stat := c.svc.AddressUpdatedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *customerQueryConsumerHandler) OnBalanceUpdatedV1(ctx context.Context, ev *event.CustomerBalanceUpdatedV1) error {
  ctx, span := c.tracer.Start(ctx, "customerQueryConsumerHandler.OnBalanceUpdatedV1")
  defer span.End()

  stat := c.svc.BalanceUpdatedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *customerQueryConsumerHandler) OnCreatedV1(ctx context.Context, ev *event.CustomerCreatedV1) error {
  ctx, span := c.tracer.Start(ctx, "customerQueryConsumerHandler.OnCreatedV1")
  defer span.End()

  stat := c.svc.CreatedV1(ctx, ev)
  return stat.ToGRPCErrorWithSpan(span)
}

func (c *customerQueryConsumerHandler) OnDefaultAddressUpdatedV1(ctx context.Context, v1 *event.CustomerDefaultAddressUpdatedV1) error {
  ctx, span := c.tracer.Start(ctx, "customerQueryConsumerHandler.OnDefaultAddressUpdatedV1")
  defer span.End()

  stat := c.svc.DefaultAddressSetV1(ctx, v1)
  return stat.ErrorWithSpan(span)
}

func (c *customerQueryConsumerHandler) OnEmailVerifiedV1(ctx context.Context, ev *event.CustomerEmailVerifiedV1) error {
  ctx, span := c.tracer.Start(ctx, "customerQueryConsumerHandler.OnEmailVerifiedV1")
  defer span.End()

  stat := c.svc.EmailVerifiedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *customerQueryConsumerHandler) OnPasswordUpdatedV1(ctx context.Context, ev *event.CustomerPasswordUpdatedV1) error {
  ctx, span := c.tracer.Start(ctx, "customerQueryConsumerHandler.OnPasswordUpdatedV1")
  defer span.End()

  stat := c.svc.PasswordUpdatedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *customerQueryConsumerHandler) OnPhotoUpdatedV1(ctx context.Context, ev *event.CustomerPhotoUpdatedV1) error {
  ctx, span := c.tracer.Start(ctx, "customerQueryConsumerHandler.OnPhotoUpdatedV1")
  defer span.End()

  stat := c.svc.PhotoUpdatedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *customerQueryConsumerHandler) OnStatusUpdatedV1(ctx context.Context, ev *event.CustomerStatusUpdatedV1) error {
  ctx, span := c.tracer.Start(ctx, "customerQueryConsumerHandler.OnStatusUpdatedV1")
  defer span.End()

  stat := c.svc.StatusUpdatedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *customerQueryConsumerHandler) OnUpdatedV1(ctx context.Context, ev *event.CustomerUpdatedV1) error {
  ctx, span := c.tracer.Start(ctx, "customerQueryConsumerHandler.OnUpdatedV1")
  defer span.End()

  stat := c.svc.UpdatedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *customerQueryConsumerHandler) OnVoucherAddedV1(ctx context.Context, ev *event.CustomerVoucherAddedV1) error {
  ctx, span := c.tracer.Start(ctx, "customerQueryConsumerHandler.OnVoucherAddedV1")
  defer span.End()

  stat := c.svc.VoucherAddedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *customerQueryConsumerHandler) OnVoucherDeletedV1(ctx context.Context, ev *event.CustomerVoucherDeletedV1) error {
  ctx, span := c.tracer.Start(ctx, "customerQueryConsumerHandler.OnVoucherDeletedV1")
  defer span.End()

  stat := c.svc.VoucherDeletedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *customerQueryConsumerHandler) OnVoucherUpdatedV1(ctx context.Context, ev *event.CustomerVoucherUpdatedV1) error {
  ctx, span := c.tracer.Start(ctx, "customerQueryConsumerHandler.OnVoucherUpdatedV1")
  defer span.End()

  stat := c.svc.VoucherUpdatedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}
