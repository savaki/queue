package queue

import (
	"io"
)

type Reader interface {
	Receiver() <-chan Message
}

type Writer interface {
	WriteMessage(data []byte)
}

type ReadWriter interface {
	Reader
	Writer
}

type Message interface {
	io.Reader

	Delete()
}
