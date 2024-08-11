package kafka

import (
  "errors"
  "github.com/IBM/sarama"
  "io"
)

type ClientWith[T io.Closer] struct {
  client sarama.Client
  closer T
}

func NewSyncProducer(client sarama.Client) SyncProducer {
  return SyncProducer{ClientWith: ClientWith[sarama.SyncProducer]{
    client: client,
    closer: nil,
  }}
}

type SyncProducer struct {
  ClientWith[sarama.SyncProducer]
}

func (s *SyncProducer) Close() error {
  var err error
  err = s.client.Close()
  err = errors.Join(err, s.closer.Close())
  return err
}

func (s *SyncProducer) Client() sarama.SyncProducer {
  return s.closer
}

type Consumer struct {
  ClientWith[sarama.Consumer]
}

func (s *Consumer) Close() error {
  var err error
  err = s.client.Close()
  err = errors.Join(err, s.closer.Close())
  return err
}

func (s *Consumer) Client() sarama.Consumer {
  return s.closer
}

type GroupConsumer struct {
  ClientWith[sarama.ConsumerGroup]
}

func (s *GroupConsumer) Close() error {
  var err error
  err = s.client.Close()
  err = errors.Join(err, s.closer.Close())
  return err
}

func (s *GroupConsumer) Client() sarama.ConsumerGroup {
  return s.closer
}
