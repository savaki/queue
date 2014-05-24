package sqs

import (
	"encoding/json"
	"github.com/savaki/queue"
	. "github.com/smartystreets/goconvey/convey"
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
	queueName := "test-queue"
	regionName := "us-west-1"

	var builder *Builder = nil
	var readWriter queue.ReadWriter = nil
	var err error = nil

	Convey("Live Integration Test", t, func() {

		Convey("Given a readWriter", func() {
			builder = New(queueName, regionName)
			builder.Timeout = 1 * time.Second

			go builder.ReadFromQueue()
			go builder.WriteToQueue()

			readWriter, err = builder.BuildReadWriter()
			So(err, ShouldBeNil)

			Convey("When I send a message to the queue", func() {
				expected := map[string]string{"hello": "world", "foo": "bar"}
				LOGGER.Printf("sending message: %+v\n", expected)
				go func() {
					data, _ := json.Marshal(expected)
					readWriter.WriteMessage(data)
				}()

				Convey("Then I expect to receive it back", func() {
					LOGGER.Printf("waiting for message\n")
					select {
					case message := <-readWriter.Receiver():
						// receive
						actual := make(map[string]string)
						json.NewDecoder(message).Decode(&actual)

						// verify
						So(len(actual), ShouldEqual, len(expected))
						So(actual, ShouldResemble, expected)

						// cleanup
						message.Delete()
						<-time.After(builder.Timeout * 2) // ensure that the delete queue has enough time to process the delete

					case <-time.After(time.Second * 15):
						t.Fail()
					}
				})
			})
		})
	})
}
