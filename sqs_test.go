package queue

import (
	"encoding/json"
	"testing"
	"time"
)

type Sample struct {
	Url  string   `json:"url"`
	Scan []string `json:"scan"`
}

func TestEnqueueMessage(t *testing.T) {
	// GIVEN
	original := Sample{"http://www.google.com", []string{"hello"}}
	bytes, _ := json.Marshal(original)
	handle := "1234"
	messages := make(chan Message)
	del := make(chan string)

	// WHEN
	reader := &SQSReader{Messages: messages, Del: del, Logger: LOGGER}
	go reader.enqueueMessage(string(bytes), handle)

	// THEN
	sample := Sample{}
	message := <-reader.Messages
	message.Unmarshal(&sample)

	if sample.Url != original.Url {
		t.Fatalf("expected url to be %s; actual was %s", original.Url, sample.Url)
	}
	if len(sample.Scan) != 1 {
		t.Fatal("expected scan to be set")
	}
	if sample.Scan[0] != original.Scan[0] {
		t.Fatalf("expected scan to be set to %s; actual was %s", original.Scan[0], sample.Scan[0])
	}

	// TEST 2 - ensure the handle to be sent on the del channel
	actualHandle := ""
	select {
	case actualHandle = <-reader.Del:
	case <-time.After(100 * time.Millisecond):
	}

	if actualHandle != "" {
		t.Fatal("expected handle to not be returned until OnComplete is called")
	}

	// call OnComplete
	go func() {
		message.OnComplete()
	}()

	// read from chan
	select {
	case actualHandle = <-reader.Del:
	case <-time.After(100 * time.Millisecond):
	}

	if actualHandle != handle {
		t.Fatal("expected handle to have been set")
	}
}
