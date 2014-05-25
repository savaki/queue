package sqs

import (
	"encoding/base64"
	gosqs "github.com/crowdmob/goamz/sqs"
	uuid "github.com/nu7hatch/gouuid"
	"log"
	"time"
)

func (c *Client) WriteToQueue() error {
	if err := c.Initialize(); err != nil {
		return err
	}

	for {
		err := c.writeToQueueOnce()
		if err != nil {
			if c.Verbose {
				log.Println(err)
			}

			<-time.After(15 * time.Second)
		}
	}

	return nil
}

func (c *Client) assembleSendMessageBatch() ([]gosqs.Message, error) {
	requests := make([]gosqs.Message, 0)

	for index := 0; index < c.BatchSize; index++ {
		var data []byte = nil
		select {
		case data = <-c.Outbound:
			encoded := base64.StdEncoding.EncodeToString(data)
			id, err := uuid.NewV4()
			if err != nil {
				return nil, err
			}
			request := gosqs.Message{
				MessageId: id.String(),
				Body:      encoded,
			}
			requests = append(requests, request)

		case <-time.After(c.Timeout):
			index = c.BatchSize
		}
	}

	return requests, nil
}

func (c *Client) writeToQueueOnce() error {
	batch, err := c.assembleSendMessageBatch()
	if err != nil {
		return err
	}

	if len(batch) > 0 {
		if c.Verbose {
			log.Printf("%s: Sending %d messages\n", c.QueueName, len(batch))
		}
		_, err := c.queue.SendMessageBatch(batch)
		if err != nil {
			return err
		}
	}

	return nil
}
