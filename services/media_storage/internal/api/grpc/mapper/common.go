package mapper

import (
  mediav1 "github.com/arcorium/rashop/proto/gen/go/media_storage/v1"
  vob "github.com/arcorium/rashop/services/media_storage/internal/domain/valueobject"
)

func toProtoUsage(usage vob.MediaUsage) mediav1.Usage {
  switch usage {
  case vob.UsageOnce:
    return mediav1.Usage_OneTime
  case vob.UsageFull:
    return mediav1.Usage_Full
  }
  panic("unknown usage")
}

func ToDomainUsage(usage mediav1.Usage) vob.MediaUsage {
  switch usage {
  case mediav1.Usage_OneTime:
    return vob.UsageOnce
  case mediav1.Usage_Full:
    return vob.UsageFull
  }
  panic("unknown usage")
}
