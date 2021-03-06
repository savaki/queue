package sqs

import (
	"io"
	"time"
)

var (
	// configure our sqs read settings
	DEFAULT_TIMEOUT time.Duration = 1 * time.Minute
	RECV_ALL                      = []string{"All"}
)

const (
	RECV_MAX_MESSAGES       = 10
	RECV_VISIBILITY_TIMEOUT = 900 // 15 minutes
	DELETE_BATCH_SIZE       = 10  // delete after this many messages have been sent
	DEFAULT_BATCH_SIZE      = 10
)

type msg struct {
	reader   io.Reader
	callback func()
	deleted  bool
}

func (m *msg) Read(p []byte) (n int, err error) {
	return m.reader.Read(p)
}

func (m *msg) Delete() {
	if !m.deleted {
		m.callback()
		m.deleted = true
	}
}
