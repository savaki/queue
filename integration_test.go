package queue

import (
	"log"
	"os"
	"testing"
	"time"
)

var (
	// a logger that can be used for testing
	LOGGER = log.New(os.Stdout, "queue:", log.Ldate|log.Ltime)
)

func TestIntegration(t *testing.T) {
	queueName := "queue-development"
	output := make(chan interface{})
	input := make(chan Message)
	timeout := 1 * time.Second

	reader := &SQSReader{QueueName: queueName, Messages: input, Logger: LOGGER, Timeout: timeout}
	go reader.ReadFromQueue()

	writer := &SQSWriter{QueueName: queueName, Messages: output, Logger: LOGGER, Timeout: timeout}
	go writer.WriteToQueue()

	expected := map[string]string{"hello": "world", "foo": "bar"}
	LOGGER.Printf("sending message: %+v\n", expected)
	go func() {
		output <- expected
	}()

	LOGGER.Printf("waiting for message\n")

	select {
	case message := <-input:
		actual := make(map[string]string)
		message.Unmarshal(&actual)
		message.OnComplete()

		if len(actual) != len(expected) {
			t.Fatalf("expected %+v; actual was %+v\n", expected, actual)
		}
		for key, value := range expected {
			if actual[key] != value {
				t.Fatalf("expected %+v; actual was %+v\n", value, actual[key])
			}
		}
		<-time.After(timeout * 2)

	case <-time.After(time.Second * 15):
		t.Fail()
	}
}
