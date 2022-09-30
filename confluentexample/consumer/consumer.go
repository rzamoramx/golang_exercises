package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var (
	topic     = "topic_1"
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
	fmt.Printf("map conf: %+v\n", conf)

	// Create instance
	consumerC = instanceConsumer()

	// Get messages
	consumer()
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

func consumer() {
	// Subscribe to topic
	err := consumerC.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("failed to subscribe: %s\n", err)
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
				//fmt.Printf("Failed to read message: %v\n", err)
				continue
			}
			fmt.Printf("START READING: %s\n", time.Now().String())

			recordKey := string(msg.Key)
			recordValue := msg.Value
			data := RecordValue{}
			err = json.Unmarshal(recordValue, &data)
			if err != nil {
				fmt.Printf("Failed to decode JSON at offset %d: %v\n", msg.TopicPartition.Offset, err)
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
