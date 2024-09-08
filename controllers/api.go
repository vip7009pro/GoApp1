package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
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
		user := body["user"].(string)
		pass := body["pass"].(string)
		result := ExcuteQuery("SELECT ZTBEMPLINFO.EMPL_IMAGE,ZTBEMPLINFO.CTR_CD,ZTBEMPLINFO.EMPL_NO,ZTBEMPLINFO.CMS_ID,ZTBEMPLINFO.FIRST_NAME,ZTBEMPLINFO.MIDLAST_NAME,ZTBEMPLINFO.DOB,ZTBEMPLINFO.HOMETOWN,ZTBEMPLINFO.SEX_CODE,ZTBEMPLINFO.ADD_PROVINCE,ZTBEMPLINFO.ADD_DISTRICT,ZTBEMPLINFO.ADD_COMMUNE,ZTBEMPLINFO.ADD_VILLAGE,ZTBEMPLINFO.PHONE_NUMBER,ZTBEMPLINFO.WORK_START_DATE,ZTBEMPLINFO.PASSWORD,ZTBEMPLINFO.EMAIL,ZTBEMPLINFO.WORK_POSITION_CODE,ZTBEMPLINFO.WORK_SHIFT_CODE,ZTBEMPLINFO.POSITION_CODE,ZTBEMPLINFO.JOB_CODE,ZTBEMPLINFO.FACTORY_CODE,ZTBEMPLINFO.WORK_STATUS_CODE,ZTBEMPLINFO.REMARK,ZTBEMPLINFO.ONLINE_DATETIME,ZTBSEX.SEX_NAME,ZTBSEX.SEX_NAME_KR,ZTBWORKSTATUS.WORK_STATUS_NAME,ZTBWORKSTATUS.WORK_STATUS_NAME_KR,ZTBFACTORY.FACTORY_NAME,ZTBFACTORY.FACTORY_NAME_KR,ZTBJOB.JOB_NAME,ZTBJOB.JOB_NAME_KR,ZTBPOSITION.POSITION_NAME,ZTBPOSITION.POSITION_NAME_KR,ZTBWORKSHIFT.WORK_SHIF_NAME,ZTBWORKSHIFT.WORK_SHIF_NAME_KR,ZTBWORKPOSITION.SUBDEPTCODE,ZTBWORKPOSITION.WORK_POSITION_NAME,ZTBWORKPOSITION.WORK_POSITION_NAME_KR,ZTBWORKPOSITION.ATT_GROUP_CODE,ZTBSUBDEPARTMENT.MAINDEPTCODE,ZTBSUBDEPARTMENT.SUBDEPTNAME,ZTBSUBDEPARTMENT.SUBDEPTNAME_KR,ZTBMAINDEPARMENT.MAINDEPTNAME,ZTBMAINDEPARMENT.MAINDEPTNAME_KR FROM ZTBEMPLINFO LEFT JOIN ZTBSEX ON (ZTBSEX.SEX_CODE = ZTBEMPLINFO.SEX_CODE) LEFT JOIN ZTBWORKSTATUS ON(ZTBWORKSTATUS.WORK_STATUS_CODE = ZTBEMPLINFO.WORK_STATUS_CODE) LEFT JOIN ZTBFACTORY ON (ZTBFACTORY.FACTORY_CODE = ZTBEMPLINFO.FACTORY_CODE) LEFT JOIN ZTBJOB ON (ZTBJOB.JOB_CODE = ZTBEMPLINFO.JOB_CODE) LEFT JOIN ZTBPOSITION ON (ZTBPOSITION.POSITION_CODE = ZTBEMPLINFO.POSITION_CODE) LEFT JOIN ZTBWORKSHIFT ON (ZTBWORKSHIFT.WORK_SHIFT_CODE = ZTBEMPLINFO.WORK_SHIFT_CODE) LEFT JOIN ZTBWORKPOSITION ON (ZTBWORKPOSITION.WORK_POSITION_CODE = ZTBEMPLINFO.WORK_POSITION_CODE) LEFT JOIN ZTBSUBDEPARTMENT ON (ZTBSUBDEPARTMENT.SUBDEPTCODE = ZTBWORKPOSITION.SUBDEPTCODE) LEFT JOIN ZTBMAINDEPARMENT ON (ZTBMAINDEPARMENT.MAINDEPTCODE = ZTBSUBDEPARTMENT.MAINDEPTCODE) WHERE ZTBEMPLINFO.EMPL_NO = '" + user + "' AND PASSWORD = '" + pass + "'")
		//convert result to map
		var resultMap map[string]interface{}
		err := json.Unmarshal([]byte(result), &resultMap)
		if err != nil {
			log.Fatal("Error unmarshalling JSON:", err.Error())
		}
		//get token from resultMap
		data, ok := resultMap["data"].([]interface{})
		if !ok || len(data) == 0 {
			log.Fatal("Invalid data format in resultMap")
		}
		//fmt.Println(data)
		loginResult, ok := data[0].(map[string]interface{})
		if !ok {
			log.Fatal("Invalid login result format")
		}
		// Set expiration time to 5 minutes from now
		expirationTime := time.Now().Add(5 * time.Minute)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"payload": loginResult,
			"exp":     expirationTime.Unix(), // Add expiration claim
		})
		//new json
		newJson := map[string]interface{}{
			"tk_status": "OK",
			"userData":  loginResult,
		}
		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			log.Fatal("Error signing token:", err.Error())
		}
		newJson["token_content"] = tokenString
		resultJson, err := json.Marshal(newJson)
		if err != nil {
			log.Fatal("Error marshalling JSON:", err.Error())
		}
		return string(resultJson)

	case "checklogin":
		query := "SELECT WORK_STATUS_CODE FROM ZTBEMPLINFO WHERE EMPL_NO='" + payload["EMPL_NO"].(string) + "'"
		result := ExcuteQuery(query)
		resultmap := make(map[string]interface{})
		err := json.Unmarshal([]byte(result), &resultmap)
		if err != nil {
			log.Fatal("Error unmarshalling JSON:", err.Error())
		}
		data, ok := resultmap["data"].([]interface{})
		if !ok || len(data) == 0 {
			log.Fatal("Invalid data format in resultmap")
		}
		workStatusCode, ok := data[0].(map[string]interface{})["WORK_STATUS_CODE"]
		if !ok {
			log.Fatal("Invalid work status code format")
		}
		if workStatusCode != 0 {
			newjson := map[string]interface{}{
				"tk_status": "OK",
				"message":   "Success",
				"data":      payload,
			}
			resultJson, err := json.Marshal(newjson)
			if err != nil {
				log.Fatal("Error marshalling JSON:", err.Error())
			}
			return string(resultJson)
		} else {
			newjson := map[string]interface{}{
				"tk_status": "NG",
				"message":   "Đã nghỉ việc",
			}
			resultJson, err := json.Marshal(newjson)
			if err != nil {
				log.Fatal("Error marshalling JSON:", err.Error())
			}
			return string(resultJson)
		}

	}
	return "error"
}
