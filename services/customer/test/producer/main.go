package main

import (
  "context"
  "github.com/IBM/sarama"
  "github.com/arcorium/rashop/shared/serde"
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  "github.com/dnwe/otelsarama"
  "log"
  "mini-shop/services/user/constant"
  "mini-shop/services/user/internal/domain/event"
  "mini-shop/services/user/internal/infra/publisher"
  "time"
)

func createKafkaConfig() (*sarama.Config, error) {
  conf := sarama.NewConfig()
  conf.Version = constant.DefaultKafkaVersion
  conf.Producer.RequiredAcks = sarama.WaitForAll
  conf.Producer.Return.Successes = true
  return conf, nil
}

func publisherSetup() (sarama.SyncProducer, error) {
  producerCfg, err := createKafkaConfig()
  if err != nil {
    return nil, err
  }

  producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, producerCfg)
  if err != nil {
    return nil, err
  }

  return otelsarama.WrapSyncProducer(producerCfg, producer), nil
}

func main() {
  publish, err := publisherSetup()
  if err != nil {
    log.Fatalln(err)
  }

  custId := types.MustCreateId().String()
  repo := publisher.NewKafka(publisher.KafkaTopic{
    DomainEvent:      constant.CUSTOMER_DOMAIN_EVENT_TOPIC,
    IntegrationEvent: constant.CUSTOMER_INTEGRATION_EVENT_TOPIC,
  }, publish, serde.JsonSerializer{})
  v1 := &event.CustomerCreatedV1{
    DomainV1:   event.NewV1(),
    CustomerId: custId,
    Email:      "something123@ymail.com",
    Password:   "asdasdasd",
    Username:   "alskdasdl",
    FirstName:  "first name213",
    LastName:   "",
    Balance:    0,
    Point:      0,
    IsVerified: false,
    IsDisabled: false,
    CreatedAt:  time.Now(),
  }

  ev2 := &event.CustomerBalanceUpdatedV1{
    DomainV1:   event.NewV1(),
    CustomerId: custId,
    Balance:    100_000,
    Point:      100,
  }

  sharedUtil.DoNothing(v1, ev2)

  ctx := context.Background()
  err = repo.PublishEvents(ctx, v1, ev2)
  if err != nil {
    log.Fatalln(err)
  }
}
