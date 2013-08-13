package queue

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

type SampleMessage struct {
	Url      string   `json:"url"`
	Scan     []string `json:"scan"`
	callback func()
}

func (s SampleMessage) OnComplete() {
	s.callback()
}

func NewSampleMessage(callback func()) Message {
	return &SampleMessage{callback: callback}
}

func TestEnqueuesRequest(t *testing.T) {
	// GIVEN
	url := "http://www.google.com"
	scan := "hello"
	text := fmt.Sprintf(`{"url":"%s", "scan":["%s"]}`, url, scan)
	handle := "1234"
	reader := make(chan Message)
	del := make(chan string)

	// WHEN
	q, err := NewReader("queueName", NewSampleMessage)
	if err != nil {
		panic(err)
	}
	go q.enqueueRequest(text, handle, reader, del)

	// THEN
	request := <-reader
	message := request.(*SampleMessage) // convert type back to what we expect

	if message.Url != url {
		t.Fatalf("expected url to be %s; actual was %s", url, message.Url)
	}
	if len(message.Scan) != 1 {
		t.Fatal("expected scan to be set")
	}
	if message.Scan[0] != scan {
		t.Fatalf("expected scan to be set to %s; actual was %s", scan, message.Scan[0])
	}

	actualHandle := ""
	select {
	case actualHandle = <-del:
	case <-time.After(100 * time.Millisecond):
	}

	if actualHandle != "" {
		t.Fatal("expected handle to not be returned until OnComplete is called")
	}

	// call OnComplete
	go func() {
		request.OnComplete()
	}()

	// read from chan
	select {
	case actualHandle = <-del:
	case <-time.After(100 * time.Millisecond):
	}

	if actualHandle != handle {
		t.Fatal("expected handle to have been set")
	}
}

func TestAssembleDeleteMessageBatch(t *testing.T) {
	// GIVEN
	del := make(chan string)
	count := 7
	go func() {
		for i := 0; i < count; i++ {
			del <- strconv.Itoa(i)
		}
	}()

	// WHEN
	batch := assembleDeleteMessageBatch(del, 100*time.Millisecond)
	if len(batch) != count {
		t.Fatalf("expected %d elements in our batch; actual was %d", count, len(batch))
	}

	for i := 0; i < count; i++ {
		if batch[i].Id != strconv.Itoa(i+1) {
			t.Fatal("expected id to be the index + 1")
		}

		if batch[i].ReceiptHandle != strconv.Itoa(i) {
			t.Fatal("expected the handle to be the index value (0 indexed)")
		}
	}
}

