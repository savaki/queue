package queue

import (
	"encoding/json"
	"errors"
	"github.com/crowdmob/goamz/sqs"
	"github.com/nu7hatch/gouuid"
	"io/ioutil"
	"log"
	"time"
)

func WriteToQueue(queueName string, regionName string, messages chan interface{}) {
	writer := &SQSWriter{
		QueueName:  queueName,
		RegionName: regionName,
		Messages:   messages,
	}
	writer.WriteToQueue()
}

func (w *SQSWriter) WriteToQueue() {
	if w.Messages == nil {
		w.Messages = make(chan interface{})
	}
	if w.Errs == nil {
		w.Errs = make(chan error)
	}
	if w.QueueName == "" {
		w.Errs <- errors.New("WriteToQueue - misssing queueName")
		return
	}
	if w.RegionName == "" {
		w.Errs <- errors.New("WriteToQueue - misssing regionName")
		return
	}
	if w.Logger == nil {
		w.Logger = log.New(ioutil.Discard, "queue", log.Ldate|log.Ltime)
	}
	if w.Timeout == 0 {
		w.Timeout = DEFAULT_TIMEOUT
	}
	if w.BatchSize == 0 {
		w.BatchSize = 1
	}

	locator := Locator{w.QueueName, w.RegionName}
	q, err := locator.LookupQueue()
	if err != nil {
		w.Logger.Printf("ERROR!  No queue with name, %s\n", w.QueueName)
		w.Errs <- err
		return
	}

	for {
		err := w.writeToQueueOnce(q)
		if err != nil {
			delay := 15 * time.Second
			w.Logger.Printf("WriteToQueue: error received while attempting to write to q -- %s\n", err.Error())
			w.Logger.Printf("WriteToQueue: waiting %s until before retrying", delay.String())
			<-time.After(delay)
		}
	}
}

func (w *SQSWriter) assembleSendMessageBatch() ([]sqs.Message, error) {
	requests := make([]sqs.Message, 0)
	results := make([]interface{}, 0)

	for index := 0; index < w.BatchSize; index++ {
		var result interface{} = nil
		select {
		case result = <-w.Messages:
			text, err := json.Marshal(result)
			if err != nil {
				return nil, err
			}
			html := string(text)
			w.Logger.Printf("assembleSendMessageBatch: queueing message, %d bytes\n", len(html))
			id, err := uuid.NewV4()
			if err != nil {
				return nil, err
			}
			request := sqs.Message{
				MessageId: id.String(),
				Body:      html,
			}
			requests = append(requests, request)
			results = append(results, result)

		case <-time.After(w.Timeout):
			index = w.BatchSize
		}
	}

	return requests, nil
}

func (w *SQSWriter) writeToQueueOnce(q *sqs.Queue) error {
	batch, err := w.assembleSendMessageBatch()
	if err != nil {
		return err
	}

	if len(batch) > 0 {
		w.Logger.Printf("%s: sending %d messages to q\n", w.QueueName, len(batch))
		result, err := q.SendMessageBatch(batch)
		if err != nil {
			w.Logger.Printf("%#v\n", result)
			return err
		}
	}

	return nil
}
