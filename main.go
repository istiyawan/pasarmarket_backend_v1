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

	rows, err := db.Query("select id, user_name, user_email,created_date, updated_date from tm_users")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	var result []Users
	for rows.Next() {
		var each = Users{}
		var err = rows.Scan(&each.ID, &each.Username, &each.Useremail, &each.Createddate, &each.Updateddate)
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

func HandleAddUsers(w http.ResponseWriter, r *http.Request) {
	db, err := Connect()

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	err = r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}

	// fmt.Println(r.FormValue)

	Username := r.FormValue("username")
	Useremail := r.FormValue("useremail")
	Password := r.FormValue("password")
	ID := r.FormValue("id")
	Createddate := r.FormValue("createddate")
	Updateddate := r.FormValue("updateddate")

	_, err = db.Exec("insert into tm_users values($1,$2,$3,$4,$5,$6)", ID, Username, Useremail, Password, Createddate, Updateddate)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err != nil {
		log.Print(err)
	}

	log.Print("Insert data to database")
}

func HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	db, err := Connect()

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	if id := r.URL.Query().Get("id"); id != "" {
		var result []Users
		rows, err := db.Query("select id, user_name, user_email,created_date from tm_users where id = $1", id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for rows.Next() {
			var each = Users{}
			var err = rows.Scan(&each.ID, &each.Username, &each.Useremail, &each.Createddate)
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

func HandleDeleteUsers(w http.ResponseWriter, r *http.Request) {
	db, err := Connect()

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	if id := r.URL.Query().Get("id"); id != "" {

		_, err := db.Query("delete from tm_users where id = $1", id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		http.Error(w, "Deleted data success", http.StatusBadRequest)
	}

}

func HandleUpdateUsers(w http.ResponseWriter, r *http.Request) {
	db, err := Connect()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	email := r.URL.Query().Get("useremail")
	name := r.URL.Query().Get("username")
	if id := r.URL.Query().Get("id"); id != "" {
		_, err := db.Query("update tm_users set user_email=$1,user_name = $2 where id=$3", email, name, id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		http.Error(w, "Update data success", http.StatusBadRequest)
	}
}

func main() {

	http.HandleFunc("/ListUsers", HandleListUsers)
	http.HandleFunc("/AddUsers", HandleAddUsers)
	http.HandleFunc("/GetUsers", HandleGetUsers)
	http.HandleFunc("/DeleteUsers", HandleDeleteUsers)
	http.HandleFunc("/UpdateUsers", HandleUpdateUsers)

	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
