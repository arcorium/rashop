package kafka

import (
  "github.com/IBM/sarama"
  "github.com/dnwe/otelsarama"
  "time"
)

type configOption interface {
  apply(*sarama.Config)
}

type funcConfigOption struct {
  f func(*sarama.Config)
}

func (f *funcConfigOption) apply(conf *sarama.Config) {
  f.f(conf)
}

func config(f func(*sarama.Config)) configOption {
  return &funcConfigOption{f: f}
}

func WithDefaultProducer() configOption {
  return &funcConfigOption{f: func(conf *sarama.Config) {
    conf.Producer.RequiredAcks = sarama.WaitForAll
    conf.Producer.Return.Successes = true
    conf.Producer.Return.Errors = true
  }}
}

func WithDefaultConsumer() configOption {
  return &funcConfigOption{f: func(conf *sarama.Config) {
    conf.Consumer.Return.Errors = true
    conf.Consumer.Offsets.Initial = sarama.OffsetOldest
    conf.Consumer.Offsets.AutoCommit.Enable = true
    conf.Consumer.Offsets.AutoCommit.Interval = time.Second * 5
  }}
}

func WithDefaultConsumerGroup(instanceId string) configOption {
  return &funcConfigOption{f: func(conf *sarama.Config) {
    conf.Consumer.Return.Errors = true
    conf.Consumer.Offsets.Initial = sarama.OffsetOldest
    conf.Consumer.Offsets.AutoCommit.Enable = true
    conf.Consumer.Offsets.AutoCommit.Interval = time.Second * 5
    conf.Consumer.Group.InstanceId = instanceId
  }}
}

// DefaultConfig create default sarama.Config with version and provided options
func DefaultConfig(version string, fallback sarama.KafkaVersion, opts ...configOption) *sarama.Config {
  conf := sarama.NewConfig()
  if version != "" {
    var err error
    conf.Version, err = sarama.ParseKafkaVersion(version)
    if err != nil {
      conf.Version = fallback
    }
  } else {
    conf.Version = fallback
  }

  for _, opt := range opts {
    opt.apply(conf)
  }

  return conf
}

type clientOption[T any] interface {
  apply(*sarama.Config, T) T
}

type funcClientOption[T any] struct {
  f func(*sarama.Config, T) T
}

func (f *funcClientOption[T]) apply(conf *sarama.Config, t T) T {
  return f.f(conf, t)
}

type syncProducerOption = clientOption[sarama.SyncProducer]
type funcSyncProducerOption = funcClientOption[sarama.SyncProducer]

func WithOTELSyncProducer() syncProducerOption {
  return &funcSyncProducerOption{f: func(conf *sarama.Config, producer sarama.SyncProducer) sarama.SyncProducer {
    return otelsarama.WrapSyncProducer(conf, producer)
  }}
}

func DefaultSyncProducer(brokers []string, conf *sarama.Config, opts ...syncProducerOption) (sarama.SyncProducer, error) {
  producer, err := sarama.NewSyncProducer(brokers, conf)
  if err != nil {
    return nil, err
  }

  for _, opt := range opts {
    producer = opt.apply(conf, producer)
  }

  return producer, nil
}

func DefaultSyncGroupConsumer(brokers []string, groupId string, conf *sarama.Config) (sarama.ConsumerGroup, error) {
  group, err := sarama.NewConsumerGroup(brokers, groupId, conf)
  if err != nil {
    return nil, err
  }

  return group, nil
}
