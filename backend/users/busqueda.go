package users

import (
	"context"
	"chat/backend/errors"
	"log"

	"cloud.google.com/go/firestore"
)

//ByID look for the user by id
func ByID(userID string) (User, error) {
	ctx := context.Background()
	snapShops, err := clienteFS.Collection("Users").Where("ID", "==", userID).Limit(1).Documents(ctx).GetAll()
	if err != nil {
		log.Println("users -> ByID:1 -> err:", err)
		return User{}, err
	}
	if err != nil {
		log.Println("users -> ByID:2 -> err:", err)
		return User{}, err
	}
	if len(snapShops) == 0 {
		return User{}, errors.Empty
	}
	var user User
	if err := snapShops[0].DataTo(&user); err != nil {
		log.Println("users -> ByID:4 -> err:", err)
		return User{}, err
	}
	return user, nil
}

//ByIDs look for the user by id
func ByIDs(userIDs map[string]bool) ([]User, error) {
	var refs []*firestore.DocumentRef
	for userID := range userIDs {
		refs = append(refs, clienteFS.Collection("Users").Doc(userID))
	}
	ctx := context.Background()
	snapshots, err := clienteFS.GetAll(ctx, refs)
	if err != nil {
		log.Println("users -> ByIDs:1 -> err:", err)
		return nil, err
	}
	var users []User
	for _, snap := range snapshots {
		var user User
		if !snap.Exists() { //posible a user just sent a post with fake ids
			log.Println("users -> ByIDs:2 -> err:", err)
			continue
		}
		if err = snap.DataTo(&user); err != nil {
			log.Println("users -> ByIDs:3 -> err:", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

//ByEmail look for the user by email
func ByEmail(email string) (User, error) {
	ctx := context.Background()
	log.Println(len(email), email)
	snapShops, err := clienteFS.Collection("Users").Where("Email", "==", email).Limit(1).Documents(ctx).GetAll()
	if err != nil {
		log.Println("users -> ByEmail:1 -> err:", err)
		return User{}, err
	}
	if len(snapShops) == 0 {
		return User{}, errors.Empty
	}
	var user User
	if err = snapShops[0].DataTo(&user); err != nil {
		log.Println("users -> ByEmail:2 -> err:", err)
		return User{}, err
	}
	return user, nil
}
