package messaging

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ConsumeKafka() {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:52014",

		"group.id":               "kafka-go-getting-started",
		"auto.offset.reset":      "earliest",
		"statistics.interval.ms": 10000})

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	topic := "messages"
	err = c.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s", err)
	}
	//defer c.Close()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	/*go func() {
		for {
			ev := c.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:
				av, err := c.ReadMessage(100 * time.Millisecond)
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
					if topics, ok := stats["topics"].(map[string]interface{}); ok {
						if topicStats, ok := topics["my-topic"].(map[string]interface{}); ok {
							if partitions, ok := topicStats["partitions"].(map[string]interface{}); ok {
								for _, p := range partitions {
									if partitionStats, ok := p.(map[string]interface{}); ok {
										if msgCount, ok := partitionStats["msg_cnt"].(float64); ok {
											log.Printf("Messages processed: %d\n", int(msgCount))
										}
										if errCount, ok := partitionStats["rxerrs"].(float64); ok {
											log.Printf("Receive errors: %d\n", int(errCount))
										}
									}
								}
							}
						}
					}
				}
			default:
				// Handle other events
			}
		}
	}() */

	run := true

	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				continue
			}
			fmt.Printf("Consumed event from topic %s: text = %s\n, time = %s\n",
				*ev.TopicPartition.Topic, string(ev.Value), ev.Timestamp)
		}
	}

	c.Close()

	//select {}
}
