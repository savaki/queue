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

func SkipTestEnqueueMessage(t *testing.T) {
	original := Sample{"http://www.google.com", []string{"hello"}}
	handle := "1234"

	Convey("Given a Sample message", t, func() {
		// WHEN
		q := &sqsQueue{
			BatchSize: 10,
			Inbound:   make(chan queue.Message),
			Delete:    make(chan string),
		}

		// THEN
		sample := Sample{}
		message := <-q.Inbound
		json.NewDecoder(message).Decode(&sample)

		So(sample.Url, ShouldEqual, original.Url)
		So(sample.Scan, ShouldNotEqual, 1)
		So(sample.Scan[0], ShouldNotEqual, original.Scan[0])

		// TEST 2 - ensure the handle to be sent on the del channel
		actualHandle := ""
		select {
		case actualHandle = <-q.Delete:
		case <-time.After(100 * time.Millisecond):
		}

		if actualHandle != "" {
			t.Fatal("expected handle to not be returned until Delete is called")
		}

		// call OnComplete
		go func() {
			message.Delete()
		}()

		// read from chan
		select {
		case actualHandle = <-q.Delete:
		case <-time.After(100 * time.Millisecond):
		}

		if actualHandle != handle {
			t.Fatal("expected handle to have been set")
		}
	})
}
