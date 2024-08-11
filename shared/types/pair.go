package types

func NewPair[T, U any](first T, second U) Pair[T, U] {
  return Pair[T, U]{
    First:  first,
    Second: second,
  }
}

func NewKeyVal[T, U any](key T, val U) KeyVal[T, U] {
  return KeyVal[T, U]{
    Key: key,
    Val: val,
  }
}

type Pair[T, U any] struct {
  First  T
  Second U
}

type KeyVal[T, U any] struct {
  Key T
  Val U
}
