package dispatcher

import (
  "context"
  "errors"
  "github.com/IBM/sarama"
  sharedErr "github.com/arcorium/rashop/shared/errors"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/logger"
  "github.com/arcorium/rashop/shared/serde"
  "github.com/arcorium/rashop/shared/types"
  "github.com/dnwe/otelsarama"
  "mini-shop/services/user/internal/api/messaging/consumer"
  "mini-shop/services/user/internal/domain/event"
  "mini-shop/services/user/pkg/util"
  "strconv"
  "sync"
  "sync/atomic"
)

var (
  ErrUnexpectedVersion = errors.New("event has unexpected version")
)

func NewQueryConsumerGroup(errBufferLimit int, customerConsumer consumer.ICustomerQueryHandler, deserializer serde.IDeserializerAny) sarama.ConsumerGroupHandler {
  return otelsarama.WrapConsumerGroupHandler(&QueryConsumerGroupHandler{
    handler:        customerConsumer,
    deserializer:   deserializer,
    concurrentLock: sync.RWMutex{},
    concurrent:     make(map[string]*sync.Mutex),
    errors:         make(chan errorMessage, errBufferLimit),
    exitChan:       make(chan struct{}),
  })
}

type errorMessage struct {
  MessageId string
  EventName string
  Error     error
}

type QueryConsumerGroupHandler struct {
  handler      consumer.ICustomerQueryHandler
  deserializer serde.IDeserializerAny

  concurrentLock sync.RWMutex
  concurrent     map[string]*sync.Mutex // Concurrent access per aggregate
  errors         chan errorMessage
  exitChan       chan struct{}
}

func (q *QueryConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
  // Run goroutine to handle error
  go q.errorHandler()
  return nil
}

func (q *QueryConsumerGroupHandler) errorHandler() {
  var shouldExit atomic.Bool

out:
  for {
    select {
    case <-q.exitChan:
      if len(q.errors) == 0 {
        break
      }
      shouldExit.Store(true)
      // Wait until all errors processed
      continue
    case err := <-q.errors:
      // Process
      logger.Infof("Failed to process message %s [%s]: %s", err.EventName, err.MessageId, err.Error)
      if shouldExit.Load() && len(q.errors) == 0 {
        break out
      }
    }
  }
}

func (q *QueryConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
  close(q.exitChan)
  close(q.errors)
  return nil
}

func (q *QueryConsumerGroupHandler) getLock(id string) *sync.Mutex {
  q.concurrentLock.RLock()
  lock, ok := q.concurrent[id]
  q.concurrentLock.RUnlock()
  if !ok {
    q.concurrentLock.Lock()
    lock = &sync.Mutex{}
    q.concurrent[id] = lock
    q.concurrentLock.Unlock()
  }
  return lock
}

func (q *QueryConsumerGroupHandler) sendError(id, eventName string, err error) {
  q.errors <- errorMessage{
    MessageId: id,
    EventName: eventName,
    Error:     err,
  }
}

func (q *QueryConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
  for message := range claim.Messages() {
    // Parse header
    md := types.NewEventMetadata(message.Headers, types.FromKafkaRecords())
    ctx := util.ExtractOTEL(md, message)
    lock := q.getLock(md.Id)
    lock.Lock()

    switch message.Topic {
    case event.CustomerCreatedEvent:
      go dispatcherHandlerV1[*event.CustomerCreatedV1](q, ctx, lock, message, func() *event.CustomerCreatedV1 {
        return &event.CustomerCreatedV1{
          DomainEventBaseV1: types.NewDomainEventV1(types.WithId(md.Id), types.WithTime(message.Timestamp)),
        }
      }, q.handler.OnCreatedV1)
    case event.CustomerUpdatedEvent:
      go dispatcherHandlerV1[*event.CustomerUpdatedV1](q, ctx, lock, message, func() *event.CustomerUpdatedV1 {
        return &event.CustomerUpdatedV1{
          DomainEventBaseV1: types.NewDomainEventV1(types.WithId(md.Id), types.WithTime(message.Timestamp)),
        }
      }, q.handler.OnUpdatedV1)
    case event.CustomerStatusUpdatedEvent:
    case event.CustomerBalanceUpdatedEvent:
    case event.CustomerPasswordUpdatedEvent:
    case event.CustomerPhotoUpdatedEvent:
    case event.CustomerEmailVerifiedEvent:
    case event.CustomerDeletedEvent:
    case event.CustomerForgotPasswordRequestedEvent:
    case event.CustomerEmailVerificationRequestedEvent:
    case event.CustomerAddressAddedEvent:
    case event.CustomerAddressesDeletedEvent:
    case event.CustomerAddressUpdatedEvent:
    case event.CustomerVouchersAddedEvent:
    case event.CustomerVouchersDeletedEvent:
    case event.CustomerVoucherUpdatedEvent:
    case event.CustomerDefaultAddressUpdatedEvent:
    }
  }

  return nil
}

// -- Helper

func deserializeMessage[T types.Event](deser serde.IDeserializerAny, message *sarama.ConsumerMessage, base T) error {
  // Parse value
  err := deser.Deserialize(message.Value, base)
  return err
}

func dispatcherHandlerV1[E types.Event](q *QueryConsumerGroupHandler, ctx context.Context, lock *sync.Mutex, message *sarama.ConsumerMessage, eventSetupFunc func() E, handle handler.ConsumerFunc[E]) {
  md := ctx.Value(types.EventMetadataCtxKey{}).(*types.EventMetadata)
  defer lock.Unlock()

  // Create base event
  ev := eventSetupFunc()

  // Version check
  if ver := ev.EventVersion(); ver != (types.V1{}).EventVersion() {
    err := sharedErr.Wrap(ErrUnexpectedVersion, sharedErr.WithPrefix("V"+strconv.Itoa(int(ver))))
    q.sendError(md.Id, ev.EventName(), err)
    return
  }

  // Deserialize event
  err := deserializeMessage[E](q.deserializer, message, ev)
  if err != nil {
    q.sendError(md.Id, ev.EventName(), err)
    return
  }

  // Process it
  err = handle(ctx, ev)
  if err != nil {
    q.sendError(md.Id, ev.EventName(), err)
    return
  }
}
