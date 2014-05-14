queue
=====
[![Build Status](https://travis-ci.org/savaki/queue.svg?branch=master)](https://travis-ci.org/savaki/queue)

Simplified interface to SQS.


### Reading Messages From SQS 

In the following example, we continuously read messages from the queue named, "your-queue-name", and push them into the channel, messages.

```
import "github.com/savaki/queue"

func ExampleReadingFromQueue() {
	queueName := "your-queue-here"
	regionName := "us-west-1"
	messages := make(chan queue.Message)

	go queue.ReadFromQueue(queueName, regionName, messages)

	properties := make(map[string]string)
	message := <-messages
	message.Unmarshal(&properties) // unwrap the json data
	message.OnComplete()           // call after you've successfully processed the message
}
```

## Writing Messages To SQS

In this example, we'll write to SQS via the channel, messages.

```
import "github.com/savaki/queue"

func ExampleWriteToQueue() {
	queueName := "your-queue-here"
	regionName := "us-west-1"
	messages := make(chan interface{})

	go queue.WriteToQueue(queueName, regionName, messages)
	messages <- map[string]string{"hello": "world"} // write your message to the queue
}
```


