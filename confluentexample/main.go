package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var (
	topic     = "topic_0"
	producerC *kafka.Producer
	consumerC *kafka.Consumer
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

	// Create Producer instance
	producerC = instanceProducer()
	consumerC = instanceConsumer()

	// Goroutine to handle message delivery reports and
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

	// Produce messages in goroutine
	go producer()

	// Get messages
	consumer()
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

func instanceConsumer() *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": conf["bootstrap.servers"],
		"sasl.mechanisms":   conf["sasl.mechanisms"],
		"security.protocol": conf["security.protocol"],
		"sasl.username":     conf["sasl.username"],
		"sasl.password":     conf["sasl.password"],
		"group.id":          "go_example_group_1",
		"auto.offset.reset": "earliest"})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	return c
}

func producer() {
	log.Println("START PRODUCING: " + time.Now().String())
	for n := 0; n < 10; n++ {
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

	// Wait for all messages to be delivered
	producerC.Flush(2 * 1000)

	fmt.Printf("10 messages were produced to topic %s!", topic)
	log.Println("FINISH PERSIST: " + time.Now().String())

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
			log.Println("START READING: " + time.Now().String())
			msg, err := consumerC.ReadMessage(100 * time.Millisecond)
			if err != nil {
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
			log.Println("FISNIH READING: " + time.Now().String())
		}
	}

	fmt.Printf("Closing consumer\n")
	consumerC.Close()
}

func readCCloudConfig(configFile string) map[string]string {
	m := make(map[string]string)

	file, err := os.Open(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %s", err)
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
		fmt.Printf("Failed to read file: %s", err)
		os.Exit(1)
	}

	return m
}
