package types

func Some[T any](data T, err error) Result[T] {
  return Result[T]{
    Data: data,
    err:  err,
  }
}

func None[T any](err error) Result[T] {
  return Result[T]{
    err: err,
  }
}

func NoneF[T any](none func() T, err error) Result[T] {
  return Result[T]{
    Data: none(),
    err:  err,
  }
}

func SomeF[T any](f func() (T, error)) Result[T] {
  return Some(f())
}

func DropError[T any](val T, err error) T {
  return Some(val, err).Data
}

func Must[T any](val T, err error) T {
  res := Some(val, err)
  if res.IsError() {
    panic(err)
  }
  return res.Data
}

func SomeF1[T, P1 any](f func(P1) (T, error), param *P1) Result[T] {
  return Some(f(*param))
}

type Result[T any] struct {
  Data T
  err  error
}

func (r *Result[T]) Value() (T, error) {
  return r.Data, r.err
}

func (r *Result[T]) IsError() bool {
  return r.err != nil
}

func (r *Result[T]) Err() error {
  return r.err
}
