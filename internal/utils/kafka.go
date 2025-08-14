package utils

import (
	"encoding/json"

	kafka "github.com/IBM/sarama"
	"github.com/go-logr/logr"
)

type Kafka struct {
	producer kafka.SyncProducer
}

func NewKafka(producer kafka.SyncProducer) *Kafka {
	return &Kafka{producer: producer}
}

func (k *Kafka) SendKafkaMessage(topic string, message any, logger *logr.Logger) error {
	data, _ := json.Marshal(message)

	msg := &kafka.ProducerMessage{
		Topic: topic,
		Value: kafka.StringEncoder(data),
	}

	if _, _, err := k.producer.SendMessage(msg); err != nil {
		return err
	}

	logger.Info("Message sent to Kafka", "message", message)

	return nil
}
