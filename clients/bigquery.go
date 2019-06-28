package clients

import (
	"context"
	"log"

	"cloud.google.com/go/bigquery"
)

//Bigquery returns client for bigquery
func Bigquery() *bigquery.Client {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, ProjectID)
	if err != nil {
		log.Fatal("clients -> Bigquery -> err:", err)
	}
	return client
}
