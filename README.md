queue
=====

Simplified interface to SQS.


### Reading Messages From SQS 

In the following example, we continuously read messages from the queue named, "your-queue-name", and push them into the channel, messages.

```
import "github.com/savaki/queue"

func ExampleReadingFromQueue() {
	queueName := "your-queue-here"
	messages := make(chan queue.Message)

	go queue.ReadFromQueue(queueName, messages)

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
	messages := make(chan interface{})

	go queue.WriteToQueue(queueName, messages)
	messages <- map[string]string{"hello": "world"} // write your message to the queue
}
```


