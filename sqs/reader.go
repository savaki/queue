package sqs

import (
	"encoding/base64"
	"github.com/savaki/queue"
	"log"
	"strings"
)

func (c *Client) readFromQueueOnce() error {
	if c.Verbose {
		log.Printf("Waiting for message (%s:%s)\n", c.QueueName, c.RegionName)
	}

	results, err := c.queue.ReceiveMessage(RECV_MAX_MESSAGES)
	if err != nil {
		if c.Verbose {
			log.Println(err)
		}
		return err
	}

	if c.Verbose {
		log.Printf("%s: Received %d messages\n", c.QueueName, len(results.Messages))
	}
	for _, message := range results.Messages {
		enqueueMessage(c.Inbound, c.Delete, message.Body, message.ReceiptHandle)
	}

	return nil
}

func enqueueMessage(inbound chan queue.Message, delete chan string, text string, handle string) {
	var message queue.Message = &msg{
		reader: base64.NewDecoder(base64.StdEncoding, strings.NewReader(text)),
		callback: func() {
			delete <- handle
		},
	}
	inbound <- message
}
