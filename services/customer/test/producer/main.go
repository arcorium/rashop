package main

import (
  "context"
  "github.com/IBM/sarama"
  "github.com/arcorium/rashop/shared/serde"
  "github.com/arcorium/rashop/shared/types"
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

  repo := publisher.NewKafka(publish, serde.JsonSerializer{})
  v1 := &event.CustomerCreatedV1{
    DomainEventBaseV1: types.NewDomainEventV1(),
    CustomerId:        types.MustCreateId().String(),
    Email:             "something@gmail.com",
    Password:          "Hello123",
    Username:          "username",
    FirstName:         "first name",
    LastName:          "",
    Balance:           0,
    Point:             0,
    IsVerified:        false,
    IsDisabled:        false,
    CreatedAt:         time.Now(),
  }
  ctx := context.Background()
  err = repo.PublishEvents(ctx, v1)
  if err != nil {
    log.Fatalln(err)
  }

  //serialize, err := serde.JsonSerializer{}.Serialize(v1)
  //if err != nil {
  //  log.Fatalln(err)
  //}
  //
  //result := event.CustomerCreatedV1{}
  //err = serde.JsonAnyDeserializer{}.Deserialize(serialize, &result)
  //if err != nil {
  //  log.Fatalln(err)
  //}
}
