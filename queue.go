package queue

import (
	"io"
)

type Client interface {
	Name() string
	Inbox() chan Message
	Outbox() chan []byte
	ReadFromQueue() error
	WriteToQueue() error
	DeleteFromQueue() error
}

type Message interface {
	io.Reader

	Delete()
}
