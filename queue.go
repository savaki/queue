package queue

import (
	"io"
)

type Message interface {
	io.Reader

	Delete()
}
