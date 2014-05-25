package sqs

import (
	"errors"
	"fmt"
	goaws "github.com/crowdmob/goamz/aws"
	gosqs "github.com/crowdmob/goamz/sqs"
	"log"
)

type Locator struct {
	QueueName  string
	RegionName string
}

func (locator Locator) LookupQueue() (*gosqs.Queue, error) {
	// login to sqs, reading our credentials from the environment
	auth, err := goaws.EnvAuth()
	if err != nil {
		return nil, err
	}

	// connect to our sqs q
	region, found := goaws.Regions[locator.RegionName]
	if !found {
		return nil, errors.New(fmt.Sprintf("no such region, '%s'", locator.RegionName))
	}
	log.Printf("Looking up sqs queue by name, '%s'\n", locator.QueueName)
	sqsService := gosqs.New(auth, region)
	q, err := sqsService.GetQueue(locator.QueueName)
	if err != nil {
		return nil, err
	}
	log.Printf("Found %s => %s: ok\n", locator.QueueName, q.Url)

	return q, nil
}
