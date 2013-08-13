package queue

import (
	"encoding/json"
	"errors"
	"github.com/nu7hatch/gouuid"
	gosqs "github.com/savaki/sqs"
	"launchpad.net/goamz/aws"
	"log"
	"strconv"
	"time"
)

type Message interface {
	OnComplete()
}

type SQSReader struct {
	QueueName  string
	NewRequest func(func()) Message
	BatchSize  int
}

type SQSWriter struct {
	QueueName string
	BatchSize int
}

var (
	// configure our sqs read settings
	DEFAULT_TIMEOUT         time.Duration = 5 * time.Minute
	RECV_ALL                              = []string{"All"}
	RECV_MAX_MESSAGES                     = 10
	RECV_VISIBILITY_TIMEOUT               = 900 // 15 minutes
	DELETE_BATCH_SIZE                     = 10  // delete after theis many messages have been sent
)

func assembleDeleteMessageBatch(del chan string, timeout time.Duration) []gosqs.DeleteMessageBatch {
	batch := make([]gosqs.DeleteMessageBatch, 0)

	for index := 0; index < DELETE_BATCH_SIZE; index++ {
		var handle string = ""
		select {
		case handle = <-del:
			log.Printf("assembleDeleteMessageBatch: received handle, %s\n", handle)
			message := gosqs.DeleteMessageBatch{
				Id:            strconv.Itoa(index + 1),
				ReceiptHandle: handle,
			}
			batch = append(batch, message)

		case <-time.After(timeout):
			return batch
		}
	}

	return batch
}

func DeleteFromQueue(q *gosqs.Queue, del chan string) error {
	for {
		err := deleteFromQueueOnce(q, del)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteFromQueueOnce(q *gosqs.Queue, del chan string) error {
	batch := assembleDeleteMessageBatch(del, DEFAULT_TIMEOUT)

	if len(batch) > 0 {
		log.Printf("deleting %d messages from q\n", len(batch))
		_, err := q.DeleteMessageBatch(batch)
		if err != nil {
			return err
		}
	}

	return nil
}

func LookupQueue(queueName string) *gosqs.Queue {
	// login to sqs, reading our credentials from the environment
	auth, err := aws.EnvAuth()
	if err != nil {
		panic(err)
	}

	// connect to our sqs q
	log.Printf("looking up sqs q, %s\n", queueName)
	sqs := gosqs.New(auth, aws.Region{Name: "USWest2", SQSEndpoint: "http://sqs.us-west-2.amazonaws.com"})
	// sqs := gosqs.New(auth, aws.USWest2)
	q, err := sqs.GetQueue(queueName)
	if err != nil {
		panic(err)
	}
	log.Printf("%s: ok\n", queueName)

	return q
}

func NewWriter(queueName string, batchSize int) (*SQSWriter, error) {
	if queueName == "" {
		return nil, errors.New("NewReader - misssing QueueName")
	}
	return &SQSWriter{QueueName: queueName, BatchSize: batchSize}, nil
}

func NewReader(queueName string, newRequest func(func()) Message) (*SQSReader, error) {
	if queueName == "" {
		return nil, errors.New("NewReader - misssing QueueName")
	}
	if newRequest == nil {
		return nil, errors.New("NewReader - missing factory, newRequest")
	}
	return &SQSReader{QueueName: queueName, NewRequest: newRequest}, nil
}

// ReadFromQueue reads data from the queue specified and pumps the results into the provided channel
// this func loops forever and should only be executed in a goroutine
func (r *SQSReader) ReadFromQueue(reader chan Message) {
	q := LookupQueue(r.QueueName)

	// create a new channel to handle the deletes
	del := make(chan string)

	// start up the delete
	go DeleteFromQueue(q, del)

	// loop forever receiving messages
	for {
		err := r.readFromQueueOnce(q, reader, del)
		if err != nil {
			delay := 15 * time.Second
			log.Printf("ReadFromQueue: error received while attempting to read from queue, %s -- %s\n", r.QueueName, err.Error())
			log.Printf("ReadFromQueue: waiting %s until before retrying", delay.String())
			<-time.After(delay)
		}
	}
}

func (r *SQSReader) enqueueRequest(text string, handle string, reader chan Message, del chan string) {
	onComplete := func() {
		log.Printf("sending to channel del <- %s\n", handle)
		del <- handle
	}

	request := r.NewRequest(onComplete)
	json.Unmarshal([]byte(text), &request)

	reader <- request // push this request into the reader
}

func (r *SQSReader) readFromQueueOnce(queue *gosqs.Queue, reader chan Message, del chan string) error {
	log.Printf("%s: reading messages from queue\n", r.QueueName)
	results, err := queue.ReceiveMessage(RECV_ALL, RECV_MAX_MESSAGES, RECV_VISIBILITY_TIMEOUT)
	if err != nil {
		return err
	}

	for _, message := range results.Messages {
		r.enqueueRequest(message.Body, message.ReceiptHandle, reader, del)
	}

	return nil
}

func (w *SQSWriter) assembleSendMessageBatch(writer chan Message) ([]gosqs.SendMessageBatchRequestEntry, func()) {
	requests := make([]gosqs.SendMessageBatchRequestEntry, 0)
	results := make([]Message, 0)

	for index := 0; index < w.BatchSize; index++ {
		var result Message = nil
		select {
		case result = <-writer:
			text, err := json.Marshal(result)
			if err != nil {
				panic(err)
			}
			html := string(text)
			log.Printf("assembleSendMessageBatch: queueing message, %d bytes\n", len(html))
			id, err := uuid.NewV4()
			if err != nil {
				panic(err)
			}
			request := gosqs.SendMessageBatchRequestEntry{
				Id:          id.String(),
				MessageBody: html,
			}
			requests = append(requests, request)
			results = append(results, result)

		case <-time.After(5 * time.Minute):
			index = w.BatchSize
		}
	}

	onComplete := func() {
		log.Printf("onComplete: invoking onComplete for %d elements\n", len(results))
		for _, result := range results {
			log.Printf("onComplete: calling OnComplete()\n")
			go result.OnComplete()
		}
		log.Printf("onComplete: done\n")
	}

	return requests, onComplete
}

func (w *SQSWriter) WriteToQueue(writer chan Message) {
	q := LookupQueue(w.QueueName)

	for {
		err := w.writeToQueueOnce(q, writer)
		if err != nil {
			delay := 15 * time.Second
			log.Printf("WriteToQueue: error received while attempting to write to q -- %s\n", err.Error())
			log.Printf("WriteToQueue: waiting %s until before retrying", delay.String())
			<-time.After(delay)
		}
	}
}

func (w *SQSWriter) writeToQueueOnce(q *gosqs.Queue, writer chan Message) error {
	batch, onComplete := w.assembleSendMessageBatch(writer)

	if len(batch) > 0 {
		log.Printf("%s: sending %d messages to q\n", w.QueueName, len(batch))
		result, err := q.SendMessageBatch(batch)
		if err != nil {
			log.Printf("%#v\n", result)
			return err
		}
		log.Printf("%s: invoking onComplete\n", w.QueueName)
		go onComplete()
	}

	return nil
}

