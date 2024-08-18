package common

import (
  "crypto/sha512"
  vob "github.com/arcorium/rashop/services/token/internal/domain/valueobject"
  "github.com/arcorium/rashop/shared/types"
  rand2 "math/rand/v2"
)

const (
  defaultPinLength   = 6
  alphanumericString = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
  fullString         = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func GenerateToken(tokenType vob.TokenType, length uint16) string {
  switch tokenType {
  case vob.TypeString:
    return generateStringToken(length)
  case vob.TypePIN:
    return generatePinToken(false, length)
  case vob.TypeAlphanumericPIN:
    return generatePinToken(true, length)
  }
  return ""
}

func generateStringToken(length uint16) string {
  if length == 0 {
    return types.MustCreateId().Hash(sha512.New())
  }

  result := generateToken(length, len(fullString), fullString)
  return string(result)
}

func generatePinToken(alphanumeric bool, length uint16) string {
  if length == 0 {
    length = defaultPinLength
  }

  count := 10
  if alphanumeric {
    count = len(alphanumericString)
  }

  result := generateToken(length, count, alphanumericString)
  return string(result)
}

func generateToken(length uint16, count int, base string) []byte {
  result := make([]byte, length)
  for i := range length {
    n := rand2.N(count)
    result[i] = base[n]
  }
  return result
}
