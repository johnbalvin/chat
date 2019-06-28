package register

import (
	"chat/clients"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

var mg = clients.Mailgun()

var clienteFS = clients.Firestore()

//Info is the struct to recover the pass
type Info struct {
	Code    string
	Email   string
	Created time.Time
	Expire  time.Time
}
