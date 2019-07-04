# About

- This is just a single chat project, it does not have automatically updates, if I someone would want to know how to do it, just create an issue and I will add it using [SSE](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events) :D.
- If you want to send emails, you should [by a domain](https://domains.google) and configure with [mailgun](https://www.mailgun.com), insert your crendentials into clients/mailgun and uncomment whenever there is a `sendEmail` function
  - Setting mailgun in your domain: [https://youtu.be/r-Qj4UWM3oM](https://youtu.be/r-Qj4UWM3oM)

# Usage
 - Create a project at [Google CLoud](https://console.cloud.google.com) 
 - Paste the project ID at clients/init.go -> variable: `ProjectID`
 - Login with your Google Account in [gcloud](https://cloud.google.com/sdk/gcloud/)  -> `gcloud init`
 - Run `gcloud beta auth application-default login`
 - Build and run code at `frontend/save`
 - Build and run code at `backend`
 - Go to [http://localhost:8070](http://localhost:8070)
# Tech
  - **Backend:**  [Go](https://golang.org)
  - **Frontend:**  HTML, CSS, Javascript
  - **Database:** [Firestore](https://firebase.google.com/docs/firestore)
  - **Static files on:** [Google Cloud Storage](https://cloud.google.com/storage)
  -  **Sending emails:** [mailgun](https://www.mailgun.com)
  - **Easy method to save static files:** [iguana](https://github.com/johnbalvin/iguana)