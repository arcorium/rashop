package persistence

import (
  "context"
  "github.com/arcorium/rashop/services/media_storage/internal/domain/entity"
  "github.com/arcorium/rashop/services/media_storage/internal/domain/repository"
  "github.com/arcorium/rashop/services/media_storage/internal/infra/model"
  "github.com/arcorium/rashop/services/media_storage/pkg/tracer"
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  "github.com/arcorium/rashop/shared/util/repo"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
)

func NewMediaPostgres(db bun.IDB) repository.IMediaPersistent {
  return &mediaPostgresPersistent{
    db:     db,
    tracer: tracer.Get(),
  }
}

type mediaPostgresPersistent struct {
  db     bun.IDB
  tracer trace.Tracer
}

func (m *mediaPostgresPersistent) AsUnit(ctx context.Context, f repo.UOWBlock[repository.IMediaPersistent]) error {
  ctx, span := m.tracer.Start(ctx, "mediaPostgresPersistent.AsUnit")
  defer span.End()

  err := m.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
    return f(ctx, m.recreate(tx))
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }

  return nil
}

func (m *mediaPostgresPersistent) recreate(db bun.IDB) repository.IMediaPersistent {
  return NewMediaPostgres(db)
}

func (m *mediaPostgresPersistent) Get(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Media], error) {
  ctx, span := m.tracer.Start(ctx, "mediaPostgresPersistent.Get")
  defer span.End()

  var dbModels []model.Metadata
  count, err := m.db.NewSelect().
    Model(&dbModels).
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    Order("created_at DESC").
    ScanAndCount(ctx)

  result := repo.CheckPaginationResult(dbModels, count, err)
  if result.IsError() {
    err = result.Err()
    spanUtil.RecordError(err, span)
    return repo.NewPaginatedResult[entity.Media](nil, uint64(count)), err
  }

  // Cast
  entities, ierr := sharedUtil.CastSliceErrsP(result.Data, repo.ToDomainErr[*model.Metadata])
  if ierr.IsError() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.Media](nil, uint64(count)), ierr
  }

  return repo.NewPaginatedResult[entity.Media](entities, uint64(count)), nil
}

func (m *mediaPostgresPersistent) FindByIds(ctx context.Context, mediaIds ...types.Id) ([]entity.Media, error) {
  ctx, span := m.tracer.Start(ctx, "mediaPostgresPersistent.GetMetadata")
  defer span.End()

  input := sharedUtil.CastSlice(mediaIds, func(id types.Id) repo.OrderedIdInput {
    return repo.OrderedIdInput{
      Id: id.String(),
    }
  })

  var dbModels []model.Metadata
  err := m.db.NewSelect().
    With("result", m.db.NewValues(&input).WithOrder()).
    Model(&dbModels).
    Join("JOIN result = result.input_id = mm.id").
    Order("result._order").
    Scan(ctx)

  result := repo.CheckSliceResult(dbModels, err)
  if result.IsError() {
    err = result.Err()
    spanUtil.RecordError(err, span)
    return nil, err
  }

  entities, ierr := sharedUtil.CastSliceErrsP(dbModels, repo.ToDomainErr[*model.Metadata])
  if ierr.IsError() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }
  return entities, nil
}

func (m *mediaPostgresPersistent) Create(ctx context.Context, media *entity.Media) error {
  ctx, span := m.tracer.Start(ctx, "mediaPostgresPersistent.Create")
  defer span.End()

  dbModel := model.FromMediaDomain(media)
  res, err := m.db.NewInsert().
    Model(&dbModel).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (m *mediaPostgresPersistent) Update(ctx context.Context, media *entity.Media) error {
  ctx, span := m.tracer.Start(ctx, "mediaPostgresPersistent.Update")
  defer span.End()
  if media.Updated() {
    dbModel := model.FromMediaDomain(media)

    res, err := m.db.NewUpdate().
      Model(&dbModel).
      WherePK().
      Exec(ctx)

    return repo.CheckResultWithSpan(res, err, span)
  }
  return nil
}

func (m *mediaPostgresPersistent) Delete(ctx context.Context, mediaIds ...types.Id) error {
  ctx, span := m.tracer.Start(ctx, "mediaPostgresPersistent.Delete")
  defer span.End()

  ids := sharedUtil.CastSlice(mediaIds, sharedUtil.ToString[types.Id])
  res, err := m.db.NewDelete().
    Model(types.Nil[model.Metadata]()).
    Where("id IN (?)", bun.In(ids)).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}
