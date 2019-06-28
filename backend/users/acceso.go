package users

import (
	"chat/backend/errors"
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"
)

//PassCompare just compare the pass
func (user *User) PassCompare(posiblePass string) (bool, error) {
	ctx := context.Background()
	snashop, err := clienteFS.Collection("P").Doc(user.Ref.ID).Get(ctx)
	if err != nil {
		log.Println("users -> (user *User) PassCompare:1 -> err:", err)
		return false, err
	}
	passRaw, err := snashop.DataAt("Value")
	if err != nil {
		log.Println("users -> (user *User) PassCompare:2 -> err:", err)
		return false, err
	}
	pass, ok := passRaw.(string)
	if !ok {
		log.Println("users -> (user *User) PassCompare:2 -> err:", err)
		return false, errors.Security
	}
	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(posiblePass))
	if err == nil {
		return true, nil
	}
	return false, nil
}
