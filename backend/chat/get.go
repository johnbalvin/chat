package chat

import (
	"chat/backend/users"
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

//GetMsgsBeforeTime get all msgs from a given time
func GetMsgsBeforeTime(when int64, size int) ([]Message, error) {
	ctx := context.Background()
	snapShops, err := clienteFS.Collection("Chat").Where("Created", "<", when).OrderBy("Created", firestore.Desc).Limit(size).Documents(ctx).GetAll()
	if err != nil {
		clienteFS.Collection("query").Add(ctx, map[string]string{"Value": err.Error()})
		log.Println("chats -> GetMsgsBeforeTime:1 -> err:", err)
		return nil, err
	}
	var msgs []Message
	for i := range snapShops {
		var msg Message
		if err = snapShops[i].DataTo(&msg); err != nil {
			log.Println("chats -> GetMsgsBeforeTime:2 -> err:", err)
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

//GetMsgs get only the chats
func GetMsgs(size int) ([]Message, error) {
	ctx := context.Background()
	snapShops, err := clienteFS.Collection("Chat").OrderBy("Created", firestore.Desc).Limit(size).Documents(ctx).GetAll()
	if err != nil {
		clienteFS.Collection("query").Add(ctx, map[string]string{"Value": err.Error()})
		log.Println("chats -> GetMsgsBeforeTime:1 -> err:", err)
		return nil, err
	}
	var msgs []Message
	tam := len(snapShops)
	if tam == 0 {
		return msgs, nil
	}
	userIDs := make(map[string]bool)
	for i := tam - 1; i >= 0; i-- { //the order most be reverted, think about it!
		var msg Message
		if err = snapShops[i].DataTo(&msg); err != nil {
			log.Println("chats -> GetMsgsFromTime:2 -> err:", err)
			return nil, err
		}
		userIDs[msg.From] = true
		msgs = append(msgs, msg)
	}
	users, err := users.ByIDs(userIDs)
	if err != nil {
		log.Println("chats -> GetMsgsFromTime:3 -> err:", err)
		return nil, err
	}
	for i := range msgs {
		for _, user := range users {
			if user.ID == msgs[i].From {
				msgs[i].PhotoUser = user.Photo
				msgs[i].UserName = user.Name
			}
		}
	}
	return msgs, nil
}
