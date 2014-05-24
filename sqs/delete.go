package sqs

import (
	gosqs "github.com/crowdmob/goamz/sqs"
	"log"
	"strconv"
	"time"
)

func assembleDeleteMessageBatch(del chan string, timeout time.Duration) []gosqs.Message {
	batch := make([]gosqs.Message, 0)

	for index := 0; index < DELETE_BATCH_SIZE; index++ {
		var handle string = ""
		select {
		case handle = <-del:
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

func deleteFromQueue(q *gosqs.Queue, del chan string, timeout time.Duration) error {
	var err error = nil
	for {
		err = deleteFromQueueOnce(q, del, timeout)
		if err != nil {
			break
		}
	}

	log.Println(err)

	return err
}

func deleteFromQueueOnce(q *gosqs.Queue, del chan string, timeout time.Duration) error {
	batch := assembleDeleteMessageBatch(del, timeout)

	if len(batch) > 0 {
		log.Printf("deleting %d messages from q\n", len(batch))
		_, err := q.DeleteMessageBatch(batch)
		if err != nil {
			return err
		}
	}

	return nil
}
