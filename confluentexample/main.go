package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var (
	topic     = "topic_0"
	producerC *kafka.Producer
)

// https://github.com/confluentinc/examples/blob/7.2.1-post/clients/cloud/go/producer.go

// RecordValue represents the struct of the value in a Kafka message
type RecordVal RecordValue

type RecordValue struct {
	Count int
}

func main() {
	// Create Producer instance
	producerC = instanceProducer()

	// Go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range producerC.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Successfully produced record to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()

	// run producer
	go producer()
}

func instanceProducer() *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "pkc-3w22w.us-central1.gcp.confluent.cloud:909",
		"sasl.mechanisms":   "PLAIN",
		"security.protocol": "SASL_SSL",
		"sasl.username":     "XCBOZPRUIH4HQGXP",
		"sasl.password":     "ShDVNow4/6PEREV1MkGIjQhoXXk3kFvKfbx2htZv2btqKfETRYOOYWwWMqC9hO2R"})
	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	return p
}

func producer() {
	for n := 0; n < 10; n++ {
		recordKey := "alice"
		data := &RecordValue{
			Count: n}
		recordValue, _ := json.Marshal(&data)
		fmt.Printf("Preparing to produce record: %s\t%s\n", recordKey, recordValue)
		producerC.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: 1}, // kafka.PartitionAny},
			Key:            []byte(recordKey),
			Value:          []byte(recordValue),
		}, nil)
	}

	// Wait for all messages to be delivered
	producerC.Flush(15 * 1000)

	fmt.Printf("10 messages were produced to topic %s!", topic)

	producerC.Close()
}
