package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
func ProcessAPI(body map[string]interface{}, payload map[string]interface{}) string {
	command := body["command"].(string)
	switch command {
	case "login":
		return Login(body, payload)
	case "checklogin":
		return CheckLogin(body, payload)
	case "workdaycheck":
		return WorkDayCheck(body, payload)
	case "tangcadaycheck":
		return OverTimeDayCheck(body, payload)
	case "countxacnhanchamcong":
		return CheckinConfirm(body, payload)
	case "countthuongphat":
		return CountThuongPhat(body, payload)
	case "checkWebVer":
		return CheckWebVersion(body, payload)
	case "checkMYCHAMCONG":
		return CheckMyChamCong(body, payload)
	case "checkLicense":
		return CheckLicense(body, payload)
	case "diemdanhnhom":
		return DiemDanhNhom(body, payload)
	default:
		newjson := map[string]interface{}{
			"data":      nil,
			"tk_status": "NG",
			"message":   "Command not found",
		}
		jsonData, err := json.Marshal(newjson)
		if err != nil {
			log.Fatal("Error marshalling JSON:", err.Error())
		}
		return string(jsonData)
	}
}
