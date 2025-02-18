package util

func DoNothing(...any) {}

func CopyWith[T any](val T, f func(*T)) T {
  f(&val)
  return val
}

func CopyWithP[T any](val T, f func(*T)) *T {
  f(&val)
  return &val
}

func Clone[T any](val *T) T {
  return *val
}
