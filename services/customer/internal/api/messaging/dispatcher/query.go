package dispatcher

import (
  "context"
  "errors"
  "fmt"
  "github.com/IBM/sarama"
  sharedErr "github.com/arcorium/rashop/shared/errors"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/logger"
  "github.com/arcorium/rashop/shared/messaging"
  "github.com/arcorium/rashop/shared/serde"
  "github.com/arcorium/rashop/shared/types"
  "github.com/avast/retry-go"
  "github.com/dnwe/otelsarama"
  "mini-shop/services/user/internal/api/messaging/consumer"
  "mini-shop/services/user/internal/domain/event"
  "mini-shop/services/user/pkg/util"
  "strconv"
  "time"
)

var (
  ErrUnexpectedVersion = errors.New("event has unexpected version")
)

const (
  maxSerializedMessageBuffer = 100
)

func NewQueryConsumerGroup(dlqPublisher messaging.IForwarder[*sarama.ConsumerMessage], customerConsumer consumer.ICustomerQueryHandler, deserializer serde.IDeserializerAny) *QueryConsumerGroupHandler {
  return &QueryConsumerGroupHandler{
    handler:      customerConsumer,
    deserializer: deserializer,
    errors:       make(chan errorMessage, maxSerializedMessageBuffer),
    dlqPublisher: dlqPublisher,
  }
}

type errorMessage struct {
  id      string
  event   string
  message *sarama.ConsumerMessage
  error   error
}

type QueryConsumerGroupHandler struct {
  handler      consumer.ICustomerQueryHandler
  deserializer serde.IDeserializerAny

  errors       chan errorMessage
  dlqPublisher messaging.IForwarder[*sarama.ConsumerMessage]
}

// Run as async
func (q *QueryConsumerGroupHandler) Run(ctx context.Context, group sarama.ConsumerGroup, topics ...string) {
  go func() {
    groupHandler := otelsarama.WrapConsumerGroupHandler(q)
    for {
      err := group.Consume(ctx, topics, groupHandler)

      if err != nil {
        if ctx.Err() != nil {
          return
        }
        logger.Infof("Failed to consume message from topic %s: %v", topics, err)
      }
      logger.Infof("Successfully consumed message from topic %s", topics)
    }
  }()
}

func (q *QueryConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
  go q.errorHandler()
  return nil
}

func (q *QueryConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
  close(q.errors)
  return nil
}

func (q *QueryConsumerGroupHandler) errorHandler() {
  for msg := range q.errors {
    logger.Infof("[%s] Error on processing %s::%s message, sent to DLQ", string(msg.message.Key), msg.event, msg.id)
    // Send to DLQ
    err := q.dlqPublisher.Forward(context.Background(), msg.message, msg.error)
    if err != nil {
      logger.Infof("Failed to forward message to DLQ: %v", err)
    }
  }
}

func (q *QueryConsumerGroupHandler) handleError(metadata *types.EventMetadata, message *sarama.ConsumerMessage, err error) {
  q.errors <- errorMessage{
    id:      metadata.Id,
    event:   metadata.Name,
    message: message,
    error:   err,
  }
}

func constructEventBase[E types.IEventBaseConstructable[T, V], T types.EventTyping, V types.EventVersioning](e E, id string, time time.Time) E {
  e.ConstructEventBase(types.WithId[T, V](id), types.WithTime[T, V](time))
  return e
}

func (q *QueryConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
  defer func() {
    if r := recover(); r != nil {
      logger.Infof("QueryConsumerGroupHandler::ConsumeClaim - Recovered from panic: %v", r)
    }
  }()

  logger.Infof("Start consuming %s on [%d] starting at %d offset", claim.Topic(), claim.Partition(), claim.InitialOffset())

  for message := range claim.Messages() {
    var err error
    logger.Infof("Got Message %s::%s ", message.Topic, message.Key)
    // Parse header
    md := types.NewEventMetadata(message.Headers, types.FromKafkaRecords())
    ctx := util.ExtractOTEL(md, message)

    switch md.Name {
    case event.CustomerCreatedEvent:
      err = dispatcherHandlerV1[*event.CustomerCreatedV1](q, ctx, message, constructEventBase(&event.CustomerCreatedV1{}, md.Id, message.Timestamp), q.handler.OnCreatedV1)
    case event.CustomerUpdatedEvent:
      err = dispatcherHandlerV1[*event.CustomerUpdatedV1](q, ctx, message, constructEventBase(&event.CustomerUpdatedV1{}, md.Id, message.Timestamp), q.handler.OnUpdatedV1)
    case event.CustomerStatusUpdatedEvent:
      err = dispatcherHandlerV1[*event.CustomerStatusUpdatedV1](q, ctx, message, constructEventBase(&event.CustomerStatusUpdatedV1{}, md.Id, message.Timestamp), q.handler.OnStatusUpdatedV1)
    case event.CustomerBalanceUpdatedEvent:
      err = dispatcherHandlerV1[*event.CustomerBalanceUpdatedV1](q, ctx, message, constructEventBase(&event.CustomerBalanceUpdatedV1{}, md.Id, message.Timestamp), q.handler.OnBalanceUpdatedV1)
    case event.CustomerPasswordUpdatedEvent:
      err = dispatcherHandlerV1[*event.CustomerPasswordUpdatedV1](q, ctx, message, constructEventBase(&event.CustomerPasswordUpdatedV1{}, md.Id, message.Timestamp), q.handler.OnPasswordUpdatedV1)
    case event.CustomerPhotoUpdatedEvent:
      err = dispatcherHandlerV1[*event.CustomerPhotoUpdatedV1](q, ctx, message, constructEventBase(&event.CustomerPhotoUpdatedV1{}, md.Id, message.Timestamp), q.handler.OnPhotoUpdatedV1)
    case event.CustomerEmailVerifiedEvent:
      err = dispatcherHandlerV1[*event.CustomerEmailVerifiedV1](q, ctx, message, constructEventBase(&event.CustomerEmailVerifiedV1{}, md.Id, message.Timestamp), q.handler.OnEmailVerifiedV1)
    case event.CustomerAddressAddedEvent:
      err = dispatcherHandlerV1[*event.CustomerAddressAddedV1](q, ctx, message, constructEventBase(&event.CustomerAddressAddedV1{}, md.Id, message.Timestamp), q.handler.OnAddressAddedV1)
    case event.CustomerAddressDeletedEvent:
      err = dispatcherHandlerV1[*event.CustomerAddressDeletedV1](q, ctx, message, constructEventBase(&event.CustomerAddressDeletedV1{}, md.Id, message.Timestamp), q.handler.OnAddressDeletedV1)
    case event.CustomerAddressUpdatedEvent:
      err = dispatcherHandlerV1[*event.CustomerAddressUpdatedV1](q, ctx, message, constructEventBase(&event.CustomerAddressUpdatedV1{}, md.Id, message.Timestamp), q.handler.OnAddressUpdatedV1)
    case event.CustomerVoucherAddedEvent:
      err = dispatcherHandlerV1[*event.CustomerVoucherAddedV1](q, ctx, message, constructEventBase(&event.CustomerVoucherAddedV1{}, md.Id, message.Timestamp), q.handler.OnVoucherAddedV1)
    case event.CustomerVoucherDeletedEvent:
      err = dispatcherHandlerV1[*event.CustomerVoucherDeletedV1](q, ctx, message, constructEventBase(&event.CustomerVoucherDeletedV1{}, md.Id, message.Timestamp), q.handler.OnVoucherDeletedV1)
    case event.CustomerVoucherUpdatedEvent:
      err = dispatcherHandlerV1[*event.CustomerVoucherUpdatedV1](q, ctx, message, constructEventBase(&event.CustomerVoucherUpdatedV1{}, md.Id, message.Timestamp), q.handler.OnVoucherUpdatedV1)
    case event.CustomerDefaultAddressUpdatedEvent:
      err = dispatcherHandlerV1[*event.CustomerDefaultAddressUpdatedV1](q, ctx, message, constructEventBase(&event.CustomerDefaultAddressUpdatedV1{}, md.Id, message.Timestamp), q.handler.OnDefaultAddressUpdatedV1)
    default:
      err = fmt.Errorf("unknown event : %s", md.Name)
    }

    if err != nil {
      // Sent to dead letter queue
      q.handleError(md, message, err)
    }
    session.MarkMessage(message, "")
  }

  return nil
}

// -- Helper

func deserializeMessage[T types.Event](deser serde.IDeserializerAny, message *sarama.ConsumerMessage, base T) error {
  // Parse value
  err := deser.Deserialize(message.Value, &base)
  return err
}

func dispatcherHandlerV1[E types.Event](q *QueryConsumerGroupHandler, ctx context.Context, message *sarama.ConsumerMessage, base E, handle handler.ConsumerFunc[E]) error {
  // Create base event
  // Version check
  if ver := base.EventVersion(); ver != (types.V1{}).EventVersion() {
    err := sharedErr.Wrap(ErrUnexpectedVersion, sharedErr.WithPrefix("V"+strconv.Itoa(int(ver))))
    return err
  }

  // Deserialize event
  err := deserializeMessage[E](q.deserializer, message, base)
  if err != nil {
    return err
  }

  logger.Infof("Deserialized Event '%s': %+v", base.EventName(), base)

  // Process it
  err = retry.Do(func() error {
    return handle(ctx, base)
  }, retry.Context(ctx), retry.OnRetry(inRetry(base.Identity(), base.EventName())))
  return err
}

func inRetry(id string, eventName string) retry.OnRetryFunc {
  return func(n uint, err error) {
    logger.Infof("[%d] Retrying on processing event %s::%s -> %s", n, eventName, id, err)
  }
}
