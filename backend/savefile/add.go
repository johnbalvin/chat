package savefile

import (
	"chat/backend/errors"
	"chat/clients"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"log"
	"mime/multipart"
	"regexp"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/googleapi"
)

var regURL *regexp.Regexp

func init() {
	var err error
	regURL, err = regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(" err: ", err)
	}
}

//AddImg adds image to cloud storage
func AddImg(nameStart string, file multipart.File, header *multipart.FileHeader) (bool, []byte, string, string, error) {
	return addFile("img", nameStart, file, header)
}

func addFile(kind, nameStart string, file multipart.File, header *multipart.FileHeader) (bool, []byte, string, string, error) {
	mimeType := header.Header.Get("Content-Type")
	switch kind {
	case "img":
		if !strings.Contains(mimeType, "image") {
			log.Println("savefile -> addFile:1 -> not supported: ", mimeType)
			return false, nil, "", "", errors.NotAllow
		}
	case "video":
		if !strings.Contains(mimeType, "video") {
			log.Println("savefile -> addFile:2 -> not supported: ", mimeType)
			return false, nil, "", "", errors.NotAllow
		}
	default:
		log.Println("savefile -> addFile:3 -> not supported: ", mimeType)
		return false, nil, "", "", errors.NotAllow
	}
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("savefile -> addFile:4 -> err: ", err)
		return false, nil, "", "", err
	}
	hasher := sha256.New()
	hasher.Write(fileByte)
	checkSum := regURL.ReplaceAllString(base64.RawURLEncoding.EncodeToString(hasher.Sum(nil)), "")
	name := nameStart + checkSum
	if err != nil {
		log.Println("savefile -> addFile:5 -> err: ", err)
		return false, nil, "", "", err
	}
	ctx := context.Background()
	w := clienteCS.Object(name).If(storage.Conditions{DoesNotExist: true}).NewWriter(ctx)
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.ContentType = mimeType
	w.CacheControl = "must-revalidate,max-age=31536000"
	w.ProgressFunc = func(copiedBytes int64) {
		log.Printf("copy %.1f%% done", float64(copiedBytes)/float64(header.Size)*100)
	}
	if _, err = w.Write(fileByte); err != nil {
		log.Println("savefile -> addFile:6 -> err: ", err)
		return false, nil, "", "", err
	}
	url := clients.BucketURL() + name
	if err = w.Close(); err != nil {
		if e, ok := err.(*googleapi.Error); ok {
			if e.Code != 412 {
				log.Println("savefile -> addFile:7 -> err: ", err)
				return false, nil, "", "", err
			}
		}
		log.Println("main -> repetive file: ", header.Filename)
		return true, nil, name, url, nil
	}
	return false, fileByte, name, url, nil
}
