package config

import (
	"os"
	"strings"
)

func GetKafkaConfig() (brokers []string, topic string, groupID string) {
	brokersEnv := os.Getenv("KAFKA_BROKERS")
	if brokersEnv == "" {
		brokersEnv = "localhost:9092"
	}
	brokers = strings.Split(brokersEnv, ",")
	topic = os.Getenv("TOPIC_NAME")
	groupID = os.Getenv("GROUP_ID")
	return
}
