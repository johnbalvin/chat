package users

import (
	"chat/clients"
	"time"

	"cloud.google.com/go/firestore"
)

var clienteFS = clients.Firestore()

var mg = clients.Mailgun()

//User contains user information
type User struct {
	ID              string `json:"-"`
	Photo           string `json:"p"`
	Name            string `json:"n"`
	Email           string `json:"-"`
	EstaBloqueado   bool   `json:"-"`
	SessionID       string
	LastMessageSent time.Time
	Ref             *firestore.DocumentRef `json:"-"`
}
