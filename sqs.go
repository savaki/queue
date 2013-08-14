package queue

import (
	"encoding/json"
	gosqs "github.com/savaki/sqs"
	"launchpad.net/goamz/aws"
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
	queueName string
	batchSize int
	queue     *gosqs.Queue
	messages  chan Message
	del       chan string
	errs      chan error
}

type SQSWriter struct {
	queueName string
	batchSize int
	messages  chan interface{}
	errs      chan error
}

var (
	// configure our sqs read settings
	DEFAULT_TIMEOUT         time.Duration = 5 * time.Minute
	RECV_ALL                              = []string{"All"}
	RECV_MAX_MESSAGES                     = 10
	RECV_VISIBILITY_TIMEOUT               = 900 // 15 minutes
	DELETE_BATCH_SIZE                     = 10  // delete after theis many messages have been sent
)

func LookupQueue(queueName string) (*gosqs.Queue, error) {
	// login to sqs, reading our credentials from the environment
	auth, err := aws.EnvAuth()
	if err != nil {
		return nil, err
	}

	// connect to our sqs q
	log.Printf("looking up sqs q, %s\n", queueName)
	sqs := gosqs.New(auth, aws.Region{Name: "USWest2", SQSEndpoint: "http://sqs.us-west-2.amazonaws.com"})
	q, err := sqs.GetQueue(queueName)
	if err != nil {
		return nil, err
	}
	log.Printf("%s: ok\n", queueName)

	return q, nil
}
