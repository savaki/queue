package sqs

import (
	gosqs "github.com/crowdmob/goamz/sqs"
	"github.com/savaki/queue"
	"time"
)

func readFromQueueForever(q *gosqs.Queue, inbound chan queue.Message, delete chan string) {
	// loop forever receiving messages
	for {
		err := readFromQueueOnce(q, inbound, delete)
		if err != nil {
			delay := 15 * time.Second
			<-time.After(delay)
		}
	}
}

func readFromQueueOnce(q *gosqs.Queue, inbound chan queue.Message, delete chan string) error {
	results, err := q.ReceiveMessage(RECV_MAX_MESSAGES)
	if err != nil {
		return err
	}

	for _, message := range results.Messages {
		enqueueMessage(inbound, delete, message.Body, message.ReceiptHandle)
	}

	return nil
}

func enqueueMessage(inbound chan queue.Message, delete chan string, text string, handle string) {
	var message queue.Message = &msg{
		data: []byte(text),
		callback: func() {
			delete <- handle
		},
	}
	inbound <- message
}
