package chat

import (
	"chat/backend/errors"
	"chat/backend/users"
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
)

//AddMsg adds a meesages giving a user and a the message itslelf
func AddMsg(user users.User, text string) (Message, error) {
	if len(text) >= 500 {
		log.Println("chat -> AddMsg:1 -> text: ", text)
		return Message{}, errors.NotAllow
	}
	now := time.Now()
	if now.Sub(user.LastMessageSent) < minDuration {
		log.Println("chat -> AddMsg:2 -> userID: ", user.ID)
		return Message{}, errors.TimesOut
	}
	ctx := context.Background()
	ref := []firestore.Update{{Path: "LastMessageSent", Value: now}}
	message := Message{From: user.ID, Text: text, Created: now.UnixNano()}
	msgRef := clienteFS.Collection("Chat").NewDoc()
	batch := clienteFS.Batch().Create(msgRef, message).Update(user.Ref, ref)
	if _, err := batch.Commit(ctx); err != nil {
		log.Println("chat -> AddMsg:3 -> err:", err)
		return Message{}, err
	}
	return message, nil
}
