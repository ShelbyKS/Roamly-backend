package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type MessageBrokerProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewMessageBrokerProducer(host string, port string, topic string) *MessageBrokerProducer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", host, port),
	})
	if err != nil {
		return nil
	}

	return &MessageBrokerProducer{
		producer: producer,
		topic:    topic,
	}
}

func (m *MessageBrokerProducer) SendMessage(msg model.NotifyMessage) error {
	jsonMessage, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("Failed to serialize message: %w", err)
	}

	deliveryChan := make(chan kafka.Event, 1)
	err = m.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &m.topic, Partition: kafka.PartitionAny},
		Value:          jsonMessage,
	}, deliveryChan)
	if err != nil {
		return fmt.Errorf("Failed to produce message: %w", err)
	}

	return nil
}
