package pg

import (
  "context"
  "github.com/arcorium/rashop/services/mailer/internal/domain/entity"
  "github.com/arcorium/rashop/services/mailer/internal/domain/repository"
  "github.com/arcorium/rashop/services/mailer/internal/infra/model"
  "github.com/arcorium/rashop/services/mailer/pkg/tracer"
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  "github.com/arcorium/rashop/shared/util/repo"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "time"
)

func NewMail(db bun.IDB) repository.IMailPersistent {
  return &mailPersistent{
    db:     db,
    tracer: tracer.Get(),
  }
}

type mailPersistent struct {
  db     bun.IDB
  tracer trace.Tracer
}

func (m *mailPersistent) Get(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Mail], error) {
  ctx, span := m.tracer.Start(ctx, "mailPersistent.Get")
  defer span.End()

  var dbModels []model.Mail
  count, err := m.db.NewSelect().
    Model(&dbModels).
    Relation("Recipients").
    Order("created_at DESC").
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    ScanAndCount(ctx)

  result := repo.CheckPaginationResult(dbModels, count, err)
  if result.IsError() {
    err = result.Err()
    spanUtil.RecordError(err, span)
    return repo.PaginatedResult[entity.Mail]{Total: uint64(count)}, err
  }

  entities, ierr := sharedUtil.CastSliceErrsP(dbModels, repo.ToDomainErr[*model.Mail, entity.Mail])
  if ierr.IsError() {
    spanUtil.RecordError(ierr, span)
    return repo.PaginatedResult[entity.Mail]{Total: uint64(count)}, ierr
  }

  return repo.NewPaginatedResult(entities, uint64(count)), nil
}

func (m *mailPersistent) FindByIds(ctx context.Context, mailIds ...types.Id) ([]entity.Mail, error) {
  ctx, span := m.tracer.Start(ctx, "mailPersistent.FindByIds")
  defer span.End()

  input := sharedUtil.CastSlice(mailIds, func(id types.Id) repo.OrderedIdInput {
    return repo.OrderedIdInput{
      Id: id.String(),
    }
  })
  var dbModels []model.Mail
  err := m.db.NewSelect().
    With("result", m.db.NewValues(&input).WithOrder()).
    Model(&dbModels).
    Relation("Recipients").
    Join("JOIN result = result.input_id = m.id").
    Order("result._order").
    Scan(ctx)

  result := repo.CheckSliceResult(dbModels, err)
  if result.IsError() {
    err = result.Err()
    spanUtil.RecordError(err, span)
    return nil, err
  }

  entities, ierr := sharedUtil.CastSliceErrsP(dbModels, repo.ToDomainErr[*model.Mail, entity.Mail])
  if ierr.IsError() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }

  return entities, nil
}

func (m *mailPersistent) Create(ctx context.Context, mail *entity.Mail) error {
  ctx, span := m.tracer.Start(ctx, "mailPersistent.Create")
  defer span.End()

  err := m.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
    dbModel := model.FromMailDomain(mail)

    res, err := tx.NewInsert().
      Model(&dbModel).
      Returning("NULL").
      Exec(ctx)

    err = repo.CheckResult(res, err)
    if err != nil {
      return err
    }

    res, err = tx.NewInsert().
      Model(&dbModel.Recipients).
      Returning("NULL").
      Exec(ctx)
    return repo.CheckResult(res, err)
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }
  return nil
}

func (m *mailPersistent) Update(ctx context.Context, mail *entity.Mail) error {
  ctx, span := m.tracer.Start(ctx, "mailPersistent.Update")
  defer span.End()

  if mail.Updated() {
    dbModel := model.FromMailDomain(mail)

    res, err := m.db.NewUpdate().
      Model(&dbModel).
      WherePK().
      Returning("NULL").
      Exec(ctx)

    err = repo.CheckResult(res, err)
    if err != nil {
      spanUtil.RecordError(err, span)
      return err
    }
  }

  return nil
}

func (m *mailPersistent) Delete(ctx context.Context, startTime, endTime time.Time) (uint64, error) {
  ctx, span := m.tracer.Start(ctx, "mailPersistent.Delete")
  defer span.End()

  // Count
  query := m.db.NewSelect().
    Model(types.Nil[model.Mail]())

  deleteQuery := m.db.NewDelete().
    Model(types.Nil[model.Mail]())

  if !startTime.IsZero() {
    query = query.Where("created_at > ?", startTime)
    deleteQuery = deleteQuery.Where("created_at > ?", endTime)
  }
  if !endTime.IsZero() {
    query = query.Where("created_at < ?", endTime)
    deleteQuery = deleteQuery.Where("created_at < ?", startTime)
  }

  count, err := query.Count(ctx)
  err = repo.CheckCount(count, err)
  if err != nil {
    spanUtil.RecordError(err, span)
    return 0, err
  }

  res, err := deleteQuery.Exec(ctx)
  err = repo.CheckResult(res, err)
  if err != nil {
    spanUtil.RecordError(err, span)
    return 0, err
  }

  return uint64(count), nil
}
