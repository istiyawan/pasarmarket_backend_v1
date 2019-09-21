package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Users struct {
	ID          int
	Username    string
	Useremail   string
	Password    string
	Createddate time.time
	Updateddate time.time
}

func HandleListUsers(w http.ResponseWriter, r *http.Request) {
	users := Users{}
	err := json.NewDecoder(r.Body).Decode(users)
	if err != nil {
		panic(err)
	}

	users.Createddate = time.Now().Local()

	userJson, err := json.Marshal(users)

}

func main() {
	http.HandleFunc("/ListUsers", HandleListUsers)
	fmt.Println("Server mulai di port 8080")
	http.ListenAndServe(":8080", nil)
}
