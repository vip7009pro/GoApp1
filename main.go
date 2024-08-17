package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	router := mux.NewRouter()

	// CORS settings
	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	router.HandleFunc("/check", handleURLEncodedForm).Methods("GET")
	router.HandleFunc("/api", YourAPIFunction).Methods("POST")
	router.HandleFunc("/uploadfile", UploadFileFunction).Methods("POST")

	API_PORT := os.Getenv("API_PORT")
	log.Printf("Server starting at port %s", API_PORT)
	log.Fatal(http.ListenAndServe(":"+API_PORT, corsOptions(router)))
}

func handleURLEncodedForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key1 := r.FormValue("hung")
	key2 := r.FormValue("vanhung")
	fmt.Fprintf(w, "Received Key1: %s, Key2: %s\n", key1, key2)
}

func YourAPIFunction(w http.ResponseWriter, r *http.Request) {
	// Xử lý request API
	w.Write([]byte("API Endpoint"))
	fmt.Println(r)
}

func UploadFileFunction(w http.ResponseWriter, r *http.Request) {
	// Xử lý upload file
	w.Write([]byte("Upload File Endpoint"))
}
