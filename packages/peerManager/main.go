package main

import (
	"fmt"
	"flag"
	"github.com/Shopify/sarama"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"
)

var (
	log core.Logger
	httpApiFlagName = "http.api"
	sessionStack core.Node
	sessionBrokers []string
	sessionKafkaConfig *sarama.Config
	nodes = make(chan string, 5)
	exit = make(chan struct{}, 1)

	Flags = *flag.NewFlagSet("peermanager-plugin", flag.ContinueOnError)
	peerBroker = Flags.String("peer-broker", "", "kafka broker for peer manager")
)

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	log.Info("Loaded peer manager plugin")
}

func InitializeNode(stack core.Node, b restricted.Backend) {
	sessionStack = stack
	log.Info("Initialized node, peer manager plugin")
}

func BlockChain() {
	if *peerBroker == "" {
		panic(fmt.Sprintf("no broker provided for peer manager plugin"))
	}
	go peeringSequence()
}

func peeringSequence() {

	sessionPeerService, err := getPeerManager()
	if err != nil {
		log.Error("session peer service unavailable, peer manager plugin", "err", err)
		return
	}
	
	selfNode, err := sessionPeerService.getEnode()
	if err != nil {
		log.Error("error calling getEnode from sessionService, peer manager plugin", "err", err)
	} 

	chainTopic, err := sessionPeerService.chainIdResolver() 
	if err != nil {
		log.Error("Error aquiring chainID, peer manager plugin", "err", err)
	} 

	producer, err := createProducer(*peerBroker, chainTopic)
	if err != nil {
		log.Error("failed to acquire kafka producer, peer manager plugin", "err", err)
		return
	}

	consumer, err := createConsumer(*peerBroker, chainTopic)
	if err != nil {
		log.Error("failed to acquire kafka consumer, peer manager plugin", "err", err)
		return
	}

	msg := &sarama.ProducerMessage{
	        Topic: chainTopic,
	        Value: sarama.StringEncoder(selfNode),
	}

	producer.Input() <- msg

	go func() {
		for message := range consumer.Messages() {
            nodes <- string(message.Value)
        }
	}()

	for message := range nodes {
		if message == selfNode {
			continue
		} else {
			sessionPeerService.attachPeers(message)
		}
	}
}