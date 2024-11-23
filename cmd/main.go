package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zmskv/websocket_stream/internal/config"
	"github.com/zmskv/websocket_stream/internal/handler"
	"github.com/zmskv/websocket_stream/internal/kafka"
)

func main() {
	brokers, topic, groupID := config.GetKafkaConfig()

	kafkaConsumer, err := kafka.NewKafkaConsumer(brokers, groupID, topic)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer kafkaConsumer.Close()

	broadcast := make(chan string)

	go func() {
		for msg := range kafkaConsumer.Output {
			broadcast <- msg
		}
	}()

	wsHandler := handler.NewWebSocketHandler(broadcast)

	r := gin.Default()

	r.GET("/ws", wsHandler.HandleWebSocket)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
