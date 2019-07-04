package userlogs

import (
	"chat/clients"

	"cloud.google.com/go/bigquery"
)

var clienteBQ = clients.Bigquery()
var tableNormal, tableBloquer *bigquery.Table

type version struct {
	Major int
	Minor int
	Patch int
}

type user struct {
	ID            string
	SessionNumber string
	SessionID     string
}

type os struct {
	Platform string
	Name     string
	Version  version
}

type browser struct {
	Name    string
	Version version
}

type userAgent struct {
	DeviceType string
	Browser    browser
	OS         os
}

//Request is the info for a request
type Request struct {
	TrackerID     string
	User          user
	UserAgent     userAgent
	RemoteAddr    string
	ForwardedFor  string
	Referer       string
	URL           string
	Method        string
	IP            string
	Port          string
	Host          string
	Timestamp     int64
	Headers       string
	ContentLength int64
}

//Preference is the info for a request
type Preference struct {
	TrackerID     string
	Kind          string
	DarkMode      bool
	Bloquer       bool
	MaliciousUser bool
	Raw           string
	Timestamp     int64
}

func init() {
	tableNormal = clienteBQ.Dataset("user").Table("requests")
	tableBloquer = clienteBQ.Dataset("user").Table("preference")
}
