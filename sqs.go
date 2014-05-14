package queue

import (
	"encoding/json"
	"github.com/crowdmob/goamz/sqs"
	"log"
	"time"
)

type Message interface {
	OnComplete()
	Unmarshal(v interface{}) error
}

type msg struct {
	data     []byte
	callback func()
}

func (m *msg) Unmarshal(v interface{}) error {
	return json.Unmarshal(m.data, v)
}

func (m *msg) OnComplete() {
	m.callback()
}

type SQSReader struct {
	QueueName  string
	RegionName string
	Queue      *sqs.Queue
	Messages   chan Message
	Del        chan string
	Errs       chan error
	Logger     *log.Logger
	Timeout    time.Duration
	Verbose    bool
}

type SQSWriter struct {
	QueueName  string
	RegionName string
	BatchSize  int
	Messages   chan interface{}
	Errs       chan error
	Logger     *log.Logger
	Timeout    time.Duration
	Verbose    bool
}

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
