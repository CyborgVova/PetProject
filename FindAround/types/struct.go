package types

import (
	bf "github.com/cenkalti/backoff/v4"

	"log"
	"time"

	es "github.com/elastic/go-elasticsearch/v7"
)

type Place struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Location struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
	} `json:"location"`
}

func CreateEsClient() *es.Client {
	retryBackoff := bf.NewExponentialBackOff()
	es, err := es.NewClient(es.Config{
		Addresses: []string{
			"http://elasticsearch:9200",
		},
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		MaxRetries: 5,
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	return es
}
