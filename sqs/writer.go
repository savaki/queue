package sqs

import (
	"encoding/json"
	gosqs "github.com/crowdmob/goamz/sqs"
	uuid "github.com/nu7hatch/gouuid"
	"time"
)

func (w *sqsQueue) assembleSendMessageBatch() ([]gosqs.Message, error) {
	requests := make([]gosqs.Message, 0)
	results := make([]interface{}, 0)

	for index := 0; index < w.BatchSize; index++ {
		var result interface{} = nil
		select {
		case result = <-w.Inbound:
			text, err := json.Marshal(result)
			if err != nil {
				return nil, err
			}
			html := string(text)
			id, err := uuid.NewV4()
			if err != nil {
				return nil, err
			}
			request := gosqs.Message{
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

func (w *sqsQueue) writeToQueueOnce(q *gosqs.Queue) error {
	batch, err := w.assembleSendMessageBatch()
	if err != nil {
		return err
	}

	if len(batch) > 0 {
		_, err := q.SendMessageBatch(batch)
		if err != nil {
			return err
		}
	}

	return nil
}
