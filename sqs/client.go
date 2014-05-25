package sqs

import (
	"errors"
	gosqs "github.com/crowdmob/goamz/sqs"
	"github.com/savaki/queue"
	"time"
)

type Client struct {
	queue      *gosqs.Queue
	QueueName  string
	RegionName string
	BatchSize  int
	Inbound    chan queue.Message
	Outbound   chan []byte
	Delete     chan string
	Timeout    time.Duration
	Verbose    bool
}

func New(queueName, regionName string) *Client {
	return &Client{
		QueueName:  queueName,
		RegionName: regionName,
	}
}

func (c *Client) initialize() error {
	if c.queue != nil {
		return nil
	}

	if c.Inbound == nil {
		c.Inbound = make(chan queue.Message)
	}
	if c.Outbound == nil {
		c.Outbound = make(chan []byte)
	}
	if c.Delete == nil {
		c.Delete = make(chan string)
	}
	if c.Timeout == 0 {
		c.Timeout = DEFAULT_TIMEOUT
	}

	if c.QueueName == "" {
		return errors.New("ERROR: QueueName is missing")

	} else if c.RegionName == "" {
		return errors.New("ERROR: RegionName is missing")
	}

	locator := Locator{QueueName: c.QueueName, RegionName: c.RegionName}
	theQueue, err := locator.LookupQueue()
	if err != nil {
		return err
	}

	c.queue = theQueue

	return nil
}
