package main

import (
    "fmt"
    "strings"

    "github.com/Shopify/sarama"

    "github.com/openrelayxyz/cardinal-streams/transports"
)

func strPtr(x string) *string {
	return &x
}

func createProducer(broker, topic string) (sarama.AsyncProducer, error) {

    brokers, config := transports.ParseKafkaURL(strings.TrimPrefix(broker, "kafka://"))
    configEntries := make(map[string]*string)
    configEntries["retention.ms"] = strPtr("3600000")
    if err := transports.CreateTopicIfDoesNotExist(strings.TrimPrefix(broker, "kafka://"), topic, 0, configEntries); err != nil {
        panic(fmt.Sprintf("Could not create topic %v on broker %v: %v", topic, broker, err.Error()))
    }
    
    producer, err := sarama.NewAsyncProducer(brokers, config)
    if err != nil {
        panic(fmt.Sprintf("Could not setup producer, peer manager plugin: %v", err.Error()))
    }

    return producer, nil
}

func consume(topic, broker string) {
    brokers := []string{broker}

    consumer, err := sarama.NewConsumer(brokers, nil)
    if err != nil {
        log.Error("Failed to start Sarama consumer", "err", err)
    }

    partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
    if err != nil {
        log.Error("Failed to start Sarama partition consumer", "err", err)
    }

    go func() {
        for message := range partitionConsumer.Messages() {
            nodes <- string(message.Value)
        }
    }()

     <-exit 
}