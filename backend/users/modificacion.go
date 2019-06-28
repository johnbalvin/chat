package users

import (
	"context"
	"chat/backend/savefile"
	"log"
	"mime/multipart"

	"cloud.google.com/go/firestore"
)

//ChangeName changes user display name
func (user *User) ChangeName(nombre string) error {
	ref := []firestore.Update{{Path: "Name", Value: nombre}}
	ctx := context.Background()
	if _, err := user.Ref.Update(ctx, ref); err != nil {
		log.Println("users -> (user *User) ChangeName:1 -> err:", err)
		return err
	}
	return nil
}

//ChangesPhoto cchanges the user photo
func (user *User) ChangesPhoto(file multipart.File, header *multipart.FileHeader) (string, error) {
	_, _, _, url, err := savefile.AddImg("users/"+user.Ref.ID+"/", file, header)
	if err != nil {
		log.Println("users -> (user *User) ChangesPhoto:1 -> err:", err)
		return "", err
	}
	ref := []firestore.Update{{Path: "Photo", Value: url}}
	ctx := context.Background()
	if _, err = user.Ref.Update(ctx, ref); err != nil {
		log.Println("users -> (user *User) ChangesPhoto:2 -> err:", err)
		return "", err
	}
	return url, nil
}
