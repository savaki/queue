package queue

func ExampleReadingFromQueue() {
	queueName := "your-queue-here"
	messages := make(chan Message)

	go ReadFromQueue(queueName, messages)

	properties := make(map[string]string)
	message := <-messages
	message.Unmarshal(&properties) // unwrap the json data
	message.OnComplete()           // call after you've successfully processed the message
}

func ExampleWriteToQueue() {
	queueName := "your-queue-here"
	messages := make(chan interface{})

	go WriteToQueue(queueName, messages)
	messages <- map[string]string{"hello": "world"} // write your message to the queue
}
