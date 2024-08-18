package config

type Database struct {
  RedisAddress  string `env:"REDIS_ADDRESS,notEmpty"`
  RedisUsername string `env:"REDIS_USERNAME,notEmpty"`
  RedisPassword string `env:"REDIS_PASSWORD,notEmpty"`
}

type MessageBroker struct {
  Addresses    []string `env:"BROKER_ADDRESSES,notEmpty"`
  KafkaVersion string   `env:"BROKER_KAFKA_VERSION"`
  GroupId      string   `env:"BROKER_GROUP_ID,notEmpty"`
}
