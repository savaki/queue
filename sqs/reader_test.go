package sqs

import (
	"encoding/json"
	"github.com/savaki/queue"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

type Sample struct {
	Url  string   `json:"url"`
	Scan []string `json:"scan"`
}

func TestEnqueueMessage(t *testing.T) {
	original := Sample{"http://www.google.com", []string{"hello"}}
	bytes, _ := json.Marshal(original)
	handle := "1234"

	Convey("Given a Sample message", t, func() {
		// WHEN
		builder := &sqsQueue
			builder.build()
		go builder.ReadFromQueue()

		// THEN
		sample := Sample{}
		message := <-reader.Inbound
		json.NewDecoder(message).Decode(&sample)

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
		case actualHandle = <-reader.Delete:
		case <-time.After(100 * time.Millisecond):
		}

		if actualHandle != "" {
			t.Fatal("expected handle to not be returned until OnComplete is called")
		}

		// call OnComplete
		go func() {
			message.Delete()
		}()

		// read from chan
		select {
		case actualHandle = <-reader.Delete:
		case <-time.After(100 * time.Millisecond):
		}

		if actualHandle != handle {
			t.Fatal("expected handle to have been set")
		}
	})
}
