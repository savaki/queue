queue
=====

Simplified interface to SQS.


### Reading Messages From SQS 

In the following example, we continuously read messages from the queue named, "your-queue-name"

```
import "github.com/savaki/queue"

func sample() {
  queueName   := "your-queue-name"
  go queue.ReadFromQueue(queueName)

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
  go queue.WriteToQueue(queueName)
  
  for {
  	ch <- map[string]string{"hello":"world"}
  }
}
```


