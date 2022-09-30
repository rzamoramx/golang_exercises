package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var (
	topic     = "topic_1"
	producerC *kafka.Producer
	conf      map[string]string
)

// RecordValue represents the struct of the value in a Kafka message
type RecordVal RecordValue

type RecordValue struct {
	Timestamp int64
	Count     int
	Foo       string
}

func main() {
	// read config
	conf = readCCloudConfig("confluent.secret")
	fmt.Printf("map conf: %+v\n", conf)

	// Create Producer instance
	producerC = instanceProducer()

	// Goroutine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	//handler()

	// Produce messages in goroutine
	producer()

}

func instanceProducer() *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": conf["bootstrap.servers"],
		"sasl.mechanisms":   conf["sasl.mechanisms"],
		"security.protocol": conf["security.protocol"],
		"sasl.username":     conf["sasl.username"],
		"sasl.password":     conf["sasl.password"]})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	return p
}

func producer() {
	for n := 0; n < 1; n++ {
		recordKey := "alice"
		data := &RecordValue{
			Count:     n,
			Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
			Foo:       "Bar"}
		recordValue, _ := json.Marshal(&data)
		fmt.Printf("Preparing to produce record: %s\t%s\n", recordKey, recordValue)

		err := producerC.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte(recordKey),
			Value:          []byte(recordValue),
		}, nil)

		if err != nil {
			fmt.Printf("Error to produce message: %v\n", err)
		}
	}

	fmt.Printf("PRODUCER: %s\n", time.Now().String())

	// Wait for all messages to be delivered
	producerC.Flush(15 * 1000)

	fmt.Printf("messages were produced to topic %s!\n", topic)

	producerC.Close()
}

func handler() {
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
}

func readCCloudConfig(configFile string) map[string]string {
	m := make(map[string]string)

	file, err := os.Open(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") && len(line) != 0 {
			kv := strings.Split(line, "=")
			parameter := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			m[parameter] = value
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Failed to read file: %s\n", err)
		os.Exit(1)
	}

	return m
}
