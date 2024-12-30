package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	_ "github.com/lib/pq"
)

type User struct {
	ID     int    `json:"id"`
	F_name string `json:"F_name"`
	L_name string `json:"L_name"`
}

func main() {
	// Connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create the table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS humans (id SERIAL PRIMARY KEY, F_name TEXT, L_name TEXT)")
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// Create router
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/humans", getUsers(db)).Methods("GET")
	router.HandleFunc("/humans/{id}", getUser(db)).Methods("GET")
	router.HandleFunc("/humans", createUser(db)).Methods("POST")
	router.HandleFunc("/humans/{id}", updateUser(db)).Methods("PUT")
	router.HandleFunc("/humans/{id}", deleteUser(db)).Methods("DELETE")

	// ใช้งาน CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // สามารถปรับ URL นี้ตามความต้องการ
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	// Start server
	log.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", corsHandler(router))) // ใช้ CORS handler
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// Root handler
func rootHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"Create": "POST: /humans",
		"ReadAll": "GET: /humans",
		"ReadOne": "GET: /humans/{id}",
		"Update": "PUT: /humans/{id}",
		"Delete": "DELETE: /humans/{id}",
	}
	json.NewEncoder(w).Encode(response)
}

// Get all users
func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, F_name, L_name FROM humans")
		if err != nil {
			http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.F_name, &u.L_name); err != nil {
				http.Error(w, "Error scanning user", http.StatusInternalServerError)
				return
			}
			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Error with rows", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(users)
	}
}

// Get user by ID
func getUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u User
		err := db.QueryRow("SELECT id, F_name, L_name FROM humans WHERE id = $1", id).Scan(&u.ID, &u.F_name, &u.L_name)
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(u)
	}
}

// Create user
func createUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := db.QueryRow("INSERT INTO humans (F_name, L_name) VALUES ($1, $2) RETURNING id", u.F_name, u.L_name).Scan(&u.ID)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(u)
	}
}

// Update user
func updateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		result, err := db.Exec("UPDATE humans SET F_name = $1, L_name = $2 WHERE id = $3", u.F_name, u.L_name, id)
		if err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(u)
	}
}

// Delete user
func deleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		result, err := db.Exec("DELETE FROM humans WHERE id = $1", id)
		if err != nil {
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "User deleted"})
	}
}
