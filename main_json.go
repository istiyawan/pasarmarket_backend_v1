package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "salam"
	dbname   = "db_pasarmarket"
)

type Users struct {
	ID          int
	Username    string
	Useremail   string
	Password    string
	Createddate string
	Updateddate string
}

func Connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func HandleAddUsers(w http.ResponseWriter, r *http.Request) {

	db, err := Connect()

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		payload := struct {
			ID          int    `json:"id"`
			Username    string `json:"username"`
			Useremail   string `json:"useremail"`
			Password    string `json:"password"`
			Createddate string `json:"createddate"`
			Updateddate string `json:"updateddate"`
		}{}

		if err := decoder.Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("insert into tm_users values($1,$2,$3,$4,$5,$6)", payload.ID, payload.Username, payload.Useremail, payload.Password, payload.Createddate, payload.Updateddate)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if err != nil {
			log.Print(err)
		}

	}
	log.Print("Insert data to database")

}

func main() {
	http.HandleFunc("/AddUsers", HandleAddUsers)
	fmt.Println("Server mulai di port 8080")
	http.ListenAndServe(":8080", nil)
}
