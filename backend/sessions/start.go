package sessions

import (
	"context"
	"chat/backend/codifications"
	"chat/backend/users"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
)

//Start starts a a new session
func Start(trackerID string, user users.User, w http.ResponseWriter, r *http.Request) (users.User, string, error) {
	sessionID := user.SessionID
	if sessionID == "" {
		sessionID = codifications.RandomNumbersLetters(10)
	}
	if err := SetCookie(trackerID, user.ID, sessionID, w); err != nil {
		log.Println("sessions -> Start:1 -> err:", err)
		return users.User{}, "", err
	}
	ref := []firestore.Update{{Path: "SessionID", Value: sessionID}}
	ctx := context.Background()
	if _, err := user.Ref.Update(ctx, ref); err != nil {
		log.Println("sessions -> Start:2 -> err:", err)
		return user, "", err
	}
	return user, sessionID, nil
}
