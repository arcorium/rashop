package config

import sharedConf "github.com/arcorium/rashop/shared/config"

type Database struct {
  sharedConf.PostgresDatabase
}

type MessageBroker struct {
  Addresses    []string `env:"BROKER_ADDRESSES,notEmpty"`
  KafkaVersion string   `env:"BROKER_KAFKA_VERSION"`
  GroupId      string   `env:"BROKER_GROUP_ID,notEmpty"`
}

type MinIO struct {
  Address   string `env:"MINIO_ADDRESS,notEmpty"`
  AccessKey string `env:"MINIO_ACCESS_KEY,notEmpty"`
  SecretKey string `env:"MINIO_SECRET_KEY,notEmpty"`
}
