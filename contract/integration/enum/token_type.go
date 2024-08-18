package enum

type TokenType uint8

const (
  TypeString TokenType = iota
  TypePIN
  TypeAlphanumericPIN
  TypeUnknown
)
