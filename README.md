queue
=====

Simplified interface to SQS.


### Reading Messages From SQS

In the following example, we continuously read messages from the queue named, "your-queue-name"

```
import "github.com/savaki/queue"

func sample() {
  queueName   := "your-queue-name"
  ch          := make(chan queue.Message)
  errs        := make(chan error)

  // read messages from sqs in the background
  go queue.ReadFromQueue(queueName, ch, errs)

  for {
    // process those messages here
    message := <-ch 
    // … do some processing …
    message.OnComplete() 
  }
}
```

## Writing Messages To SQS

In this example, we'll do the following, we'll write messags to SQS.

```
import "github.com/savaki/queue"

func sample() {
  queueName   := "your-queue-name"
  ch          := make(chan interface{})
  errs        := make(chan error)

  // write messages to sqs in the background
  go queue.WriteToQueue(queueName, ch, errs)
  
  for {
  	ch <- map[string]string{"hello":"world"}
  }
}
```


