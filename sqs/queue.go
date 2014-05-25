package sqs

import (
	gosqs "github.com/crowdmob/goamz/sqs"
	"github.com/savaki/queue"
	"time"
)

type sqsQueue struct {
	BatchSize int
	Queue     *gosqs.Queue
	Inbound   chan queue.Message
	Outbound  chan []byte
	Delete    chan string
	Timeout   time.Duration
	Verbose   bool
}

func (c *Client) ReadFromQueue() error {
	if err := c.initialize(); err != nil {
		return err
	}

	readFromQueueForever(c.queue, c.Inbound, c.Delete)
	return nil
}

func (c *Client) DeleteFromQueue() error {
	if err := c.initialize(); err != nil {
		return err
	}

	return deleteFromQueue(c.queue, c.Delete, c.Timeout)
}

func (c *Client) WriteToQueue() error {
	if err := c.initialize(); err != nil {
		return err
	}

	for {
		err := writeToQueueOnce(c.queue, c.Outbound, c.BatchSize, c.Timeout)
		if err != nil {
			delay := 15 * time.Second
			<-time.After(delay)
		}
	}

	return nil
}
