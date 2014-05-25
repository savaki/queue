queue
=====
[![Build Status](https://travis-ci.org/savaki/queue.svg?branch=master)](https://travis-ci.org/savaki/queue)

Simplified interface to SQS.

### Setting the environment

Currently, queue assumes your AWS credentials have been set via the environment

```
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-access-key"
```

### Reading Messages From SQS 

In the following example, we continuously read messages from the queue named, "your-queue-name", and push them into the channel, messages.

```
import (
  "github.com/savaki/queue/sqs"
  "json"
)

func ExampleReadingFromQueue() {
	queueName := "your-queue-here"
	regionName := "us-west-1"
	
	client := sqs.ReadFromQueue(queueName, regionName)

	go client.ReadFromQueue()

	properties := make(map[string]string)
	message := <-client.Inbound
	json.NewDecoder(message).Decode(&properties)
	message.Delete() // call when you've successfully processed the message
}
```

## Writing Messages To SQS

In this example, we'll write to SQS via the channel, messages.

```
import (
  "github.com/savaki/queue/sqs"
  "json"
)

func ExampleWriteToQueue() {
	queueName := "your-queue-here"
	regionName := "us-west-1"
	
	client := sqs.New(queueName, regionName)

	go queue.WriteToQueue()

	value := map[string]string{"hello": "world"}
	data, _ := json.Marshal(value)
	messages <- data  // write your message to the queue
}
```


