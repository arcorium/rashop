package enum

type TokenUsage uint8

const (
  UsageEmailVerification TokenUsage = iota
  UsageResetPassword
  UsageLogin
  UsageGeneral
  UsageUnknown
)
