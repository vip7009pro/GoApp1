package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"

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
	ConnectSQL()
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

func ConnectSQL() {
	server := os.Getenv("DB_SERVER")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	database := os.Getenv("DB_NAME")

	connString := fmt.Sprintf("server=%s;port=%s;user id=%s;password=%s;database=%s",
		server, port, user, password, database)

	// Mở kết nối tới cơ sở dữ liệu
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer db.Close()

	// Kiểm tra kết nối
	err = db.Ping()
	if err != nil {
		log.Fatal("Ping failed:", err.Error())
	}

	fmt.Println("Connected to SQL Server successfully")

	// Truy vấn dữ liệu
	query := "SELECT TOP 100 EMPL_NO, EMPL_NAME FROM M010"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("Query failed:", err.Error())
	}
	defer rows.Close()

	// Lặp qua từng dòng kết quả và in ra console
	for rows.Next() {
		var column1 string
		var column2 string

		err := rows.Scan(&column1, &column2)
		if err != nil {
			log.Fatal("Scan failed:", err.Error())
		}

		fmt.Printf("EMPL_NO: %s, EMPL_NAME: %s\n", column1, column2)
	}

	// Kiểm tra lỗi sau khi duyệt xong các dòng
	if err = rows.Err(); err != nil {
		log.Fatal("Rows iteration error:", err.Error())
	}
}
