package main

import (
  "fmt"
	"log"
	"net/http"
  "os"

	"github.com/7empestx/GoHTMXToDoList/internal/app"
	"github.com/7empestx/GoHTMXToDoList/internal/store"
	"github.com/gorilla/mux"
  "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
   cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func main() {

	dbHost := os.Getenv("RDS_HOSTNAME")
	dbName := os.Getenv("RDS_DB_NAME")
	dbUser := os.Getenv("RDS_USERNAME")
	dbPassword := os.Getenv("RDS_PASSWORD")

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)
  fmt.Println(dataSourceName)
	store.InitDB(dataSourceName)

  // AWS Cognito initialization
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Replace with your AWS region
	}))
	cognitoClient := cognito.New(sess)

	userPoolID := os.Getenv("COGNITO_USER_POOL_ID")
	appClientID := os.Getenv("COGNITO_APP_CLIENT_ID")
	appClientSecret := os.Getenv("COGNITO_APP_CLIENT_SECRET")

	// Create App instance
	appInstance := &app.App{
		CognitoClient:   cognitoClient,
		UserPoolID:      userPoolID,
		AppClientID:     appClientID,
		AppClientSecret: appClientSecret,
	}

  app.InitApp(appInstance)

	r := mux.NewRouter()
  r.HandleFunc("/login", app.Login).Methods("GET")
  r.HandleFunc("/login", app.Login).Methods("POST")
  r.HandleFunc("/home", app.Home).Methods("GET")
	r.HandleFunc("/tasks", app.GetTasks).Methods("GET")
	r.HandleFunc("/tasks", app.AddTask).Methods("POST")
	r.HandleFunc("/completed", app.FilterCompletedTasks).Methods("GET")
	r.HandleFunc("/incomplete", app.FilterIncompleteTasks).Methods("GET")
	r.HandleFunc("/checked", app.Checked).Methods("POST")
	r.HandleFunc("/delete/{id}", app.DeleteTask).Methods("POST")

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/").Handler(fs)

  // Serve static files  
	log.Fatal(http.ListenAndServe(":8080", r))
}
