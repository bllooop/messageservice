package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"messageservice/prometheus"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	c *kafka.Consumer
}

func (kc *KafkaConsumer) ConsumeKafka() error {

	var err error
	kc.c, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",

		"group.id":               "kafka-go-getting-started",
		"auto.offset.reset":      "earliest",
		"statistics.interval.ms": 10000})

	if err != nil {
		return fmt.Errorf("failed to create consumer: %s", err)
	}

	topic := "messages"
	err = kc.c.Subscribe(topic, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %s", err)
	}
	//defer c.Close()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			ev := kc.c.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:
				prometheus.KafkaMessages.WithLabelValues(*e.TopicPartition.Topic).Inc()
				messageData := map[string]interface{}{
					"topic":     *e.TopicPartition.Topic,
					"text":      string(e.Value),
					"timestamp": e.Timestamp,
				}
				jsonData, err := json.MarshalIndent(messageData, "", "  ")
				if err != nil {
					log.Printf("Error marshaling JSON: %s", err)
					continue
				}
				fmt.Println("Consumed event " + string(jsonData))

			/*case kafka.Error:
			prometheus.KafkaErrors.WithLabelValues(*e.Code().String()).Inc()
			fmt.Printf("Error: %v\n", e) */
			default:
			}
		}
	}()
	return nil
}

func (kc *KafkaConsumer) Stop() {
	if kc.c != nil {
		kc.c.Close()
		kc.c = nil
		log.Println("Consumer stopped")
	}
}
