package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

func ExcuteQuery(querystring string) string {
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
		//log.Fatal("Open connection failed:", err.Error())

		//return combined json data
		result := map[string]interface{}{
			"data":      nil,
			"tk_status": "NG",
			"message":   err.Error(),
		}
		jsonData, err := json.Marshal(result)
		if err != nil {
			log.Fatal("Error marshalling JSON:", err.Error())
		}
		return string(jsonData)
	}
	defer db.Close()

	// Kiểm tra kết nối
	err = db.Ping()
	if err != nil {
		//log.Fatal("Ping failed:", err.Error())
		//return combined json data
		result := map[string]interface{}{
			"data":      nil,
			"tk_status": "NG",
			"message":   err.Error(),
		}
		jsonData, err := json.Marshal(result)
		if err != nil {
			log.Fatal("Error marshalling JSON:", err.Error())
		}
		return string(jsonData)
	}
	// Truy vấn dữ liệu

	rows, err := db.Query(querystring)
	if err != nil {
		//log.Fatal("Query failed:", err.Error())
		//return combined json data
		result := map[string]interface{}{
			"data":      nil,
			"tk_status": "NG",
			"message":   err.Error(),
		}
		jsonData, err := json.Marshal(result)
		if err != nil {
			log.Fatal("Error marshalling JSON:", err.Error())
		}
		return string(jsonData)
	}
	defer rows.Close()

	// Create a slice to store the results
	var results []map[string]interface{}

	// Iterate through the rows
	columns, _ := rows.Columns()
	for rows.Next() {
		// Create a slice of interface{} to hold each row's values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into the slice of interface{}
		if err := rows.Scan(valuePtrs...); err != nil {
			//log.Fatal("Error scanning row:", err.Error())
			//return combined json data
			result := map[string]interface{}{
				"data":      nil,
				"tk_status": "NG",
				"message":   err.Error(),
			}
			jsonData, err := json.Marshal(result)
			if err != nil {
				log.Fatal("Error marshalling JSON:", err.Error())
			}
			return string(jsonData)

		}

		// Create a map for the current row
		row := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			row[col] = v
		}

		// Append the row to the results
		results = append(results, row)
	}

	// Kiểm tra lỗi sau khi duyệt xong các dòng
	if err = rows.Err(); err != nil {
		//log.Fatal("Rows iteration error:", err.Error())
		//return combined json data
		result := map[string]interface{}{
			"data":      nil,
			"tk_status": "NG",
			"message":   err.Error(),
		}
		jsonData, err := json.Marshal(result)
		if err != nil {
			log.Fatal("Error marshalling JSON:", err.Error())
		}
		return string(jsonData)
	}

	//combine results, tk_status, message
	result := map[string]interface{}{
		"data":      results,
		"tk_status": "OK",
		"message":   "Success",
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error marshalling JSON:", err.Error())
	}
	return string(jsonData)
}
func ProcessAPI(c echo.Context) string {
	requestBody := c.Request().Body
	body, err := io.ReadAll(requestBody)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		return "Error reading request body"
	}

	parsedBody := make(map[string]interface{})
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		return "Error unmarshalling JSON"
	}
	/* fmt.Println(string(body))
	fmt.Println(parsedBody) */
	fmt.Println(parsedBody["command"])
	DATA := parsedBody["DATA"].(map[string]interface{})
	fmt.Println(DATA["command1"])
	fmt.Println(DATA["command2"])
	result := ExcuteQuery("SELECT TOP 100 * FROM AMAZONE_DATA")
	return result
}
