package sqs

import (
	"encoding/json"
)

func ExampleReadingFromQueue() {
	queueName := "your-queue-here"
	regionName := "us-east-1"

	builder := New(queueName, regionName)
	go builder.ReadFromQueue()
	go builder.ReadFromQueue()

	reader, _ := builder.BuildReader()

	message := <-reader.Receiver()
	properties := make(map[string]string)
	json.NewDecoder(message).Decode(&properties)
	message.Delete()
}

func ExampleWriteToQueue() {
	queueName := "your-queue-here"
	regionName := "us-east-1"

	builder := New(queueName, regionName)
	go builder.WriteToQueue()

	writer, _ := builder.BuildWriter()
	data, _ := json.Marshal(map[string]string{"hello": "world"})
	writer.WriteMessage(data)
}
