package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	// if len(os.Args) < 2 {
	// 	fmt.Fprintf(os.Stderr, "Usage: %s <config-file-path>\n",
	// 		os.Args[0])
	// 	os.Exit(1)
	// }

	configFile := "getting-started.properties"
	conf := ReadConfig(configFile)
	conf["group.id"] = "conv" //"kafka-go-getting-started"
	conf["auto.offset.reset"] = "earliest"

	c, err := kafka.NewConsumer(&conf)

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	topic := "purchases"
	convertedTopic := "purchases-converted"
	err = c.SubscribeTopics([]string{topic}, nil)
	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	p, err := kafka.NewProducer(&conf)
	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}
			fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &convertedTopic, Partition: kafka.PartitionAny},
				Key:            []byte(strings.ToUpper(string(ev.Key))),
				Value:          []byte(strings.ToUpper(string(ev.Value))),
			}, nil)
		}
	}

	c.Close()

}
