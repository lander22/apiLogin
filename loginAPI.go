package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "gopkg.in/yaml.v2"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//DB object
var db *sql.DB

func initDB() {

	cfg := mysql.Config{
		//User:   os.Getenv("DBUSER"),
		//Passwd: os.Getenv("DBPASS"),
		User:   "lander",
		Passwd: "password",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "apiLogin",
	}

	var err error

	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		panic(err.Error())
	}
}

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
	Passwd  string `json:"password`
}

var users []User

func createUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint Hit: Create User")

	decoder := json.NewDecoder(r.Body)
	var data User
	decoder.Decode(&data)

	if len(data.Name) > 0 && len(data.Surname) > 0 && len(data.Email) > 0 && len(data.Passwd) > 0 {

		sql := fmt.Sprintf("INSERT INTO users (name,lastname,email,password) VALUES ('%s','%s','%s',SHA1('%s'));", data.Name, data.Surname, data.Email, data.Passwd)

		res, err := db.Exec(sql)

		if err != nil {
			panic(err)
		}

		lastId, err := res.LastInsertId()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("The last inserted row id: %d\n", lastId)

		w.WriteHeader(201)

	} else {
		fmt.Println("Malformed Data Json. Needed: Name, Surname, Email and Passwd")
		w.WriteHeader(400)
		return
	}
}

func checkUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Check User")

	decoder := json.NewDecoder(r.Body)
	var data User
	decoder.Decode(&data)

	if len(data.Email) > 0 {

		sql := fmt.Sprintf("SELECT email FROM users WHERE email = '%s'", data.Email)

		res, err := db.Query(sql)

		if err != nil {
			panic(err)
		}

		if res.Next() {
			fmt.Println("El usuario ya existe")
			w.WriteHeader(200)
		} else {
			fmt.Println("El usuario no existe")
			w.WriteHeader(204) //no content
		}

	} else {
		fmt.Println("Malformed Data Json. Needed: Name")
		w.WriteHeader(400)
		return
	}
}

func checkUserCredentials(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Check User Credentials")

	decoder := json.NewDecoder(r.Body)
	var data User
	decoder.Decode(&data)

	if len(data.Email) > 0 && len(data.Passwd) > 0 {

		sql := fmt.Sprintf("SELECT email FROM users WHERE email = '%s' and password = SHA1('%s')", data.Email, data.Passwd)

		res, err := db.Query(sql)

		if err != nil {
			panic(err)
		}

		if res.Next() {
			fmt.Println("Credenciales correctas")
			w.WriteHeader(202)
		} else {
			fmt.Println("Credenciales incorrectas")
			w.WriteHeader(401)
		}

	} else {
		fmt.Println("Malformed Data Json. Needed: Email and Password")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Delete User")

	decoder := json.NewDecoder(r.Body)
	var data User
	decoder.Decode(&data)

	if len(data.Email) > 0 && len(data.Passwd) > 0 {

		sql := fmt.Sprintf("DELETE FROM users WHERE email = '%s' and password = SHA1('%s')", data.Email, data.Passwd)

		res, err := db.Exec(sql)

		if err != nil {
			panic(err)
		}

		var count, err2 = res.RowsAffected()

		if err2 != nil {
			panic(err2)
		}

		if count > 0 {
			fmt.Printf("User %s eliminated", data.Email)
			w.WriteHeader(202)
		} else {
			fmt.Println("No one was eliminated")
			w.WriteHeader(401)
		}

	} else {
		fmt.Println("Malformed Data Json. Needed: Email and Password")
		w.WriteHeader(400)
		return
	}
}

func main() {

	router := mux.NewRouter()

	//endpoints
	router.HandleFunc("/api/v1/createUser", createUser).Methods("POST")
	router.HandleFunc("/api/v1/checkUser", checkUser).Methods("GET")
	router.HandleFunc("/api/v1/checkUserCredentials", checkUserCredentials).Methods("GET")
	router.HandleFunc("/api/v1/deleteUser", deleteUser).Methods("DELETE")

	//db
	initDB()

	//log.Fatal(http.ListenAndServeTLS(":8965", "localhost.crt", "localhost.key", router))
	log.Fatal(http.ListenAndServe(":8965", router))
}
