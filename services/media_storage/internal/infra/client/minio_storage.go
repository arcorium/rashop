package client

import (
  "bytes"
  "context"
  "crypto/sha256"
  "fmt"
  "github.com/arcorium/rashop/services/media_storage/internal/domain/entity"
  "github.com/arcorium/rashop/services/media_storage/internal/domain/repository"
  vob "github.com/arcorium/rashop/services/media_storage/internal/domain/valueobject"
  "github.com/arcorium/rashop/services/media_storage/pkg/tracer"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "github.com/minio/minio-go/v7"
  "github.com/minio/minio-go/v7/pkg/credentials"
  "go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
  "go.opentelemetry.io/otel/trace"
  "io"
  "path/filepath"
  "time"
)

const policy = `
    {
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": {
                    "AWS": [
                        "*"
                    ]
                },
                "Action": [
                    "s3:GetObject"
                ],
                "Resource": [
                    "arn:aws:s3:::%s/public/*"
                ]
            }
        ]
    }`

func upsertBucket(ctx context.Context, client *minio.Client, bucket string) error {
  exists, err := client.BucketExists(ctx, bucket)
  if err != nil {
    return err
  }
  if !exists {
    // Create bucket
    err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
    if err != nil {
      return err
    }

    // Set policy for /public
    err = client.SetBucketPolicy(ctx, bucket, fmt.Sprintf(policy, bucket))
    if err != nil {
      return err
    }
  }
  return nil
}

func NewMinIOStorageClient(isSecure bool, bucket string, config MinIOStorageConfig) (repository.IStorageClient, error) {
  // Using context.Background to use global tracer
  ctx := context.Background()
  clientTrace := otelhttptrace.NewClientTrace(ctx)
  client, err := minio.New(config.Address, &minio.Options{
    Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
    Secure: isSecure,
    Trace:  clientTrace,
  })
  if err != nil {
    return nil, err
  }

  err = upsertBucket(ctx, client, bucket)
  if err != nil {
    return nil, err
  }

  return &minIOStorageClient{
    bucket: bucket,
    config: config,
    client: client,
    tracer: tracer.Get(),
  }, nil
}

type MinIOStorageConfig struct {
  Address         string
  AccessKeyID     string
  SecretAccessKey string
}

type minIOStorageClient struct {
  bucket string

  config MinIOStorageConfig
  client *minio.Client
  tracer trace.Tracer
}

func (m *minIOStorageClient) getObjectName(media *entity.Media) string {
  // Hash the id for object name
  filename := media.Id.Hash(sha256.New())
  // Add extension to it
  ext := filepath.Ext(media.Name)
  filename += ext

  if media.IsPublic {
    return fmt.Sprintf("public/%s", filename)
  }
  // Disallow public prefix on non-public media
  //return strings.TrimPrefix(filename, "public/")
  return filename
}

func (m *minIOStorageClient) Get(ctx context.Context, media *entity.Media) (*entity.Media, error) {
  ctx, span := m.tracer.Start(ctx, "minIOStorageClient.Get")
  defer span.End()

  if media.Provider != m.GetProvider() {
    err := fmt.Errorf("media is stored in different provider")
    spanUtil.RecordError(err, span)
    return nil, err
  }

  obj, err := m.client.GetObject(ctx, m.bucket, media.ProviderPath, minio.GetObjectOptions{})
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  bytes, err := io.ReadAll(obj)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  media.Data = bytes
  return media, nil
}

func (m *minIOStorageClient) Store(ctx context.Context, media *entity.Media) error {
  ctx, span := m.tracer.Start(ctx, "minIOStorageClient.Store")
  defer span.End()

  reader := bytes.NewReader(media.Data)
  objectName := m.getObjectName(media)

  opt := minio.PutObjectOptions{
    ContentType: media.ContentType,
  }

  // Timed object
  //if media.Usage.Type == vob.UsageTimed && !media.Usage.Time.IsZero() {
  //  opt.Expires = media.Usage.Time
  //}

  info, err := m.client.PutObject(ctx,
    m.bucket,
    objectName,
    reader,
    int64(media.Size),
    opt,
  )
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }

  // Modify media
  media.ProviderPath = info.Key
  //media.FullPath, err = m.GetFullPath(ctx, objectName, media.IsPublic)
  return nil
}

func (m *minIOStorageClient) Delete(ctx context.Context, name string) error {
  ctx, span := m.tracer.Start(ctx, "minIOStorageClient.Delete")
  defer span.End()

  err := m.client.RemoveObject(ctx, m.bucket, name, minio.RemoveObjectOptions{
    ForceDelete: true,
  })
  return err
}

func (m *minIOStorageClient) GetFullPath(ctx context.Context, name string, public bool) (string, error) {
  ctx, span := m.tracer.Start(ctx, "minIOStorageClient.GetFullPath")
  defer span.End()

  // TODO: For public it should get the object, because the name is permanent
  stat, err := m.client.StatObject(ctx, m.bucket, name, minio.StatObjectOptions{})
  if err != nil {
    return "", err
  }
  sharedUtil.DoNothing(stat)

  url, err := m.client.PresignedGetObject(ctx, m.bucket, name, time.Hour*24, nil)
  if err != nil {
    return "", err
  }
  return url.String(), nil
}

func (m *minIOStorageClient) GetProvider() vob.StorageProvider {
  return vob.ProviderMinIO
}
