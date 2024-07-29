package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	KafkaMessages = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_messages_consumed_total",
			Help: "Total number of messages consumed from Kafka",
		},
		[]string{"topic"},
	)
	KafkaErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_errors_total",
			Help: "Total number of errors encountered by the Kafka consumer",
		},
		[]string{"topic"},
	)
)

func init() {
	prometheus.MustRegister(KafkaMessages)
	prometheus.MustRegister(KafkaErrors)
}
