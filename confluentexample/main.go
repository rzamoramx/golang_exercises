package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var (
	topic      = "topic_0"
	producerC  *kafka.Producer
	consumerC  *kafka.Consumer
	mechanisms = "PLAIN"
	protocol   = "SASL_SSL"
	username   = "username"
	password   = "password"
	server     = "server:9092"
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
	consumerC = instanceConsumer()

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

	go producer()

	consumer()
}

func instanceProducer() *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": server,
		"sasl.mechanisms":   mechanisms,
		"security.protocol": protocol,
		"sasl.username":     username,
		"sasl.password":     password})
	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	return p
}

func instanceConsumer() *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": server,
		"sasl.mechanisms":   mechanisms,
		"security.protocol": protocol,
		"sasl.username":     username,
		"sasl.password":     password,
		"group.id":          "go_example_group_1", // es arbitrario?
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	return c
}

func producer() {
	for n := 0; n < 10; n++ {
		recordKey := "alice"
		data := &RecordValue{
			Count: n}
		recordValue, _ := json.Marshal(&data)
		fmt.Printf("Preparing to produce record: %s\t%s\n", recordKey, recordValue)

		err := producerC.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte(recordKey),
			Value:          []byte(recordValue),
		}, nil)

		if err != nil {
			fmt.Printf("Error to produce message: %v", err)
		}
	}

	// Wait for all messages to be delivered
	producerC.Flush(15 * 1000)

	fmt.Printf("10 messages were produced to topic %s!", topic)

	producerC.Close()
}

func consumer() {
	// Subscribe to topic
	err := consumerC.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("failed to subscribe: %s", err)
	}
	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	totalCount := 0
	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			msg, err := consumerC.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				//fmt.Printf("Failed to read message: %v\n", err)
				continue
			}
			recordKey := string(msg.Key)
			recordValue := msg.Value
			data := RecordValue{}
			err = json.Unmarshal(recordValue, &data)
			if err != nil {
				fmt.Printf("Failed to decode JSON at offset %d: %v", msg.TopicPartition.Offset, err)
				continue
			}
			count := data.Count
			totalCount += count
			fmt.Printf("Consumed record with key %s and value %s, and updated total count to %d\n", recordKey, recordValue, totalCount)
		}
	}

	fmt.Printf("Closing consumer\n")
	consumerC.Close()
}
