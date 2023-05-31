package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Response struct {
	Code    int         `json:"code"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func connectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/go_crud_db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	router := mux.NewRouter()

	// Get
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		db := connectDB()
		defer db.Close()

		rows, err := db.Query("SELECT * FROM tb_user")
		if err != nil {
			log.Fatal(err)
		}

		var users []User
		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Username, &user.Email)
			if err != nil {
				log.Fatal(err)
			}
			users = append(users, user)
		}

		response := Response{
			Code:    http.StatusOK,
			Success: true,
			Message: "Users retrieved",
			Data:    users,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("GET")

	// Get User ID
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		db := connectDB()
		defer db.Close()

		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			log.Fatal(err)
		}

		var user User
		err = db.QueryRow("SELECT * FROM tb_user WHERE id = ?", id).Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				response := Response{
					Code:    http.StatusNotFound,
					Success: false,
					Message: "User not found",
					Data:    nil,
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}
			log.Fatal(err)
		}

		response := Response{
			Code:    http.StatusOK,
			Success: true,
			Message: "User retrieved",
			Data:    user,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("GET")

	// Create User
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		db := connectDB()
		defer db.Close()

		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Fatal(err)
		}

		result, err := db.Exec("INSERT INTO tb_user (username, email) VALUES (?, ?)", user.Username, user.Email)
		if err != nil {
			log.Fatal(err)
		}

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}

		response := Response{
			Code:    http.StatusOK,
			Success: true,
			Message: "User created",
			Data:    lastInsertID,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("POST")

	// Update User
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		db := connectDB()
		defer db.Close()

		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			log.Fatal(err)
		}

		var user User
		err = json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("UPDATE tb_user SET username = ?, email = ? WHERE id = ?", user.Username, user.Email, id)
		if err != nil {
			log.Fatal(err)
		}

		response := Response{
			Code:    http.StatusOK,
			Success: true,
			Message: "User updated",
			Data:    nil,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("PUT")

	// Delete User
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		db := connectDB()
		defer db.Close()

		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("DELETE FROM tb_user WHERE id = ?", id)
		if err != nil {
			log.Fatal(err)
		}

		response := Response{
			Code:    http.StatusOK,
			Success: true,
			Message: "User deleted",
			Data:    nil,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
