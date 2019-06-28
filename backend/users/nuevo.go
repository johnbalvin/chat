package users

import (
	"context"
	"chat/backend/codifications"
	"chat/clients"
	"log"
)

//New add user to the database
func New(email, name, password string) (User, error) {
	ctx := context.Background()
	user := User{Name: name, Email: email}
	user.Ref = clienteFS.Collection("Users").NewDoc()
	user.ID = user.Ref.ID
	user.Photo = clients.BucketURL() + "default/user"
	passRef := clienteFS.Collection("P").Doc(user.Ref.ID)
	hash, err := codifications.GetHash(password)
	if err != nil {
		log.Println("users -> New:8 -> err:", err)
		return User{}, err
	}
	batch := clienteFS.Batch().Create(user.Ref, user).Create(passRef, map[string]string{"Value": hash})
	if _, err := batch.Commit(ctx); err != nil {
		log.Println("users -> New:9 -> err:", err)
		return User{}, err
	}
	return user, nil
}
