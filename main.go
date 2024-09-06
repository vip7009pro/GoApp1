package main

import (
	"fmt"
	"go-app/controllers"
	"io"
	"log"
	"net/http"
	"os"

	"encoding/json"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
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
		AllowOrigins:     []string{"http://localhost:3001"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization, echo.HeaderOrigin},
		AllowCredentials: true,
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	//install echo jwt

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("JWT_SECRET")),
		TokenLookup: "cookie:token",
	}))

	//create custom middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestBody := c.Request().Body
			body, err := io.ReadAll(requestBody)
			if err != nil {
				log.Printf("Error reading request body: %v", err)
				return err
			}
			parsedBody := make(map[string]interface{})
			err = json.Unmarshal(body, &parsedBody)
			if err != nil {
				log.Printf("Error unmarshalling JSON: %v", err)
				return err
			}

			fmt.Println("vao day")
			command := parsedBody["command"].(string)
			fmt.Println(command)
			if command == "login" {
				return next(c)
			}
			cookie, err := c.Cookie("token")
			if err != nil {
				return err
			}
			parsedToken, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			if err != nil {
				return err
			}
			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusInternalServerError, "Invalid token claims")
			}
			parsedTokenMap := make(map[string]interface{})
			for key, value := range claims {
				parsedTokenMap[key] = value
			}
			payload := parsedTokenMap["payload"]
			if payloadStr, ok := payload.(string); ok {
				// Try to unmarshal as a map first
				var payloadMap map[string]interface{}
				err := json.Unmarshal([]byte(payloadStr), &payloadMap)
				if err == nil {
					c.Set("payload", payloadMap)
					return next(c)
				}

				// If that fails, try to unmarshal as an array
				var payloadArray []interface{}
				err = json.Unmarshal([]byte(payloadStr), &payloadArray)
				if err == nil {
					c.Set("payload", payloadArray[0])
					return next(c)
				}
				// If both fail, return an error
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse payload: " + err.Error()})
			}
			// If it's neither a string, map, nor array, return an error
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid payload type"})

		}
	})

	e.GET("/", handleURLEncodedForm)
	e.GET("/check", handleURLEncodedForm)
	e.POST("/api", YourAPIFunction)
	e.POST("/uploadfile", UploadFileFunction)

	API_PORT := os.Getenv("API_PORT")
	log.Printf("Server starting at port %s", API_PORT)
	e.Logger.Fatal(e.Start(":" + API_PORT))

}
func handleURLEncodedForm(c echo.Context) error {
	payload := c.Get("payload")
	payloadMap, ok := payload.(map[string]interface{})
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid payload type"})
	}
	emplNo, ok := payloadMap["EMPL_NO"].(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid EMPL_NO type"})
	}
	return c.JSON(http.StatusOK, emplNo)
}

func YourAPIFunction(c echo.Context) error {
	result := controllers.ProcessAPI(c)
	return c.String(http.StatusOK, result)
}

func UploadFileFunction(c echo.Context) error {
	// Handle file upload
	return c.String(http.StatusOK, "Upload File Endpoint")
}
