package chat

import (
	"chat/clients"
	"time"
)

//Message contains the message, if its type text, the text is in texto field, is not if it in url
type Message struct {
	From      string `json:"f"`
	Text      string `json:"t"`
	Created   int64  `json:"w"` //unix in nanoseconds
	PhotoUser string `json:"-" firestore:"-"`
	UserName  string ` json:"-" firestore:"-"`
}

var minDuration = time.Second * 1
var clienteFS = clients.Firestore()
