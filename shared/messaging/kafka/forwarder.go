package kafka

import (
  "context"
  "github.com/IBM/sarama"
  "github.com/arcorium/rashop/shared/messaging"
  sharedUtil "github.com/arcorium/rashop/shared/util"
)

func NewForwarder(topic string, producer sarama.SyncProducer) messaging.IForwarder[*sarama.ConsumerMessage] {
  return &kafkaForwarder{
    producer: producer,
    topic:    topic,
  }
}

type kafkaForwarder struct {
  producer sarama.SyncProducer
  topic    string
}

func (d *kafkaForwarder) Forward(ctx context.Context, message *sarama.ConsumerMessage, err error) error {
  // Append error into headers
  headers := sharedUtil.CastSlice(message.Headers, sharedUtil.ToDeref[sarama.RecordHeader])
  headers = append(headers, sarama.RecordHeader{
    Key:   []byte(messaging.HEADER_ERROR_KEY),
    Value: []byte(err.Error()),
  })

  msg := &sarama.ProducerMessage{
    Topic:     d.topic,
    Key:       sarama.ByteEncoder(message.Key),
    Value:     sarama.ByteEncoder(message.Value),
    Headers:   headers,
    Timestamp: message.Timestamp,
  }

  _, _, err2 := d.producer.SendMessage(msg)
  return err2
}
