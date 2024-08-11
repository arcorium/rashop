package types

type Preservable[R any] interface {
  PreserveFields() R
  RestorePreserved(*R)
}

func InPreservedUnaryTx[T Preservable[R], R any](val T, f func(T)) {
  preserved := val.PreserveFields()
  f(val)
  val.RestorePreserved(&preserved)
}

func InPreservedTx[T Preservable[R], R any](val T, f func()) {
  preserved := val.PreserveFields()
  f()
  val.RestorePreserved(&preserved)
}
