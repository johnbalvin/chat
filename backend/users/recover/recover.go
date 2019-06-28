package recover

import (
	"context"
	"fmt"
	"chat/backend/errors"
	"chat/backend/users"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"chat/backend/codifications"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

//GetHash codify the hash
func GetHash(rawPass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPass), bcrypt.MinCost)
	if err != nil {
		log.Println("users -> recover:1 -> err:", err)
		return "", err
	}
	return string(hash), nil
}

//Step1 send the recover to the email
func Step1(emailUser string) error {
	user, err := users.ByEmail(emailUser)
	if err != nil {
		log.Println("recover -> Step1:1 -> err:", err)
		return err
	}
	code := codifications.RandomNumbersLetters(5)
	ctx := context.Background()
	recoverPass := Info{Code: code, UserID: user.ID, Created: time.Now().UTC()}
	recoverPass.Expire = recoverPass.Created.Add(time.Minute * 10)
	ref := clienteFS.Collection("RecoverPassStep1").Doc(user.ID)
	if _, err = ref.Set(ctx, recoverPass); err != nil {
		log.Println("recover -> Step1:3 -> err:", err)
		return err
	}
	fmt.Println(code)
	/*	body := strings.Replace(bodyCode, "{{code}}", code, 1) //uncoment if you already configure mailgun
		if err := mg.Send(email, subject, body); err != nil {
			log.Println("register -> Step1:3 -> err:", err)
			return err
	}*/
	return nil
}

//Step2 is the second step to check if the code is ok
func Step2(emailUser, codeUser string) error {
	ctx := context.Background()
	user, err := users.ByEmail(emailUser)
	if err != nil {
		log.Println("recover -> Step2:1 -> err:", err)
		return err
	}
	snapshot, err := clienteFS.Collection("RecoverPassStep1").Doc(user.ID).Get(ctx)
	if err != nil {
		if grpc.Code(err) != codes.NotFound {
			log.Println("recover -> Step2:2 -> err:", err)
			return errors.Empty
		}
		return err
	}
	var recoverPass Info
	if err = snapshot.DataTo(&recoverPass); err != nil {
		log.Println("recover -> Step2:3 -> err:", err)
		return err
	}
	if _, err = snapshot.Ref.Delete(ctx); err != nil {
		log.Println("recover -> Step2:4 -> err:", err)
		return err
	}
	tiempoAhora := time.Now()
	if tiempoAhora.After(recoverPass.Expire) {
		log.Println("recover -> Step2:5 -> err:", errors.TimesOut)
		return errors.TimesOut
	}
	if recoverPass.Code != codeUser {
		log.Println("recover -> Step2:6 -> err:", err)
		return errors.NotTheSame
	}
	recover := Info{UserID: user.ID, Created: tiempoAhora}
	recover.Code = recoverPass.Code
	recover.Expire = recover.Created.Add(time.Minute * 20)
	ref := clienteFS.Collection("RecoverPassStep2").Doc(user.ID)
	if _, err = ref.Set(ctx, recover); err != nil {
		log.Println("recover -> Step2:7 -> err:", err)
		return err
	}
	return nil
}

//Step3 is the last step, is when the user pass all the process to change his key
func Step3(emailUser, codeUser, rawPass string) (users.User, error) {
	ctx := context.Background()
	user, err := users.ByEmail(emailUser)
	if err != nil {
		log.Println("recover -> Step3:1 -> err:", err)
		return users.User{}, err
	}
	snapshot, err := clienteFS.Collection("RecoverPassStep2").Doc(user.ID).Get(ctx)
	if err != nil {
		if grpc.Code(err) != codes.NotFound {
			log.Println("recover -> Step3:2 -> err:", err)
			return users.User{}, errors.Empty
		}
		return users.User{}, err
	}
	var recoverPass Info
	if err = snapshot.DataTo(&recoverPass); err != nil {
		log.Println("recover -> Step3:3 -> err:", err)
		return users.User{}, err
	}
	if _, err = snapshot.Ref.Delete(ctx); err != nil {
		log.Println("recover -> Step3:4 -> err:", err)
		return users.User{}, err
	}
	tiempoAhora := time.Now()
	if tiempoAhora.After(recoverPass.Expire) {
		log.Println("recover -> Step3:5 -> err:", errors.TimesOut)
		return users.User{}, errors.TimesOut
	}
	if recoverPass.Code != codeUser {
		log.Println("recover -> Step3:6 -> err:", err)
		return users.User{}, errors.NotTheSame
	}
	hash, err := GetHash(rawPass)
	if err != nil {
		log.Println("users -> Step3:7 -> err:", err)
		return users.User{}, err
	}
	updates := []firestore.Update{{Path: "Value", Value: hash}}
	if _, err = clienteFS.Collection("P").Doc(user.Ref.ID).Update(ctx, updates); err != nil {
		log.Println("recover -> Step3:8 -> err:", err)
		return users.User{}, err
	}
	return user, nil
}

const bodyCode = "Your code is {{code}}"
const subject = "Account verification"
