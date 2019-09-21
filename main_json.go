package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

func HashAndSalt(password []byte) string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
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
			Password    []byte `json:"password"`
			Createddate string `json:"createddate"`
			Updateddate string `json:"updateddate"`
		}{}

		if err := decoder.Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		hashPassword := HashAndSalt(payload.Password)

		_, err = db.Exec("insert into tm_users values($1,$2,$3,$4,$5,$6)", payload.ID, payload.Username, payload.Useremail, hashPassword, payload.Createddate, payload.Updateddate)
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

func HandleListUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	db, err := Connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query("select id, user_name, user_email,created_date, password, updated_date from tm_users")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	var result []Users
	for rows.Next() {
		var each = Users{}
		var err = rows.Scan(&each.ID, &each.Username, &each.Useremail, &each.Password, &each.Createddate, &each.Updateddate)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)
	}

	if r.Method == "GET" {
		var result, err = json.Marshal(result)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		w.Write(result)
		return
	}

	http.Error(w, "", http.StatusBadRequest)

}

func HandleListUserById(w http.ResponseWriter, r *http.Request) {
	db, err := Connect()

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	if r.Method == "GET" {

		decoder := json.NewDecoder(r.Body)
		payload := struct {
			ID int `json:"id"`
		}{}

		if err := decoder.Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var result []Users
		rows, err := db.Query("select id, user_name, user_email,password, created_date, updated_date from tm_users where id = $1", payload.ID)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for rows.Next() {
			var each = Users{}
			var err = rows.Scan(&each.ID, &each.Username, &each.Useremail, &each.Password, &each.Createddate, &each.Updateddate)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			result = append(result, each)
		}

		if r.Method == "GET" {
			var result, err = json.Marshal(result)

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			w.Write(result)
			return
		}

		http.Error(w, "", http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc("/AddUsers", HandleAddUsers)
	http.HandleFunc("/ListUsers", HandleListUsers)
	http.HandleFunc("/ListUsersById", HandleListUserById)
	fmt.Println("Server mulai di port 8080")
	http.ListenAndServe(":8080", nil)
}
