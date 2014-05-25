package sqs

import (
	gosqs "github.com/crowdmob/goamz/sqs"
	"github.com/savaki/queue"
	"log"
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
	if err := c.Initialize(); err != nil {
		return err
	}

	// loop forever receiving messages
	for {
		err := c.readFromQueueOnce()
		if err != nil {
			if c.Verbose {
				log.Println(err)
			}
			<-time.After(15 * time.Second)
		}
	}
	return nil
}

func (c *Client) DeleteFromQueue() error {
	if err := c.Initialize(); err != nil {
		return err
	}

	err := c.deleteFromQueue()
	if c.Verbose {
		log.Println(err)
	}
	return err
}
