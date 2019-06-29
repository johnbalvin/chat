package clients

import "log"

//ProjectID is the proyect ID
var ProjectID = "" //inset your project ID
var bucketName, url string

//dev feature production
func init() {
	if ProjectID == "" {
		log.Fatalf("You haven't assing a projectID, create one at https://console.cloud.google.com")
	}
	bucketName = ProjectID + ".appspot.com"
	url = "https://storage.googleapis.com/" + bucketName + "/"
}
