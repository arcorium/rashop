package errors

import "fmt"

type wrapBuilder func(err error) error

func WithPrefix(msg string) wrapBuilder {
  return func(err error) error {
    return fmt.Errorf("%s: %s", msg, err)
  }
}

func WithSuffix(msg string) wrapBuilder {
  return func(err error) error {
    return fmt.Errorf("%s %s", err, msg)
  }
}

func Wrap(err error, builders ...wrapBuilder) error {
  err2 := err
  for _, builder := range builders {
    err2 = builder(err2)
  }
  return err2
}
