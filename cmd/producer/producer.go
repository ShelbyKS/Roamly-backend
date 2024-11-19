package main

import (
	"fmt"
	"log"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka1:9092",
	})
	if err != nil {
		log.Fatalf("Ошибка при создании продюсера: %s\n", err)
	}
	defer producer.Close()

	// это имя топика, которйы в докер композе указан
	topic := "topic1" 
	message := "Скачивайте ромали.ру на мобильные устройства!"

	deliveryChan := make(chan kafka.Event)
	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)

	e := <-deliveryChan
	msg := e.(*kafka.Message)
	if msg.TopicPartition.Error != nil {
		log.Fatalf("Ошибка отправки сообщения: %v\n", msg.TopicPartition.Error)
	} else {
		fmt.Println("Сообщение успешно отправлено!")
	}
}
