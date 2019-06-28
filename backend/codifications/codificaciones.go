package codifications

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const fullAleatorio = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

//GetHash codify password
func GetHash(rawPass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

//RandomNumbersLetters returns random combination of letters and numbers
func RandomNumbersLetters(tamano int) string {
	b := make([]byte, tamano)
	for i := range b {
		b[i] = fullAleatorio[rand.Intn(len(fullAleatorio))]
	}
	return string(b)
}
