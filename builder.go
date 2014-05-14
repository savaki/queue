package queue

import (
	"errors"
	"fmt"
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/sqs"
	"log"
)

type Locator struct {
	QueueName  string
	RegionName string
}

func (locator Locator) LookupQueue() (*sqs.Queue, error) {
	// login to sqs, reading our credentials from the environment
	auth, err := aws.EnvAuth()
	if err != nil {
		return nil, err
	}

	// connect to our sqs q
	region, found := aws.Regions[locator.RegionName]
	if !found {
		return nil, errors.New(fmt.Sprintf("no such region, '%s'", locator.RegionName))
	}
	log.Printf("looking up sqs queue by name, '%s'\n", locator.QueueName)
	sqsService := sqs.New(auth, region)
	q, err := sqsService.GetQueue(locator.QueueName)
	if err != nil {
		return nil, err
	}
	log.Printf("%s: ok\n", locator.QueueName)

	return q, nil
}
