package messaging

import (
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
	/*go func() {
		for {
			ev := kc.c.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:
				av, err := kc.c.ReadMessage(100 * time.Millisecond)
				if err != nil {
					continue
				}
				fmt.Printf("Consumed event from topic %s: text = %s\n, time = %s\n",
					*av.TopicPartition.Topic, string(av.Value), av.Timestamp)
			case kafka.Error:
				fmt.Printf("Error: %v\n", e)
			case *kafka.Stats:
				stats := map[string]interface{}{}
				if err := json.Unmarshal([]byte(e.String()), &stats); err != nil {
					log.Printf("Failed to parse statistics: %v\n", err)
				} else {
					log.Printf("Consumer statistics: %v\n", stats)
				}
			}
		}
	}()
	go func() {
		for {
			ev, err := kc.c.ReadMessage(80 * time.Millisecond)
			if err != nil {
				continue
			}
			prometheus.KafkaMessages.WithLabelValues(*e.TopicPartition.Topic).Inc()
			fmt.Printf("Consumed event from topic %s: text = %s\n, time = %s\n",
				*ev.TopicPartition.Topic, string(ev.Value), ev.Timestamp)
		}
	}() */
	go func() {
		for {
			ev := kc.c.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:
				prometheus.KafkaMessages.WithLabelValues(*e.TopicPartition.Topic).Inc()
				fmt.Printf("Consumed event from topic %s: text = %s\n, time = %s\n",
					*e.TopicPartition.Topic, string(e.Value), e.Timestamp)
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
