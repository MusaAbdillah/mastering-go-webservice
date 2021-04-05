package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type API struct {
	Message string "json:message"
}

type User struct {
	ID    int    "json:id"
	Name  string "json:username"
	Email string "json:email"
	First string "json:first"
	Last  string "json:last"
}

var (
	ctx context.Context
)

func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/social_network")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// connect to db first
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	urlParams := mux.Vars(r)
	id := urlParams["id"]
	ReadUser := User{}
	errQueryRow := db.
		QueryRow("SELECT user_id, user_nickname, user_email, user_first, user_last FROM users WHERE user_id = ?", id).
		Scan(&ReadUser.ID, &ReadUser.Name, &ReadUser.Email, &ReadUser.First, &ReadUser.Last)
	switch {
	case errQueryRow == sql.ErrNoRows:
		fmt.Fprintf(w, "No such user.")
	case errQueryRow != nil:
		log.Fatal(errQueryRow)
		fmt.Fprintf(w, "error")
	default:
		output, _ := json.Marshal(ReadUser)
		fmt.Fprintf(w, string(output))
	}

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// connect to db first
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	NewUser := User{}
	NewUser.Name = r.FormValue("name")
	NewUser.Email = r.FormValue("email")
	NewUser.First = r.FormValue("first")
	NewUser.Last = r.FormValue("last")

	output, err := json.Marshal(NewUser)
	fmt.Println(output)
	if err != nil {
		fmt.Println("Something went wrong!")
	}

	sql := "INSERT INTO users set user_nickname='" + NewUser.Name + "', user_first='" + NewUser.First + "', user_last='" + NewUser.Last + "', user_email='" + NewUser.Email + "' "
	q, err := db.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(q)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	name := urlParams["user"]
	HelloMessage := "Hello, " + name
	message := API{HelloMessage}
	output, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Something went wrong!")
	}
	fmt.Fprintf(w, string(output))
}

func main() {
	gorillaRoute := mux.NewRouter()
	gorillaRoute.HandleFunc("/api/{user:[0-9]+}", Hello)
	gorillaRoute.HandleFunc("/api/user/create", CreateUser).Methods("GET")
	gorillaRoute.HandleFunc("/api/user/get/{id}", GetUser).Methods("GET")
	http.Handle("/", gorillaRoute)
	http.ListenAndServe(":8080", nil)
}
