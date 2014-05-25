package main

import (
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
	MAX_FILE_SIZE = 32768
)

func DirectoryCommand(context *cli.Context) {
	client := getClient(context)

	go client.WriteToQueue()

	dir := context.String("dir")
	if dir == "" {
		log.Fatalf("ERROR: no directory specified")
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("ERROR: unable to read directory, %s => %s\n", dir, err)
	}

	for _, file := range files {
		if file.IsDir() || strings.HasPrefix(".", file.Name()) {
			continue
		}

		if file.Size() > MAX_FILE_SIZE {
			log.Printf("%s: skipping.  file size exceeds maximum file size, %d\n", file.Name(), MAX_FILE_SIZE)
			continue
		}

		data, err := ioutil.ReadFile(file.Name())
		if err != nil {
			log.Fatalf("ERROR: unable to read file, %s => %s\n", file.Name(), err)
		}

		log.Printf("%s: Queuing file, %s", client.Name(), file.Name())
		client.Outbox() <- data
	}

	// let the message get sent ... this is a terrible way to do this
	<-time.After(5 * time.Second)
	os.Exit(0)
}
