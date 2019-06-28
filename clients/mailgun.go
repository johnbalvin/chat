package clients

import (
	"context"
	"log"

	mailgun "github.com/mailgun/mailgun-go/v3"
)

//you should put your own credentials
const publicValidationKey = ""
const domain = ""
const privateAPIKey = ""
const sender = ""

//EmailInfo it's just an auxiliar to send emails
type EmailInfo struct {
	Client mailgun.Mailgun
	Sender string
}

//Mailgun retorna el cliente para trabajar con firestore
func Mailgun() EmailInfo {
	mg := mailgun.NewMailgun(domain, privateAPIKey)
	emailInfo := EmailInfo{Client: mg, Sender: sender}
	return emailInfo
}

//Send sends message to user given an email, subject and the body
func (emailInfo EmailInfo) Send(email, subject, body string) error {
	message := emailInfo.Client.NewMessage(emailInfo.Sender, subject, body, email)
	ctx := context.Background()
	_, _, err := emailInfo.Client.Send(ctx, message)
	if err != nil {
		log.Println("clients -> (emailInfo EmailInfo) Send:1 -> err:", err)
		return err
	}
	//fmt.Printf("ID: %s Resp: %s\n", id, resp)
	return nil
}
