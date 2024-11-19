package main

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka1:9092",
		"group.id":          "test-consumer-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Ошибка при создании потребителя: %s\n", err)
	}
	defer consumer.Close()

	err = consumer.Subscribe("topic1", nil)
	if err != nil {
		log.Fatalf("Ошибка при подписке на топик: %s\n", err)
	}

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			log.Printf("Получено сообщение: %s\n", string(msg.Value))
		} else {
			log.Printf("Ошибка чтения сообщения: %v\n", err)
		}
	}
}
