package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	consumer sarama.Consumer
	Output   chan string
}

func NewKafkaConsumer(brokers []string, groupID, topic string) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, err
	}

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}

	kafkaConsumer := &KafkaConsumer{
		consumer: consumer,
		Output:   make(chan string),
	}

	go kafkaConsumer.consumeMessages(topic)

	return kafkaConsumer, nil
}

func (kc *KafkaConsumer) consumeMessages(topic string) {
	partitionList, err := kc.consumer.Partitions(topic)
	if err != nil {
		log.Printf("Error fetching partitions: %v", err)
		return
	}

	for _, partition := range partitionList {
		pc, err := kc.consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
		if err != nil {
			log.Printf("Error starting consumer for partition %d: %v", partition, err)
			continue
		}

		go func(partitionConsumer sarama.PartitionConsumer) {
			defer partitionConsumer.Close()
			for msg := range partitionConsumer.Messages() {
				kc.Output <- string(msg.Value)
			}
		}(pc)
	}
}

func (kc *KafkaConsumer) Close() {
	if err := kc.consumer.Close(); err != nil {
		log.Printf("Error closing Kafka consumer: %v", err)
	}
}
