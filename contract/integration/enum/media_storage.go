package enum

type MediaUsage uint8

const (
  UsageOnce MediaUsage = iota
  UsageFull
)

func (m MediaUsage) Underlying() uint8 {
  return uint8(m)
}
