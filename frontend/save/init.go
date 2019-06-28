package main

import (
	"chat/clients"
	"time"

	"github.com/johnbalvin/iguana"
)

var clienteFS = clients.Firestore()
var clienteCS = clients.Storage()
var docDirectory = clienteFS.Collection("frontEnd").Doc("Directory")
var docServiceWorkers = clienteFS.Collection("frontEnd").Doc("SW")

const constStartInPath = "../../frontEnd"

func idURLNormal(static iguana.Static) (string, string) {
	id := "s/" + static.Content.Checksum + "." + static.Extension
	url := clients.BucketURL() + id
	return id, url
}
func idURLObf(static iguana.Static) (string, string) {
	id := "s/" + static.ContentObf.Checksum + "." + static.Extension
	url := clients.BucketURL() + id
	return id, url
}

type html struct {
	Checksum     string
	DataGenerate bool
	LastModify   time.Time
}

type swInfo struct {
	Checksum string
	ID       string
}
