package clients

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

//Firestore returns client to work with firestore
func Firestore() *firestore.Client {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, ProjectID)
	if err != nil {
		log.Fatal("clients -> Firestore -> err:", err)
	}
	return client
}
