package consumer

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "mini-shop/services/user/internal/app/service"
  "mini-shop/services/user/internal/domain/event"
  "mini-shop/services/user/pkg/tracer"
)

func NewCustomerQueryHandler(svc service.ICustomerQueryConsumer) ICustomerQueryHandler {
  return &CustomerQueryConsumer{
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
  OnDefaultAddressSetV1(ctx context.Context, v1 *event.CustomerDefaultAddressUpdatedV1) error
  OnDisabledV1(ctx context.Context, ev *event.CustomerStatusUpdatedV1) error
  OnEmailVerifiedV1(ctx context.Context, ev *event.CustomerEmailVerifiedV1) error
  OnEnabledV1(ctx context.Context, ev *event.CustomerStatusUpdatedV1) error
  OnPasswordUpdatedV1(ctx context.Context, ev *event.CustomerPasswordUpdatedV1) error
  OnPhotoUpdatedV1(ctx context.Context, ev *event.CustomerPhotoUpdatedV1) error
  OnUpdatedV1(ctx context.Context, ev *event.CustomerUpdatedV1) error
  OnVoucherAddedV1(ctx context.Context, ev *event.CustomerVoucherAddedV1) error
  OnVoucherDeletedV1(ctx context.Context, ev *event.CustomerVoucherDeletedV1) error
  OnVoucherUpdatedV1(ctx context.Context, ev *event.CustomerVoucherUpdatedV1) error
}

type CustomerQueryConsumer struct {
  svc    service.ICustomerQueryConsumer
  tracer trace.Tracer
}

func (c *CustomerQueryConsumer) OnAddressAddedV1(ctx context.Context, ev *event.CustomerAddressAddedV1) error {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryConsumer.OnAddressAddedV1")
  defer span.End()

  stat := c.svc.AddressAddedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *CustomerQueryConsumer) OnAddressDeletedV1(ctx context.Context, ev *event.CustomerAddressDeletedV1) error {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryConsumer.OnAddressDeletedV1")
  defer span.End()

  stat := c.svc.AddressDeletedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *CustomerQueryConsumer) OnAddressUpdatedV1(ctx context.Context, ev *event.CustomerAddressUpdatedV1) error {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryConsumer.OnAddressUpdatedV1")
  defer span.End()

  stat := c.svc.AddressUpdatedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *CustomerQueryConsumer) OnBalanceUpdatedV1(ctx context.Context, ev *event.CustomerBalanceUpdatedV1) error {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryConsumer.OnBalanceUpdatedV1")
  defer span.End()

  stat := c.svc.BalanceUpdatedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *CustomerQueryConsumer) OnCreatedV1(ctx context.Context, ev *event.CustomerCreatedV1) error {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryConsumer.OnCreatedV1")
  defer span.End()

  stat := c.svc.CreatedV1(ctx, ev)
  return stat.ToGRPCErrorWithSpan(span)
}

func (c *CustomerQueryConsumer) OnDefaultAddressSetV1(ctx context.Context, v1 *event.CustomerDefaultAddressUpdatedV1) error {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryConsumer.OnDefaultAddressSetV1")
  defer span.End()

  stat := c.svc.DefaultAddressSetV1(ctx, v1)
  return stat.ErrorWithSpan(span)
}

func (c *CustomerQueryConsumer) OnDisabledV1(ctx context.Context, ev *event.CustomerStatusUpdatedV1) error {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryConsumer.OnDisabledV1")
  defer span.End()

  stat := c.svc.DisabledV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *CustomerQueryConsumer) OnEmailVerifiedV1(ctx context.Context, ev *event.CustomerEmailVerifiedV1) error {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryConsumer.OnEmailVerifiedV1")
  defer span.End()

  stat := c.svc.EmailVerifiedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *CustomerQueryConsumer) OnEnabledV1(ctx context.Context, ev *event.CustomerStatusUpdatedV1) error {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryConsumer.OnEnabledV1")
  defer span.End()

  stat := c.svc.EnabledV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *CustomerQueryConsumer) OnPasswordUpdatedV1(ctx context.Context, ev *event.CustomerPasswordUpdatedV1) error {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryConsumer.OnPasswordUpdatedV1")
  defer span.End()

  stat := c.svc.PasswordUpdatedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *CustomerQueryConsumer) OnPhotoUpdatedV1(ctx context.Context, ev *event.CustomerPhotoUpdatedV1) error {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryConsumer.OnPhotoUpdatedV1")
  defer span.End()

  stat := c.svc.PhotoUpdatedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *CustomerQueryConsumer) OnUpdatedV1(ctx context.Context, ev *event.CustomerUpdatedV1) error {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryConsumer.OnUpdatedV1")
  defer span.End()

  stat := c.svc.UpdatedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *CustomerQueryConsumer) OnVoucherAddedV1(ctx context.Context, ev *event.CustomerVoucherAddedV1) error {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryConsumer.OnVoucherAddedV1")
  defer span.End()

  stat := c.svc.VoucherAddedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *CustomerQueryConsumer) OnVoucherDeletedV1(ctx context.Context, ev *event.CustomerVoucherDeletedV1) error {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryConsumer.OnVoucherDeletedV1")
  defer span.End()

  stat := c.svc.VoucherDeletedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (c *CustomerQueryConsumer) OnVoucherUpdatedV1(ctx context.Context, ev *event.CustomerVoucherUpdatedV1) error {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryConsumer.OnVoucherUpdatedV1")
  defer span.End()

  stat := c.svc.VoucherUpdatedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}
