package sqs

import (
	"errors"
	"github.com/savaki/queue"
	"time"
)

type Builder struct {
	QueueName  string
	RegionName string
	BatchSize  int
	Inbound    chan queue.Message
	Outbound   chan []byte
	Delete     chan string
	Timeout    time.Duration
	Verbose    bool
}

func New(queueName, regionName string) *Builder {
	return &Builder{
		QueueName:  queueName,
		RegionName: regionName,
	}
}

func (b *Builder) build() (*sqsQueue, error) {
	if b.Inbound == nil {
		b.Inbound = make(chan queue.Message)
	}
	if b.Outbound == nil {
		b.Outbound = make(chan []byte)
	}
	if b.Delete == nil {
		b.Delete = make(chan string)
	}
	if b.Timeout == 0 {
		b.Timeout = DEFAULT_TIMEOUT
	}

	if b.QueueName == "" {
		return nil, errors.New("ERROR: QueueName is missing")
	} else if b.RegionName == "" {
		return nil, errors.New("ERROR: RegionName is missing")
	}

	locator := Locator{QueueName: b.QueueName, RegionName: b.RegionName}
	theQueue, err := locator.LookupQueue()
	if err != nil {
		return nil, err
	}

	return &sqsQueue{
		Queue:     theQueue,
		BatchSize: b.BatchSize,
		Inbound:   b.Inbound,
		Outbound:  b.Outbound,
		Delete:    b.Delete,
		Timeout:   b.Timeout,
		Verbose:   b.Verbose,
	}, nil
}

func (b *Builder) BuildReader() (queue.Reader, error) {
	return b.build()
}

func (b *Builder) BuildWriter() (queue.Writer, error) {
	return b.build()
}

func (b *Builder) BuildReadWriter() (queue.ReadWriter, error) {
	return b.build()
}

