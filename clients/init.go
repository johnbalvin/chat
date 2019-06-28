package clients

import (
	"log"
	"os"
)

//ProjectID is the proyect ID
var ProjectID, bucketName, url string

//dev feature production
func init() {
	gitFlow := os.Getenv("GitFlow")
	log.Println("GitFlow: ", gitFlow)
	switch gitFlow {
	case "master":
		ProjectID = "" //Project ID, for this git branch
	case "dev":
		ProjectID = "" //Project ID, for this git branch
	case "qa":
		ProjectID = "" //Project ID, for this git branch
	case "feature":
		ProjectID = "" //Project ID, for this git branch
	default:
		ProjectID = "" //Project ID, for this git branch
	}
	bucketName = ProjectID + ".appspot.com"
	url = "https://storage.googleapis.com/" + bucketName + "/"

	log.Println("ProjectID: ", ProjectID)
}
