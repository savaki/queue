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

func (s *sqsQueue) Receiver() <-chan queue.Message {
	return s.Inbound
}

func (s *sqsQueue) WriteMessage(data []byte) {
	s.Outbound <- data
}

func (b *Builder) ReadFromQueue() error {
	w, err := b.build()
	if err != nil {
		return err
	}

	readFromQueueForever(w.Queue, w.Inbound, w.Delete)
	return nil
}

func (b *Builder) DeleteFromQueue() error {
	w, err := b.build()
	if err != nil {
		return err
	}

	return deleteFromQueue(w.Queue, w.Delete, w.Timeout)
}

func (b *Builder) WriteToQueue() error {
	w, err := b.build()
	if err != nil {
		return err
	}

	for {
		err := w.writeToQueueOnce(w.Queue)
		if err != nil {
			delay := 15 * time.Second
			<-time.After(delay)
		}
	}

	return nil
}
