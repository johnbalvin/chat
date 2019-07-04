package clients

import (
	"log"
	"os"
)

//ProjectID is the proyect ID
var ProjectID = "" //inset your project ID
var bucketName, url string

//dev feature production
func init() {
	if ProjectID == "" {
		log.Fatalf("You haven't assing a projectID, create one at https://console.cloud.google.com")
	}
	projectID := os.Getenv("projectID")
	if projectID != "" {
		ProjectID = projectID
	}
	bucketName = ProjectID + ".appspot.com"
	url = "https://storage.googleapis.com/" + bucketName + "/"
}
