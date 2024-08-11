package pg

import (
  "context"
  algo "github.com/arcorium/rashop/shared/algorithm"
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  "github.com/arcorium/rashop/shared/util/repo"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "mini-shop/services/user/internal/domain/entity"
  "mini-shop/services/user/internal/domain/repository"
  "mini-shop/services/user/internal/infra/model"
  "mini-shop/services/user/pkg/tracer"
)

func NewCustomer(db bun.IDB) repository.ICustomer {
  return &customerRepository{
    db:     db,
    tracer: tracer.Get(),
  }
}

type customerRepository struct {
  db     bun.IDB
  tracer trace.Tracer
}

func (c *customerRepository) Create(ctx context.Context, customer *entity.Customer) error {
  ctx, span := c.tracer.Start(ctx, "customerRepository.Create")
  defer span.End()

  dbModel := model.FromCustomerDomain(customer)
  err := c.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
    // Insert the user
    res, err := c.db.NewInsert().
      Model(dbModel.User).
      Returning("NULL").
      Exec(ctx)
    err = repo.CheckResult(res, err)
    if err != nil {
      return err
    }

    // Insert the customer
    res, err = c.db.NewInsert().
      Model(&dbModel).
      Returning("NULL").
      Exec(ctx)
    if err != nil {
      return err
    }
    return nil
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }
  return nil
}

func (c *customerRepository) addAddress(ctx context.Context, db bun.IDB, addresses ...model.ShippingAddress) error {
  ctx, span := c.tracer.Start(ctx, "customerRepository.addAddress")
  defer span.End()

  res, err := db.NewInsert().
    Model(&addresses).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResult(res, err)
}

func (c *customerRepository) addVouchers(ctx context.Context, db bun.IDB, vouchers ...model.Voucher) error {
  ctx, span := c.tracer.Start(ctx, "customerRepository.addVouchers")
  defer span.End()

  res, err := db.NewInsert().
    Model(&vouchers).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResult(res, err)
}

func (c *customerRepository) Get(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Customer], error) {
  ctx, span := c.tracer.Start(ctx, "customerRepository.Get")
  defer span.End()

  var dbModels []model.Customer
  count, err := c.db.NewSelect().
    Model(&dbModels).
    Relation("User").
    Relation("Vouchers").
    Relation("ShippingAddresses").
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    Order("u.username").
    ScanAndCount(ctx, &dbModels)

  res := repo.CheckPaginationResult(dbModels, count, err)
  if res.IsError() {
    spanUtil.RecordError(err, span)
    return repo.NewPaginatedResult[entity.Customer](nil, uint64(count)), res.Err()
  }

  entities, ierr := sharedUtil.CastSliceErrsP(dbModels, repo.ToDomainErr[*model.Customer, entity.Customer])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.Customer](nil, uint64(count)), res.Err()
  }

  return repo.NewPaginatedResult(entities, uint64(count)), nil
}

func (c *customerRepository) FindByIds(ctx context.Context, ids ...types.Id) ([]entity.Customer, error) {
  ctx, span := c.tracer.Start(ctx, "customerRepository.FindByIds")
  defer span.End()

  input := sharedUtil.CastSlice(ids, func(id types.Id) repo.OrderedIdInput {
    return repo.OrderedIdInput{
      Id: id.String(),
    }
  })
  var dbModels []model.Customer
  err := c.db.NewSelect().
    With("result", c.db.NewValues(&input).WithOrder()).
    Model(&dbModels).
    Relation("User").
    Relation("Vouchers").
    Relation("ShippingAddresses").
    Join("JOIN result ON result.input_id = c.id").
    Order("result._order").
    Scan(ctx)

  result := repo.CheckSliceResult(dbModels, err)
  if result.IsError() {
    spanUtil.RecordError(result.Err(), span)
    return nil, err
  }

  entities, ierr := sharedUtil.CastSliceErrsP(dbModels, repo.ToDomainErr[*model.Customer, entity.Customer])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }

  return entities, nil
}

func (c *customerRepository) FindByEmails(ctx context.Context, emails ...types.Email) ([]entity.Customer, error) {
  ctx, span := c.tracer.Start(ctx, "customerRepository.FindByEmails")
  defer span.End()

  type orderedEmailInput struct {
    Email string `bun:"input_email"`
  }

  input := sharedUtil.CastSlice(emails, func(email types.Email) orderedEmailInput {
    return orderedEmailInput{
      Email: email.String(),
    }
  })
  var dbModels []model.Customer
  err := c.db.NewSelect().
    With("result", c.db.NewValues(&input).WithOrder()).
    Model(&dbModels).
    Relation("User").
    Relation("Vouchers").
    Relation("ShippingAddresses").
    Order("result._order").
    Join("JOIN result ON u.email = result.input_email").
    Scan(ctx)

  result := repo.CheckSliceResult(dbModels, err)
  if result.IsError() {
    spanUtil.RecordError(result.Err(), span)
    return nil, err
  }

  entities, ierr := sharedUtil.CastSliceErrsP(dbModels, repo.ToDomainErr[*model.Customer, entity.Customer])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }

  return entities, nil
}

func (c *customerRepository) updateCustomer(ctx context.Context, db bun.IDB, customer *model.Customer) error {
  ctx, span := c.tracer.Start(ctx, "customerRepository.updateCustomer")
  defer span.End()

  res, err := db.NewUpdate().
    Model(&customer).
    WherePK().
    Returning("NULL").
    Exec(ctx)

  if err = repo.CheckResult(res, err); err != nil {
    return err
  }

  res, err = db.NewUpdate().
    Model(customer.User).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResult(res, err)
}

func (c *customerRepository) Update(ctx context.Context, customer *entity.Customer) error {
  ctx, span := c.tracer.Start(ctx, "customerRepository.Update")
  defer span.End()

  err := c.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
    dbModel := model.FromCustomerDomain(customer)
    // Update the customer and user
    if customer.Updated() {
      err := c.updateCustomer(ctx, tx, &dbModel)
      if err != nil {
        return err
      }
    }

    // Shipping addresses
    // Insert
    if total := customer.ShippingAddresses.Added(); total > 0 {
      err := c.addAddress(ctx, tx, dbModel.ShippingAddresses[len(dbModel.ShippingAddresses)-int(total):]...)
      if err != nil {
        return err
      }
    }
    // Update
    if indices := customer.ShippingAddresses.Updated(); indices != nil {
      views := algo.FilterByIndicesPointing(dbModel.ShippingAddresses, indices...)
      err := c.updateAddress(ctx, tx, views...)
      if err != nil {
        return err
      }
    }
    // Delete
    if ids := customer.ShippingAddresses.Deleted(); ids != nil {
      err := c.deleteAddresses(ctx, tx, customer.Id, ids...)
      if err != nil {
        return err
      }
    }

    // Vouchers
    // Insert
    if total := customer.Vouchers.Added(); total > 0 {
      err := c.addVouchers(ctx, tx, dbModel.Vouchers[len(dbModel.Vouchers)-int(total):]...)
      if err != nil {
        return err
      }
    }
    // Update
    if indices := customer.Vouchers.Updated(); indices != nil {
      views := algo.FilterByIndicesPointing(dbModel.Vouchers, indices...)
      err := c.updateVouchers(ctx, tx, views...)
      if err != nil {
        return err
      }
    }
    // Delete
    if ids := customer.Vouchers.Deleted(); ids != nil {
      err := c.deleteVouchers(ctx, tx, customer.Id, ids...)
      if err != nil {
        return err
      }
    }

    return nil
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }
  return nil
}

func (c *customerRepository) updateAddress(ctx context.Context, db bun.IDB, addresses ...*model.ShippingAddress) error {
  ctx, span := c.tracer.Start(ctx, "customerRepository.updateAddress")
  defer span.End()

  res, err := db.NewUpdate().
    Model(&addresses).
    WherePK().
    Bulk().
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (c *customerRepository) updateVouchers(ctx context.Context, db bun.IDB, vouchers ...*model.Voucher) error {
  ctx, span := c.tracer.Start(ctx, "customerRepository.updateVouchers")
  defer span.End()

  res, err := db.NewUpdate().
    Model(&vouchers).
    WherePK().
    Bulk().
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResult(res, err)
}

func (c *customerRepository) deleteAddresses(ctx context.Context, db bun.IDB, customerId types.Id, addressIds ...types.Id) error {
  ctx, span := c.tracer.Start(ctx, "customerRepository.deleteAddresses")
  defer span.End()

  ids := sharedUtil.CastSlice(addressIds, sharedUtil.ToString[types.Id])
  query := db.NewDelete().
    Model(types.Nil[model.ShippingAddress]()).
    Where("user_id = ?", customerId.String())

  if len(addressIds) > 0 {
    query = query.Where("id IN (?) ", customerId.String(), bun.In(ids))
  }

  res, err := query.Exec(ctx)
  return repo.CheckResult(res, err)
}

func (c *customerRepository) deleteVouchers(ctx context.Context, db bun.IDB, customerId types.Id, voucherIds ...types.Id) error {
  ctx, span := c.tracer.Start(ctx, "customerRepository.deleteVouchers")
  defer span.End()

  ids := sharedUtil.CastSlice(voucherIds, sharedUtil.ToString[types.Id])
  query := db.NewDelete().
    Model(types.Nil[model.Voucher]()).
    Where("user_id = ?", customerId.String())

  if len(voucherIds) > 0 {
    query = query.Where("voucher_id IN (?) ", customerId.String(), bun.In(ids))
  }

  res, err := query.Exec(ctx)
  return repo.CheckResult(res, err)
}
