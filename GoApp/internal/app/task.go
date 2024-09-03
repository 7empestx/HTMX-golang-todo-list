package app 

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/7empestx/GoHTMXToDoList/internal/store"
	"github.com/7empestx/GoHTMXToDoList/internal/views"
	"github.com/7empestx/GoHTMXToDoList/internal/login"
  "github.com/7empestx/GoHTMXToDoList/internal/store/sqlc"
	"github.com/gorilla/mux"
  "github.com/aws/aws-sdk-go/aws"
  cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider" 

  "crypto/hmac"
  "crypto/sha256"
  "encoding/base64"
)

var app *App

// InitApp initializes the App struct
func InitApp(a *App) {
    app = a
}

func Home(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Home handler called")
  component := login.Home()
  err := component.Render(r.Context(), w)
  if err != nil {
    fmt.Printf("Error rendering home component: %v", err)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
  }
}


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
  component := login.Login()
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

  secretHash := computeSecretHash(app.AppClientSecret, *aws.String(email), app.AppClientID)

  authTry := &cognito.InitiateAuthInput{
    AuthFlow: aws.String("USER_PASSWORD_AUTH"),
    AuthParameters: map[string]*string{
      "USERNAME": aws.String(email),
      "PASSWORD": aws.String(password),
      "SECRET_HASH": aws.String(secretHash), 
    },
    ClientId: aws.String(app.AppClientID),
  }


  _, err := app.CognitoClient.InitiateAuth(authTry)

	if err != nil {
    //http.Error(w, "Authentication failed", http.StatusUnauthorized)
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

func renderComponent(w http.ResponseWriter, tasks []storedb.Task) {
	component := views.Tasks(tasks)
	component.Render(context.Background(), w)
}

func FilterIncompleteTasks(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tasks, err := store.FilterIncompleteTasks(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderComponent(w, tasks)
}

func FilterCompletedTasks(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tasks, err := store.FilterCompletedTasks(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderComponent(w, tasks)
}

func Checked(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("Error parsing form:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("taskID")
	if idStr == "" {
		fmt.Println("No taskID received")
		http.Error(w, "No taskID received", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error converting taskID to integer:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err = store.Checked(ctx, int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	GetTasks(w, r)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tasks, err := store.GetTasks(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderComponent(w, tasks)
}

func logIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("Error parsing form:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ip := logIPAddress(r)
	description := r.FormValue("description")

	ctx := context.Background()
	err := store.AddTask(ctx, description, ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	GetTasks(w, r)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err = store.DeleteTask(ctx, int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	GetTasks(w, r)
}
