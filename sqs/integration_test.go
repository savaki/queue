package sqs

import (
	"encoding/json"
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

	var client *Client = nil

	Convey("Live Integration Test", t, func() {

		Convey("Given a readWriter", func() {
			client = New(queueName, regionName)
			client.Timeout = 1 * time.Second

			go client.ReadFromQueue()
			go client.WriteToQueue()

			Convey("When I send a message to the queue", func() {
				expected := map[string]string{"hello": "world", "foo": "bar"}
				LOGGER.Printf("sending message: %+v\n", expected)
				go func() {
					data, _ := json.Marshal(expected)
					client.Outbound <- data
				}()

				Convey("Then I expect to receive it back", func() {
					LOGGER.Printf("waiting for message\n")
					select {
					case message := <-client.Inbound:
						// receive
						actual := make(map[string]string)
						json.NewDecoder(message).Decode(&actual)

						// verify
						So(len(actual), ShouldEqual, len(expected))
						So(actual, ShouldResemble, expected)

						// cleanup
						message.Delete()
						<-time.After(client.Timeout * 2) // ensure that the delete queue has enough time to process the delete

					case <-time.After(time.Second * 15):
						t.Fail()
					}
				})
			})
		})
	})
}
