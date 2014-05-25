package sqs

import (
	"encoding/json"
	gosqs "github.com/crowdmob/goamz/sqs"
	uuid "github.com/nu7hatch/gouuid"
	"time"
)

func assembleSendMessageBatch(outbound chan []byte, batchSize int, timeout time.Duration) ([]gosqs.Message, error) {
	requests := make([]gosqs.Message, 0)
	results := make([]interface{}, 0)

	for index := 0; index < batchSize; index++ {
		var result []byte = nil
		select {
		case result = <-outbound:
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

		case <-time.After(timeout):
			index = batchSize
		}
	}

	return requests, nil
}

func writeToQueueOnce(q *gosqs.Queue, outbound chan []byte, batchSize int, timeout time.Duration) error {
	batch, err := assembleSendMessageBatch(outbound, batchSize, timeout)
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
