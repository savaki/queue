package sqs

import "time"

var (
	// configure our sqs read settings
	DEFAULT_TIMEOUT time.Duration = 1 * time.Minute
	RECV_ALL                      = []string{"All"}
)

const (
	RECV_MAX_MESSAGES       = 10
	RECV_VISIBILITY_TIMEOUT = 900 // 15 minutes
	DELETE_BATCH_SIZE       = 10  // delete after theis many messages have been sent
)

type msg struct {
	data     []byte
	callback func()
}

func (m *msg) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (m *msg) Delete() {
}
