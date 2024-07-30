package messaging

import (
	"fmt"
	"messageservice"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func MessageToKafka(messsages messageservice.MessageItem) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",

		"acks": "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}
	timestamp := time.Now()

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	topic := "messages"

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(messsages.Text),
		Timestamp:      timestamp,
	}, nil)

	p.Flush(15 * 1000)
	p.Close()
}
