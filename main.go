package main

import (
	"fmt"
	"go-app/controllers"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	e := echo.New()
	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3001"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	e.GET("/", handleURLEncodedForm)
	e.GET("/check", handleURLEncodedForm)
	e.POST("/api", YourAPIFunction)
	e.POST("/uploadfile", UploadFileFunction)

	API_PORT := os.Getenv("API_PORT")
	log.Printf("Server starting at port %s", API_PORT)
	e.Logger.Fatal(e.Start(":" + API_PORT))

}
func handleURLEncodedForm(c echo.Context) error {
	result := controllers.ExcuteQuery("SELECT * FROM ZTB_REL_TESTPOINT")
	fmt.Println(result)
	return c.String(http.StatusOK, result)
}

func YourAPIFunction(c echo.Context) error {
	// Handle API request
	fmt.Println("YourAPIFunction")
	requestBody := c.Request().Body
	body, err := io.ReadAll(requestBody)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		return c.String(http.StatusInternalServerError, "Error reading request body")
	}

	fmt.Println(string(body))
	result := controllers.ProcessAPI()
	return c.String(http.StatusOK, result)
}

func UploadFileFunction(c echo.Context) error {
	// Handle file upload
	return c.String(http.StatusOK, "Upload File Endpoint")
}
