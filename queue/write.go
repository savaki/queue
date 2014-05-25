package main

import (
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func WriteCommand(context *cli.Context) {
	client := getClient(context)

	go client.WriteToQueue()

	filename := context.String("file")
	if filename == "" {
		log.Fatalf("ERROR: no file specified")
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("ERROR: unable to read file, %s => %s\n", filename, err)
	}

	log.Printf("%s: Queuing file, %s", client.Name(), filename)
	client.Outbox() <- data

	// let the message get sent ... this is a terrible way to do this
	<-time.After(5 * time.Second)
	os.Exit(0)
}
