package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const CURRENT_API_URL = "https://script.google.com/macros/s/AKfycbyD_LRqVLETu8IvuiqDSsbItdmzRw3p_q9gCv12UOer0V-5OnqtbJvKjK86bfgGbUM1NA/exec"

func Login(body map[string]interface{}, payload map[string]interface{}) string {
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
}

func CheckLogin(body map[string]interface{}, payload map[string]interface{}) string {
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

func WorkDayCheck(body map[string]interface{}, payload map[string]interface{}) string {
	startOfYear := time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	endOfYear := time.Date(time.Now().Year(), time.December, 31, 23, 59, 59, 0, time.UTC)
	query := "SELECT COUNT(EMPL_NO) AS WORK_DAY FROM ZTBATTENDANCETB WHERE EMPL_NO='" + payload["EMPL_NO"].(string) + "' AND ON_OFF=1 AND APPLY_DATE >='" + startOfYear.Format("2006-01-02") + "' AND APPLY_DATE <='" + endOfYear.Format("2006-01-02") + "'"
	result := ExcuteQuery(query)
	return result
}

func OverTimeDayCheck(body map[string]interface{}, payload map[string]interface{}) string {
	startOfYear := time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	query := `SELECT COUNT(EMPL_NO) AS TANGCA_DAY FROM ZTBATTENDANCETB WHERE EMPL_NO='` + payload["EMPL_NO"].(string) + `' AND ON_OFF=1 AND APPLY_DATE >='` + startOfYear.Format("2006-01-02") + `' AND OVERTIME=1`
	result := ExcuteQuery(query)
	return result
}

func CheckinConfirm(body map[string]interface{}, payload map[string]interface{}) string {
	startOfYear := time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	query := `SELECT COUNT(XACNHAN) AS COUTNXN FROM ZTBATTENDANCETB WHERE EMPL_NO='` + payload["EMPL_NO"].(string) + `' AND XACNHAN is not null AND APPLY_DATE >='` + startOfYear.Format("2006-01-02") + `'`
	result := ExcuteQuery(query)
	return result
}

func CountThuongPhat(body map[string]interface{}, payload map[string]interface{}) string {
	query := `SELECT TP_EMPL_NO, SUM(CASE WHEN PL_HINHTHUC='KT' THEN isnull(DIEM,0) ELSE 0 END) AS THUONG, SUM(CASE WHEN PL_HINHTHUC='KL' THEN isnull(DIEM,0)  ELSE 0 END) AS PHAT FROM ZTBTHUONGPHATTB WHERE TP_EMPL_NO='` + payload["EMPL_NO"].(string) + `' GROUP BY TP_EMPL_NO`
	result := ExcuteQuery(query)
	return result
}

func CheckWebVersion(body map[string]interface{}, payload map[string]interface{}) string {
	query := `SELECT VERWEB, VERMOBILE FROM ZBTVERTABLE`
	result := ExcuteQuery(query)
	return result
}

func CheckMyChamCong(body map[string]interface{}, payload map[string]interface{}) string {
	query1 := `SELECT MIN(C001.CHECK_DATETIME) AS MIN_TIME, MAX(C001.CHECK_DATETIME) AS MAX_TIME  FROM C001 LEFT JOIN ZTBEMPLINFO ON (C001.NV_CCID = ZTBEMPLINFO.NV_CCID) WHERE C001.CHECK_DATE = CAST(GETDATE() as date) AND ZTBEMPLINFO.EMPL_NO='` + payload["EMPL_NO"].(string) + `'`
	result1 := ExcuteQuery(query1)
	query2 := `SELECT ZTBEMPLINFO.EMPL_IMAGE,ZTBEMPLINFO.CTR_CD,ZTBEMPLINFO.EMPL_NO,ZTBEMPLINFO.CMS_ID,ZTBEMPLINFO.FIRST_NAME,ZTBEMPLINFO.MIDLAST_NAME,ZTBEMPLINFO.DOB,ZTBEMPLINFO.HOMETOWN,ZTBEMPLINFO.SEX_CODE,ZTBEMPLINFO.ADD_PROVINCE,ZTBEMPLINFO.ADD_DISTRICT,ZTBEMPLINFO.ADD_COMMUNE,ZTBEMPLINFO.ADD_VILLAGE,ZTBEMPLINFO.PHONE_NUMBER,ZTBEMPLINFO.WORK_START_DATE,ZTBEMPLINFO.PASSWORD,ZTBEMPLINFO.EMAIL,ZTBEMPLINFO.WORK_POSITION_CODE,ZTBEMPLINFO.WORK_SHIFT_CODE,ZTBEMPLINFO.POSITION_CODE,ZTBEMPLINFO.JOB_CODE,ZTBEMPLINFO.FACTORY_CODE,ZTBEMPLINFO.WORK_STATUS_CODE,ZTBEMPLINFO.REMARK,ZTBEMPLINFO.ONLINE_DATETIME,ZTBSEX.SEX_NAME,ZTBSEX.SEX_NAME_KR,ZTBWORKSTATUS.WORK_STATUS_NAME,ZTBWORKSTATUS.WORK_STATUS_NAME_KR,ZTBFACTORY.FACTORY_NAME,ZTBFACTORY.FACTORY_NAME_KR,ZTBJOB.JOB_NAME,ZTBJOB.JOB_NAME_KR,ZTBPOSITION.POSITION_NAME,ZTBPOSITION.POSITION_NAME_KR,ZTBWORKSHIFT.WORK_SHIF_NAME,ZTBWORKSHIFT.WORK_SHIF_NAME_KR,ZTBWORKPOSITION.SUBDEPTCODE,ZTBWORKPOSITION.WORK_POSITION_NAME,ZTBWORKPOSITION.WORK_POSITION_NAME_KR,ZTBWORKPOSITION.ATT_GROUP_CODE,ZTBSUBDEPARTMENT.MAINDEPTCODE,ZTBSUBDEPARTMENT.SUBDEPTNAME,ZTBSUBDEPARTMENT.SUBDEPTNAME_KR,ZTBMAINDEPARMENT.MAINDEPTNAME,ZTBMAINDEPARMENT.MAINDEPTNAME_KR FROM ZTBEMPLINFO LEFT JOIN ZTBSEX ON (ZTBSEX.SEX_CODE = ZTBEMPLINFO.SEX_CODE) LEFT JOIN ZTBWORKSTATUS ON(ZTBWORKSTATUS.WORK_STATUS_CODE = ZTBEMPLINFO.WORK_STATUS_CODE) LEFT JOIN ZTBFACTORY ON (ZTBFACTORY.FACTORY_CODE = ZTBEMPLINFO.FACTORY_CODE) LEFT JOIN ZTBJOB ON (ZTBJOB.JOB_CODE = ZTBEMPLINFO.JOB_CODE) LEFT JOIN ZTBPOSITION ON (ZTBPOSITION.POSITION_CODE = ZTBEMPLINFO.POSITION_CODE) LEFT JOIN ZTBWORKSHIFT ON (ZTBWORKSHIFT.WORK_SHIFT_CODE = ZTBEMPLINFO.WORK_SHIFT_CODE) LEFT JOIN ZTBWORKPOSITION ON (ZTBWORKPOSITION.WORK_POSITION_CODE = ZTBEMPLINFO.WORK_POSITION_CODE) LEFT JOIN ZTBSUBDEPARTMENT ON (ZTBSUBDEPARTMENT.SUBDEPTCODE = ZTBWORKPOSITION.SUBDEPTCODE) LEFT JOIN ZTBMAINDEPARMENT ON (ZTBMAINDEPARMENT.MAINDEPTCODE = ZTBSUBDEPARTMENT.MAINDEPTCODE) WHERE ZTBEMPLINFO.EMPL_NO = '` + payload["EMPL_NO"].(string) + `' AND PASSWORD = '` + payload["PASSWORD"].(string) + `'`
	result2 := ExcuteQuery(query2)

	resultmap1 := make(map[string]interface{})
	err := json.Unmarshal([]byte(result1), &resultmap1)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err.Error())
	}

	resultmap2 := make(map[string]interface{})
	err = json.Unmarshal([]byte(result2), &resultmap2)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err.Error())
	}

	var resultMap map[string]interface{}
	err = json.Unmarshal([]byte(result2), &resultMap)
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
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Fatal("Error signing token:", err.Error())
	}

	//add REFRESH_TOKEN to resultmap1
	resultmap1["REFRESH_TOKEN"] = tokenString

	resultJson, err := json.Marshal(resultmap1)
	if err != nil {
		log.Fatal("Error marshalling JSON:", err.Error())
	}
	return string(resultJson)
}

func CheckLicense(body map[string]interface{}, payload map[string]interface{}) string {
	//fetch data from CURRENT_API_URL
	DATA := body["DATA"].(map[string]interface{})
	resp, err := http.Get(CURRENT_API_URL)
	if err != nil {
		log.Fatal("Error fetching data:", err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err.Error())
	}
	//convert bodyBytes to array of string array
	var bodyArray [][]interface{}
	//bodyBytes is array of string array but in byte, unmarshal it into bodyArray
	err1 := json.Unmarshal(bodyBytes, &bodyArray)
	if err1 != nil {
		log.Fatal("Error unmarshalling JSON:", err1.Error())
	}
	//filter bodyarray elements which have its first element is equal to DATA["COMPANY"]
	var filteredArray [][]interface{}
	for _, v := range bodyArray {
		if v[0] == DATA["COMPANY"] {
			filteredArray = append(filteredArray, v)
		}
	}
	if len(filteredArray) == 0 {
		newJson := map[string]interface{}{
			"tk_status": "NG",
			"message":   "Kiểm tra license thất bại !",
		}
		resultJson, err := json.Marshal(newJson)
		if err != nil {
			log.Fatal("Error marshalling JSON:", err.Error())
		}
		return string(resultJson)
	}
	expdate := filteredArray[0][1]
	parsedExpdate, err := time.Parse("2006-01-02T15:04:05.000Z", expdate.(string))
	if err != nil {
		log.Fatal("Error parsing expdate:", err.Error())
	}
	currentDate := time.Now()
	if parsedExpdate.Before(currentDate) {
		newJson := map[string]interface{}{
			"tk_status": "NG",
			"message":   "License đã hết hạn !",
		}
		resultJson, err := json.Marshal(newJson)
		if err != nil {
			log.Fatal("Error marshalling JSON:", err.Error())
		}
		return string(resultJson)
	}
	newJson := map[string]interface{}{
		"tk_status": "OK",
		"message":   "License còn hạn !",
	}
	resultJson, err := json.Marshal(newJson)
	if err != nil {
		log.Fatal("Error marshalling JSON:", err.Error())
	}
	return string(resultJson)

}
