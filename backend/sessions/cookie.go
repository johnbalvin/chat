package sessions

import (
	"chat/backend/codifications"
	"chat/backend/users"
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/securecookie"
)

var hashKey = []byte("absd9sdfsdfsdfsoipiod1sdfsdfjhm0")
var blockKey = []byte("ab1dfsdf5dssdfsoipiodfsdfsd0jhmz")

//Key is the key to decode and encode cookie
var Key = securecookie.New(hashKey, blockKey)

const cookieDuration = 7776000 // 3 moths, it is in seconds

//ReadCookie just reads the cookie
func ReadCookie(r *http.Request) (string, string, string, error) {
	value := make(map[string]string)
	cookie, err := r.Cookie("u")
	if err != nil {
		return "", "", "", nil
	}
	err = Key.Decode("u", cookie.Value, &value)
	if err != nil {
		log.Println("sessions -> readCookie:1 -> err:", err)
		return "", "", "", err
	}
	return value["t"], value["u"], value["i"], nil
}

//ReadAndAsignCookie reads and assing a cookie with a tracker ID, it returns trackerID,userID,sessionID
func ReadAndAsignCookie(w http.ResponseWriter, r *http.Request) (string, string, string, error) {
	cookie, err := r.Cookie("u")
	if err != nil { //not actually an error but means there is not cookie
		trackerID := codifications.RandomNumbersLetters(15)
		SetCookie(trackerID, "", "", w)
		return trackerID, "", "", nil
	}
	value := make(map[string]string)
	err = Key.Decode("u", cookie.Value, &value)
	if err != nil { //maybe beacause a change in the keys for decoding
		trackerID := codifications.RandomNumbersLetters(15)
		SetCookie(trackerID, "", "", w)
		log.Println("sessions -> ReadAndAsignCookie:1 -> err:", err)
		return trackerID, "", "", nil
	}
	userID := value["u"]
	sessionID := value["i"]
	trackerID := value["t"]
	return trackerID, userID, sessionID, nil
}

//SetCookie set cookie sets the cookie with it's parameteres
func SetCookie(trackerID, userID, ID string, w http.ResponseWriter) error {
	value := make(map[string]string)
	value["t"] = trackerID
	value["u"] = userID
	value["i"] = ID
	encoded, err := Key.Encode("u", value)
	if err != nil {
		log.Println("sessions -> setCookie:1 -> err:", err)
		return err
	}
	cookie := &http.Cookie{
		Name:  "u",
		Value: encoded,
		Path:  "/",
		//Secure: true, //activate laterIMPORTEANTEEEEEEEEEEEEEEEEEE<-----------
		HttpOnly: true,
		MaxAge:   cookieDuration,
	}
	http.SetCookie(w, cookie)
	return nil
}

//Close close the session without removing the tracker ID
func Close(w http.ResponseWriter, r *http.Request) {
	trackerID, userID, sessionID, err := ReadCookie(r)
	if err != nil {
		log.Println("sessions -> Close:1 -> err:", err)
		return
	}
	SetCookie(trackerID, "", "", w)
	if userID == "" {
		return
	}
	go func() {
		user, err := users.ByID(userID)
		if err != nil {
			log.Println("sessions -> Close:2 -> err:", err)
			return
		}
		if user.SessionID != sessionID {
			log.Println("sessions -> Close:3 -> err:", err)
			return
		}
		update := []firestore.Update{{Path: "SessionID", Value: ""}}
		ctx := context.Background()
		if _, err := user.Ref.Update(ctx, update); err != nil {
			log.Println("sessions -> Close:4 -> err:", err)
			return
		}
	}()
}
