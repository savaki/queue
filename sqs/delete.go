package sqs

import (
	gosqs "github.com/crowdmob/goamz/sqs"
	"log"
	"strconv"
	"time"
)

func assembleDeleteMessageBatch(delete chan string, timeout time.Duration) []gosqs.Message {
	batch := make([]gosqs.Message, 0)

	for index := 0; index < DELETE_BATCH_SIZE; index++ {
		var handle string = ""
		select {
		case handle = <-delete:
			log.Printf("assembleDeleteMessageBatch: received handle, %s\n", handle)
			message := gosqs.Message{
				MessageId:     strconv.Itoa(index + 1),
				ReceiptHandle: handle,
			}
			batch = append(batch, message)

		case <-time.After(timeout):
			return batch
		}
	}

	return batch
}

func (c *Client) deleteFromQueue() error {
	var err error = nil
	for {
		err = c.deleteFromQueueOnce()
		if err != nil {
			break
		}
	}

	return err
}

func (c *Client) deleteFromQueueOnce() error {
	batch := assembleDeleteMessageBatch(c.Delete, c.Timeout)

	if len(batch) > 0 {
		_, err := c.queue.DeleteMessageBatch(batch)
		if err != nil {
			return err
		}

		if c.Verbose {
			log.Printf("%s: Deleted %d messages\n", c.QueueName, len(batch))
		}
	}

	return nil
}
