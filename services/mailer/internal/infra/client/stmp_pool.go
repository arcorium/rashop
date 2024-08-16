package client

import (
  "context"
  "errors"
  "gopkg.in/gomail.v2"
  "sync"
)

func newSMTPClientPool(config *SMTPConfig, total uint) (*smtpClientPool, error) {
  wg := sync.WaitGroup{}
  result := &smtpClientPool{}
  lock := sync.Mutex{}

  // Use goroutine to make multiple connections concurrently
  for range total {
    wg.Add(1)
    go func() {
      defer wg.Done()

      dialer := gomail.NewDialer(config.Host, int(config.Port), config.Username, config.Password)
      dialer.TLSConfig = config.TLSConfig

      sender, err := dialer.Dial()
      if err != nil {
        return
      }

      lock.Lock()
      result.Put(sender)
      lock.Unlock()
    }()
  }

  wg.Wait()
  return result, nil
}

type smtpClientPool struct {
  maxClient uint
  pool      sync.Pool
  cond      sync.Cond
}

// Get client and wait for it until one client available
func (c *smtpClientPool) Get() gomail.SendCloser {
  data := c.TryGet()
  if data == nil {
    data = c.wait()
  }
  return data
}

// TryGet works like Get, instead of waiting it will return whatever it is.
// Caller should check either the return value is nil or not
func (c *smtpClientPool) TryGet() gomail.SendCloser {
  return c.pool.Get().(gomail.SendCloser)
}

func (c *smtpClientPool) wait() gomail.SendCloser {
  c.cond.L.Lock()
  var data gomail.SendCloser
  for data == nil {
    data = c.TryGet()
    c.cond.Wait()
  }
  c.cond.L.Unlock()

  return data
}

// Put store client into pool
func (c *smtpClientPool) Put(dialer gomail.SendCloser) {
  c.pool.Put(dialer)
}

func (c *smtpClientPool) Do(ctx context.Context, f func(context.Context, gomail.SendCloser) error) error {
  dialer := c.Get()
  defer c.Put(dialer)
  return f(ctx, dialer)
}

func (c *smtpClientPool) Close() error {
  var err error
  for range c.maxClient {
    client := c.Get()
    if client == nil {
      panic("Bad implementation of SMTP client, client and pool have different length")
    }
    err2 := client.Close()
    if err2 != nil {
      err = errors.Join(err, err2)
    }
  }
  return err
}
