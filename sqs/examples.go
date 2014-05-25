package sqs

import (
	"encoding/json"
)

func ExampleReadingFromQueue() {
	queueName := "your-queue-here"
	regionName := "us-east-1"

	client := New(queueName, regionName)
	go client.ReadFromQueue()
	go client.ReadFromQueue()

	message := <-client.Inbound
	properties := make(map[string]string)
	json.NewDecoder(message).Decode(&properties)
	message.Delete()
}

func ExampleWriteToQueue() {
	queueName := "your-queue-here"
	regionName := "us-east-1"

	client := New(queueName, regionName)
	go client.WriteToQueue()

	data, _ := json.Marshal(map[string]string{"hello": "world"})
	client.Outbound <- data
}
