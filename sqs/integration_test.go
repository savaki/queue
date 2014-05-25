package sqs

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	queueName := "test-queue"
	regionName := "us-west-1"

	var client *Client = nil

	Convey("Live Integration Test", t, func() {

		Convey("Given a readWriter", func() {
			client = New(queueName, regionName)
			client.Timeout = 1 * time.Second
			client.BatchSize = 1
			client.Verbose = true
			err := client.Initialize()
			So(err, ShouldBeNil)

			go client.ReadFromQueue()
			go client.WriteToQueue()
			go client.DeleteFromQueue()

			Convey("When I send a message to the queue", func() {
				expected := map[string]string{"hello": "world", "foo": "bar"}
				go func() {
					data, _ := json.Marshal(expected)
					client.Outbound <- data
				}()

				Convey("Then I expect to receive it back", func() {
					select {
					case message := <-client.Inbound:
						log.Println("TEST: Received message")
						// receive
						actual := make(map[string]string)
						err = json.NewDecoder(message).Decode(&actual)
						So(err, ShouldBeNil)

						// verify
						So(len(actual), ShouldEqual, len(expected))
						So(actual, ShouldResemble, expected)

						// cleanup
						message.Delete()
						<-time.After(client.Timeout * 2) // ensure that the delete queue has enough time to process the delete

					case <-time.After(time.Second * 300):
						t.Fail()
					}
				})
			})
		})
	})
}
