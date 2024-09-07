package login

import (
  cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider" 
  "github.com/aws/aws-sdk-go/aws"
	"github.com/7empestx/GoHTMXToDoList/internal/views/login"
	"github.com/7empestx/GoHTMXToDoList/internal/app"

  "crypto/hmac"
  "crypto/sha256"
  "encoding/base64"
  "fmt"
  "net/http"
  "context"
)

var appInstance *app.App

// InitApp initializes the App struct
func InitApp(a *app.App) {
    appInstance = a
}

// Login handles the login page
func Login(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" {
    loginGet(w, r)
  } else if r.Method == "POST" {
    loginPost(w, r);
  } else {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
  }
}

func loginGet(w http.ResponseWriter, r *http.Request) {
  component := login.LoginView()
  err := component.Render(r.Context(), w)
  if err != nil {
    fmt.Printf("Error rendering login component: %v", err)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
  }
}

func computeSecretHash(clientSecret string, username string, clientId string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func loginPost(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Login POST handler called")

	if err := r.ParseForm(); err != nil {
		fmt.Println("Error parsing form:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
  password := r.FormValue("password")

	if email == "" {
    http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
  if password == "" {
    http.Error(w, "Password is required", http.StatusBadRequest)
    return
  }

  fmt.Println(email);
  fmt.Println(password);		

  secretHash := computeSecretHash(appInstance.AppClientSecret, 
  *aws.String(email), appInstance.AppClientID)

  authTry := &cognito.InitiateAuthInput{
    AuthFlow: aws.String("USER_PASSWORD_AUTH"),
    AuthParameters: map[string]*string{
      "USERNAME": aws.String(email),
      "PASSWORD": aws.String(password),
      "SECRET_HASH": aws.String(secretHash), 
    },
    ClientId: aws.String(appInstance.AppClientID),
  }

  _, err := appInstance.CognitoClient.InitiateAuth(authTry)

	if err != nil {
    if(err.Error() == "NotAuthorizedException: Incorrect username or password.") {
      fmt.Println("Incorrect username or password")
      component := login.IncorrectLogin()
      component.Render(context.Background(), w)
    }
		return
	}

  component := login.SuccessfulLogin()
  component.Render(context.Background(), w)

}

