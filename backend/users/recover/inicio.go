package recover

import (
	"chat/clients"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

var clienteFS = clients.Firestore()

var mg = clients.Mailgun()

//Info is the struct to recover the pass
type Info struct {
	Code    string
	UserID  string
	Created time.Time
	Expire  time.Time
}
