package queue

import (
	"encoding/json"
	"errors"
	"github.com/nu7hatch/gouuid"
	gosqs "github.com/savaki/sqs"
	"log"
	"time"
)

func WriteToQueue(queueName string) {
	writer := &SQSWriter{queueName: queueName}
	writer.WriteToQueue()
}

func (w *SQSWriter) WriteToQueue() {
	if w.messages == nil {
		w.messages = make(chan interface{})
	}
	if w.errs == nil {
		w.errs = make(chan error)
	}
	if w.queueName == "" {
		w.errs <- errors.New("WriteToQueue - misssing queueName")
		return
	}

	q, err := LookupQueue(w.queueName)
	if err != nil {
		panic(err)
	}

	for {
		err := w.writeToQueueOnce(q)
		if err != nil {
			delay := 15 * time.Second
			log.Printf("WriteToQueue: error received while attempting to write to q -- %s\n", err.Error())
			log.Printf("WriteToQueue: waiting %s until before retrying", delay.String())
			<-time.After(delay)
		}
	}
}

func (w *SQSWriter) assembleSendMessageBatch() ([]gosqs.SendMessageBatchRequestEntry, error) {
	requests := make([]gosqs.SendMessageBatchRequestEntry, 0)
	results := make([]interface{}, 0)

	for index := 0; index < w.batchSize; index++ {
		var result interface{} = nil
		select {
		case result = <-w.messages:
			text, err := json.Marshal(result)
			if err != nil {
				return nil, err
			}
			html := string(text)
			log.Printf("assembleSendMessageBatch: queueing message, %d bytes\n", len(html))
			id, err := uuid.NewV4()
			if err != nil {
				return nil, err
			}
			request := gosqs.SendMessageBatchRequestEntry{
				Id:          id.String(),
				MessageBody: html,
			}
			requests = append(requests, request)
			results = append(results, result)

		case <-time.After(5 * time.Minute):
			index = w.batchSize
		}
	}

	return requests, nil
}

func (w *SQSWriter) writeToQueueOnce(q *gosqs.Queue) error {
	batch, err := w.assembleSendMessageBatch()
	if err != nil {
		return err
	}

	if len(batch) > 0 {
		log.Printf("%s: sending %d messages to q\n", w.queueName, len(batch))
		result, err := q.SendMessageBatch(batch)
		if err != nil {
			log.Printf("%#v\n", result)
			return err
		}
	}

	return nil
}
