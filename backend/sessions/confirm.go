package sessions

import (
	"chat/backend/errors"
	"chat/backend/users"
	"log"
)

//Confirm confirms is to confirm if session is valid, checkid is id from cookie and id from database
func Confirm(userID, sessionID string) (users.User, error) {
	user, err := users.ByID(userID)
	if err != nil {
		log.Println("sessions -> Confirm:1 -> err:", err)
		return users.User{}, err
	}
	if user.SessionID != sessionID {
		log.Println("sessions -> Confirm:2 -> session.ID:", user.SessionID, " sessionID:", sessionID)
		return users.User{}, errors.NotTheSame
	}
	return user, nil
}
