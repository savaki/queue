package queue

import (
	"errors"
	gosqs "github.com/savaki/sqs"
	"log"
	"time"
)

func ReadFromQueue(queueName string) {
	queue := &SQSReader{QueueName: queueName}
	queue.ReadFromQueue()
}

// ReadFromQueue reads data from the queue specified and pumps the results into the provided channel
// this func loops forever and should only be executed in a goroutine
func (r *SQSReader) ReadFromQueue() {
	if r.Messages == nil {
		r.Messages = make(chan Message)
	}
	if r.Errs == nil {
		r.Errs = make(chan error)
	}
	if r.Del == nil {
		r.Del = make(chan string)
	}
	if r.QueueName == "" {
		r.Errs <- errors.New("NewReader - misssing QueueName")
		return
	}

	q, err := LookupQueue(r.QueueName)
	if err != nil {
		r.Errs <- err
		return
	}

	// create a new channel to handle the deletes
	go deleteFromQueue(q, r.Del)

	r.readFromQueueForever(q)
}

func (r *SQSReader) readFromQueueForever(q *gosqs.Queue) {
	// loop forever receiving messages
	for {
		err := r.readFromQueueOnce(q)
		if err != nil {
			r.Errs <- err
			delay := 15 * time.Second
			log.Printf("ReadFromQueue: error received while attempting to read from queue, %s -- %s\n", r.QueueName, err.Error())
			log.Printf("ReadFromQueue: waiting %s until before retrying", delay.String())
			<-time.After(delay)
		}
	}
}

func (r *SQSReader) readFromQueueOnce(q *gosqs.Queue) error {
	log.Printf("%s: reading messages from queue\n", r.QueueName)
	results, err := q.ReceiveMessage(RECV_ALL, RECV_MAX_MESSAGES, RECV_VISIBILITY_TIMEOUT)
	if err != nil {
		return err
	}

	for _, message := range results.Messages {
		r.enqueueMessage(message.Body, message.ReceiptHandle)
	}

	return nil
}

func (r *SQSReader) enqueueMessage(text string, handle string) {
	onComplete := func() {
		log.Printf("sending to channel del <- %s\n", handle)
		r.Del <- handle
	}

	var message Message = &msg{data: []byte(text), callback: onComplete}
	r.Messages <- message // push this request into the reader
}
