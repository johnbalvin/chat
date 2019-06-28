package register

import (
	"context"
	"fmt"
	"chat/backend/errors"
	"chat/backend/users"
	"math/rand"
	"strconv"
	"time"

	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

//ExisteEmailUser is just to verify is user exist
func ExisteEmailUser(email string) bool {
	ctx := context.Background()
	snapShops, err := clienteFS.Collection("Users").Where("Email", "==", email).Select().Limit(1).Documents(ctx).GetAll()
	if err != nil {
		log.Println("register -> ExisteEmailUser:1 -> err:", err)
		return false
	}
	if len(snapShops) == 0 {
		return false
	}
	return true
}

//Step1 sends the code to the user email
func Step1(email string) error { //obviamente el codigo no retornarlo :V
	if ExisteEmailUser(email) {
		log.Println("register -> Step1:1 -> Duplicado: ", email)
		return errors.Duplicated
	}
	ctx := context.Background()
	recoverCode := Info{Code: strconv.Itoa(rand.Intn(89999) + 100000), Email: email, Created: time.Now()}
	recoverCode.Expire = recoverCode.Created.Add(time.Minute * 10)
	ref := clienteFS.Collection("RegisterStep1").Doc(email)
	if _, err := ref.Set(ctx, recoverCode); err != nil {
		log.Println("register -> Step1:2 -> err:", err)
		return err
	}
	fmt.Println(recoverCode.Code)
	/*	body := strings.Replace(bodyCode, "{{code}}", code, 1) //uncoment if you already configure mailgun
		if err := mg.Send(email, subject, body); err != nil {
			log.Println("register -> Step1:3 -> err:", err)
			return err
	}*/
	return nil
}

//Step2 is the second step to check if the code is ok
func Step2(email, codeUser string) error {
	ctx := context.Background()
	snapshot, err := clienteFS.Collection("RegisterStep1").Doc(email).Get(ctx)
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			log.Println("register -> Step2:1 -> err:", err)
			return errors.Empty
		}
		return err
	}
	var registroCode Info
	if err = snapshot.DataTo(&registroCode); err != nil {
		log.Println("register -> Step2:3 -> err:", err)
		return err
	}
	if _, err = snapshot.Ref.Delete(ctx); err != nil {
		log.Println("register -> Step2:4 -> err:", err)
		return err
	}
	tiempoAhora := time.Now()
	if tiempoAhora.After(registroCode.Expire) {
		log.Println("recover -> Step2:5 -> err:", errors.TimesOut)
		return errors.TimesOut
	}
	if registroCode.Code != codeUser {
		log.Println("register -> Step2:5 -> err:", err)
		return errors.NotTheSame
	}
	registroCode.Created = tiempoAhora
	registroCode.Expire = registroCode.Created.Add(time.Minute * 20)
	ref := clienteFS.Collection("RegisterStep2").Doc(email)
	if _, err = ref.Set(ctx, registroCode); err != nil {
		log.Println("register -> Step2:6 -> err:", err)
		return err
	}
	return nil
}

//Step3 is the last step, is when the user pass all the process to change his key
func Step3(emailUser, codeUser, rawPass, name string) (users.User, error) {
	ctx := context.Background()
	snapshot, err := clienteFS.Collection("RegisterStep2").Doc(emailUser).Get(ctx)
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			log.Println("register -> Step3:1 -> err:", err)
			return users.User{}, errors.Empty
		}
		return users.User{}, err
	}

	var registroCode Info
	if err = snapshot.DataTo(&registroCode); err != nil {
		log.Println("register -> Step3:2 -> err:", err)
		return users.User{}, err
	}
	if _, err = snapshot.Ref.Delete(ctx); err != nil {
		log.Println("register -> Step3:3 -> err:", err)
		return users.User{}, err
	}
	tiempoAhora := time.Now()
	if tiempoAhora.After(registroCode.Expire) {
		log.Println("register -> Step2:4 -> err:", errors.TimesOut)
		return users.User{}, errors.TimesOut
	}
	if registroCode.Code != codeUser {
		log.Println("register -> Step2:5 -> err:", err)
		return users.User{}, errors.Empty
	}
	var user users.User
	user, err = users.New(emailUser, name, rawPass)
	if err != nil {
		log.Println("register -> Step3:6 -> err:", err)
		return users.User{}, err
	}
	return user, nil
}

const bodyCode = "Your code is {{code}}"
const subject = "Account verification"
