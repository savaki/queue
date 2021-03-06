package sqs

import (
	"errors"
	"fmt"
	gosqs "github.com/crowdmob/goamz/sqs"
	"github.com/savaki/queue"
	"log"
	"sync"
	"time"
)

type Client struct {
	mutex      sync.Mutex
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

func (c *Client) Name() string {
	return fmt.Sprintf("%s:%s", c.QueueName, c.RegionName)
}

func (c *Client) Inbox() chan queue.Message {
	return c.Inbound
}

func (c *Client) Outbox() chan []byte {
	return c.Outbound
}

func (c *Client) Initialize() error {
	if c.queue != nil {
		return nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

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
	if c.BatchSize == 0 {
		c.BatchSize = DEFAULT_BATCH_SIZE
	}

	if c.QueueName == "" {
		err := errors.New("ERROR: QueueName is missing")
		if c.Verbose {
			log.Println(err)
		}
		return err

	} else if c.RegionName == "" {
		err := errors.New("ERROR: RegionName is missing")
		if c.Verbose {
			log.Println(err)
		}
		return err
	}

	locator := Locator{
		QueueName:  c.QueueName,
		RegionName: c.RegionName,
		Verbose:    c.Verbose,
	}
	theQueue, err := locator.LookupQueue()
	if err != nil {
		if c.Verbose {
			log.Println(err)
		}
		return err
	}

	c.queue = theQueue

	return nil
}
