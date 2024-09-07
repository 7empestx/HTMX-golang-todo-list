package auth

import (
  "os"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
  "github.com/7empestx/GoHTMXToDoList/internal/app"
  "github.com/7empestx/GoHTMXToDoList/internal/router/login"
)

func Init() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), 
	}))
	cognitoClient := cognito.New(sess)

	userPoolID := os.Getenv("COGNITO_USER_POOL_ID")
	appClientID := os.Getenv("COGNITO_APP_CLIENT_ID")
	appClientSecret := os.Getenv("COGNITO_APP_CLIENT_SECRET")

	// Create App instance
	appInstance := &app.App {
		CognitoClient:   cognitoClient,
		UserPoolID:      userPoolID,
		AppClientID:     appClientID,
		AppClientSecret: appClientSecret,
	}

  login.InitApp(appInstance)
}

