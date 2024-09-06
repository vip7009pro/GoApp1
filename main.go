package main

import (
	"go-app/controllers"
	"log"
	"net/http"
	"os"

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
			c.Set("payload", payload)
			return next(c)
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
	//result := controllers.ExcuteQuery("SELECT * FROM ZTB_REL_TESTPOINT")
	//read token from cookie
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

	return c.JSON(http.StatusOK, parsedTokenMap)

}

func YourAPIFunction(c echo.Context) error {
	result := controllers.ProcessAPI(c)
	return c.String(http.StatusOK, result)
}

func UploadFileFunction(c echo.Context) error {
	// Handle file upload
	return c.String(http.StatusOK, "Upload File Endpoint")
}
