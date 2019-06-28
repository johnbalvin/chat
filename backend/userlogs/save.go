package userlogs

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"chat/backend/sessions"
	"net"
	"net/http"
	"strings"
	"time"
	"github.com/avct/uasurfer"
)

func getNormal(r *http.Request) Request {
	tiempo := time.Now().Unix()
	var objeto Request
	trackerID, userID, ID, err := sessions.ReadCookie(r)
	if err != nil { //error in decoding
		trackerID = "error"
	}

	objeto.TrackerID = trackerID
	objeto.User.ID = userID
	objeto.User.SessionID = ID
	objeto.Timestamp = tiempo
	objeto.RemoteAddr = r.RemoteAddr
	objeto.Host = r.Host
	objeto.Referer = r.Referer()
	objeto.Method = r.Method
	objeto.ContentLength = r.ContentLength
	objeto.ForwardedFor = r.Header.Get("X-Forwarded-For")
	var headers string
	jsonString, err := json.Marshal(r.Header)
	if err == nil {
		headers = string(jsonString)
	} else {
		headers = "error"
	}
	ip, port, err := net.SplitHostPort(r.RemoteAddr)
	objeto.IP = ip
	objeto.Port = port
	objeto.Headers = headers
	objeto.URL = r.RequestURI
	ua := uasurfer.Parse(r.UserAgent())
	objeto.UserAgent.Browser.Name = ua.Browser.Name.String()
	objeto.UserAgent.Browser.Version.Major = ua.Browser.Version.Major
	objeto.UserAgent.Browser.Version.Minor = ua.Browser.Version.Minor
	objeto.UserAgent.Browser.Version.Patch = ua.Browser.Version.Patch
	objeto.UserAgent.OS.Platform = ua.OS.Platform.String()
	objeto.UserAgent.OS.Name = ua.OS.Name.String()
	objeto.UserAgent.OS.Version.Major = ua.OS.Version.Major
	objeto.UserAgent.OS.Version.Minor = ua.OS.Version.Minor
	objeto.UserAgent.OS.Version.Patch = ua.OS.Version.Patch
	objeto.UserAgent.DeviceType = ua.DeviceType.String()

	return objeto
}

//SaveNormal saves all request made by user
func SaveNormal(r *http.Request) {
	objeto := getNormal(r)
	u := tableNormal.Inserter()
	ctx := context.Background()
	if err := u.Put(ctx, objeto); err != nil {
		log.Println("userlogs -> SaveNormal:1 -> err:", err)
		return
	}
}
//SaveBloquer saves all request made by user
func SaveBloquer(r *http.Request) {
	tiempo := time.Now().Unix()
	var objeto Preference
	objeto.Timestamp = tiempo
	trackerID, _, _, err := sessions.ReadCookie(r)
	if err != nil { //error in decoding
		trackerID = "error"
	}
	objeto.TrackerID = trackerID
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("userlogs -> SaveBloquer:1 -> err:", err)
		return
	}
	bodyString := string(bodyBytes)
	slice := strings.Split(bodyString, "")
	if len(slice) == 3 {
		objeto.Kind = slice[0]
		if slice[1] == "y" { //Darkmode
			objeto.DarkMode = true
		}
		if slice[2] == "y" { //bloquer
			objeto.Bloquer = true
		}
	} else {
		objeto.MaliciousUser = true
	}

	objeto.Raw = bodyString
	u := tableBloquer.Inserter()
	ctx := context.Background()
	if err := u.Put(ctx, objeto); err != nil {
		log.Println("userlogs -> SaveBloquer:2 -> err:", err)
		return
	}
}
