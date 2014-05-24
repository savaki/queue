package sqs

import (
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
	"time"
)

func TestAssembleDeleteMessageBatch(t *testing.T) {
	Convey("Given a channel to handle deletes", t, func() {
		delete := make(chan string)
		count := 7
		go func() {
			for i := 0; i < count; i++ {
				delete <- strconv.Itoa(i)
			}
		}()

		Convey("When I call #assembleDeleteMessageBatch", func() {
			batch := assembleDeleteMessageBatch(delete, 100*time.Millisecond)

			Convey("Then I expect the same number of entries to be returned", func() {
				So(len(batch), ShouldEqual, count)
			})

			Convey("And I expect the messages to be properly constructed", func() {
				for i := 0; i < count; i++ {
					if batch[i].MessageId != strconv.Itoa(i+1) {
						t.Fatal("expected id to be the index + 1")
					}

					if batch[i].ReceiptHandle != strconv.Itoa(i) {
						t.Fatal("expected the handle to be the index value (0 indexed)")
					}
				}
			})
		})
	})
}
