package main

import (
	"errors"
	"github.com/codegangsta/cli"
	"github.com/savaki/queue"
	"github.com/savaki/queue/sqs"
	"log"
	"time"
)

func getClient(context *cli.Context) queue.Client {
	if context.String("provider") == "sqs" {
		queueName := context.String("name")
		region := context.String("region")
		client := sqs.New(queueName, region)

		client.Timeout = 1 * time.Second
		if context.Bool("verbose") {
			client.Verbose = true
		}

		err := client.Initialize()
		if err != nil {
			log.Fatalf("ERROR: unable to connect to queue, %s:%s => %v\n", queueName, region, err)
		}

		return client
	}

	log.Fatalln(errors.New("ERROR: only the sqs client is currently supported"))
	return nil
}
