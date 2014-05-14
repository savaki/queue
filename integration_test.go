package queue

import (
	"log"
	"os"
	"testing"
	"time"
	"fmt"
)

var (
	// a logger that can be used for testing
	LOGGER = log.New(os.Stdout, "queue:", log.Ldate|log.Ltime)
)

func TestIntegration(t *testing.T) {
	queueName := "test-queue"
	regionName := "us-west-1"

	output := make(chan interface{})
	input := make(chan Message)
	timeout := 1 * time.Second

	reader := &SQSReader{
		QueueName: queueName,
		RegionName: regionName,
		Messages: input,
		Timeout: timeout,
		Verbose: true,
	}
	go reader.ReadFromQueue()

	writer := &SQSWriter{
		QueueName: queueName,
		RegionName: regionName,
		Messages: output,
		Timeout: timeout,
		Verbose: true,
	}
	go writer.WriteToQueue()

	go func() {
		select {
		case err := <-reader.Errs:
			fmt.Println(err)
		case err := <-writer.Errs:
			fmt.Println(err)
		}
	}()

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
		Verify(t, expected, actual)
		<-time.After(timeout * 2) // ensure that the delete queue has enough time to process the delete

	case <-time.After(time.Second * 15):
		t.Fail()
	}
}

func Verify(t *testing.T, expected, actual map[string]string) {
	if len(actual) != len(expected) {
		t.Fatalf("expected %+v; actual was %+v\n", expected, actual)
	}
	for key, value := range expected {
		if actual[key] != value {
			t.Fatalf("expected %+v; actual was %+v\n", value, actual[key])
		}
	}
}
