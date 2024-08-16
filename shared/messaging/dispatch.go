package messaging

import (
  "context"
  "errors"
  "github.com/IBM/sarama"
  sharedErr "github.com/arcorium/rashop/shared/errors"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/logger"
  "github.com/arcorium/rashop/shared/serde"
  "github.com/arcorium/rashop/shared/types"
  "github.com/avast/retry-go"
  "strconv"
)

var (
  ErrUnexpectedVersion = errors.New("event has unexpected version")
)

func Dispatch[E types.Event, V types.EventVersioning](deser serde.IDeserializerAny, ctx context.Context, message *sarama.ConsumerMessage, base E, handle handler.ConsumerFunc[E]) error {
  // Version check
  var v V
  if ver := base.EventVersion(); ver != v.EventVersion() {
    err := sharedErr.Wrap(ErrUnexpectedVersion, sharedErr.WithPrefix("V"+strconv.Itoa(int(ver))))
    return err
  }

  // Deserialize event
  err := deser.Deserialize(message.Value, &base)
  if err != nil {
    return err
  }

  logger.Infof("Deserialized Event '%s': %+v", base.EventName(), base)

  // Process it
  err = retry.Do(func() error {
    defer func() {
      if r := recover(); r != nil {
        logger.Infof("DispatchV1 - Recovered from panic: %v", r)
      }
    }()

    return handle(ctx, base)
  }, retry.Context(ctx), retry.OnRetry(inRetry(base.Identity(), base.EventName())))
  return err
}

func DispatchV1[E types.Event](deser serde.IDeserializerAny, ctx context.Context, message *sarama.ConsumerMessage, base E, handle handler.ConsumerFunc[E]) error {
  return Dispatch[E, types.V1](deser, ctx, message, base, handle)
}

//func DispatchV1[E types.Event](deser serde.IDeserializerAny, ctx context.Context, message *sarama.ConsumerMessage, base E, handle handler.ConsumerFunc[E]) error {
//  // Version check
//  if ver := base.EventVersion(); ver != (types.V1{}).EventVersion() {
//    err := sharedErr.Wrap(ErrUnexpectedVersion, sharedErr.WithPrefix("V"+strconv.Itoa(int(ver))))
//    return err
//  }
//
//  // Deserialize event
//  err := deser.Deserialize(message.Value, &base)
//  if err != nil {
//    return err
//  }
//
//  logger.Infof("Deserialized Event '%s': %+v", base.EventName(), base)
//
//  // Process it
//  err = retry.Do(func() error {
//    defer func() {
//      if r := recover(); r != nil {
//        logger.Infof("DispatchV1 - Recovered from panic: %v", r)
//      }
//    }()
//
//    return handle(ctx, base)
//  }, retry.Context(ctx), retry.OnRetry(inRetry(base.Identity(), base.EventName())))
//  return err
//}

func inRetry(id string, eventName string) retry.OnRetryFunc {
  return func(n uint, err error) {
    logger.Infof("[%d] Retrying on processing event %s::%s -> %s", n, eventName, id, err)
  }
}
