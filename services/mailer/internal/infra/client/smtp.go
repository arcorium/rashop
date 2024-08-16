package client

import (
  "context"
  "crypto/tls"
  "github.com/arcorium/rashop/services/mailer/internal/domain/entity"
  "github.com/arcorium/rashop/services/mailer/internal/domain/repository"
  vob "github.com/arcorium/rashop/services/mailer/internal/domain/valueobject"
  "github.com/arcorium/rashop/services/mailer/pkg/tracer"
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  "go.opentelemetry.io/otel/trace"
  "gopkg.in/gomail.v2"
  "io"
)

const defaultMaxClient = 3

func NewSMTP(config SMTPConfig) (repository.IMailClient, error) {
  total := config.MaxClients
  if total == 0 {
    total = defaultMaxClient
  }

  pool, err := newSMTPClientPool(&config, total)
  if err != nil {
    return nil, err
  }

  return &smtpMailClient{
    pool:   pool,
    tracer: tracer.Get(),
  }, nil
}

type SMTPConfig struct {
  Host      string
  Port      uint16
  Username  string
  Password  string
  TLSConfig *tls.Config

  MaxClients uint
}

type smtpMailClient struct {
  pool   *smtpClientPool
  tracer trace.Tracer
}

func (s *smtpMailClient) Send(ctx context.Context, mail *entity.Mail, attachments []vob.Attachment) error {
  ctx, span := s.tracer.Start(ctx, "smtpMailClient.Send")
  defer span.End()

  msg := gomail.NewMessage(func(m *gomail.Message) {
    m.SetHeader("Subject", mail.Subject)
    m.SetBody(mail.BodyType.String(), mail.Body)

    for _, attachment := range attachments {
      if attachment.IsEmbedded {
        m.Embed(attachment.Filename, gomail.SetCopyFunc(copyFunc(&attachment)))
      } else {
        m.Attach(attachment.Filename, gomail.SetCopyFunc(copyFunc(&attachment)))
      }
    }
  })

  to := sharedUtil.CastSlice(mail.Recipients, sharedUtil.ToString[types.Email])
  return s.pool.Do(ctx, func(ctx context.Context, sender gomail.SendCloser) error {
    return sender.Send(mail.Sender.String(), to, msg)
  })
}

func copyFunc(file *vob.Attachment) func(io.Writer) error {
  return func(writer io.Writer) error {
    _, err := writer.Write(file.Data)
    return err
  }
}

func (s *smtpMailClient) Close(ctx context.Context) error {
  return s.pool.Close()
}
