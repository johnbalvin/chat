package clients

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
)

//Storage returns the google cloud storage client
func Storage() *storage.BucketHandle {
	ctx := context.Background()
	cliente, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal("clients -> Storage -> err:", err)
	}
	return cliente.Bucket(bucketName)
}

//BucketURL returns the first full path of any url file
func BucketURL() string {
	return url
}
