package main

import (
	"fmt"
	"go-app/controllers"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"encoding/json"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/golang-jwt/jwt/v5"
	socketio "github.com/googollee/go-socket.io"
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
	f := echo.New()
	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3001", "http://cms.ddns.net:3002"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization, echo.HeaderOrigin},
		AllowCredentials: true,
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	f.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3001", "http://cms.ddns.net:3002"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization, echo.HeaderOrigin},
		AllowCredentials: true,
	}))
	f.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	//install echo jwt

	/* 	f.Use(middleware.Logger())
	   	f.Use(middleware.Recover()) */

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
			command := parsedBody["command"].(string)
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), command)
			if command == "login" {
				c.Set("body", parsedBody)
				return next(c)
			}
			cookie, err := c.Cookie("token")
			if err != nil {
				return err
			}
			//fmt.Println(cookie.Value)
			parsedToken, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			if err != nil {
				return err
			}
			//fmt.Printf("parsedToken: %+v\n", parsedToken)
			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusInternalServerError, "Invalid token claims")
			}
			//fmt.Printf("claims: %+v\n", claims)
			parsedTokenMap := make(map[string]interface{})
			for key, value := range claims {
				parsedTokenMap[key] = value
			}
			//fmt.Printf("parsedTokenMap: %+v\n", parsedTokenMap)
			payload := parsedTokenMap["payload"]
			//fmt.Println(payload)
			c.Set("body", parsedBody)
			c.Set("payload", payload)
			return next(c)
		}

	})

	/*
		e.Use(echojwt.WithConfig(echojwt.Config{
			SigningKey: []byte(os.Getenv("JWT_SECRET")),
			TokenLookup: "cookie:token",
		})) */

	// Initialize Socket.IO server
	server := socketio.NewServer(nil)

	// Handle new connections
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("New connection:", s.ID())
		return nil
	})

	server.OnEvent("/", "message", func(s socketio.Conn, msg string) {
		log.Printf("Received message: %s", msg)
		s.Emit("reply", "Message received: "+msg)
	})

	// Handle disconnections
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("Disconnected:", s.ID(), "Reason:", reason)
	})

	// Serve Socket.IO
	go server.Serve()
	defer server.Close()

	// Attach Socket.IO server to HTTP server
	/* f.GET("/socket.io/*", func(c echo.Context) error {
		fmt.Println("Socket.IO connection")
		return handleSocketIO(c, server)
	}) */
	f.Any("/socket.io/*", echo.WrapHandler(server))

	e.GET("/", handleURLEncodedForm)
	e.GET("/check", handleURLEncodedForm)
	e.POST("/api", YourAPIFunction)
	e.POST("/uploadfile", UploadFileFunction)

	API_PORT := os.Getenv("API_PORT")
	SOCKET_PORT := os.Getenv("SOCKET_PORT")
	log.Printf("Http Server starting at port %s", API_PORT)
	log.Printf("Socket.IO Server starting at port %s", SOCKET_PORT)

	go func() {
		f.Logger.Fatal(f.Start(":3006"))
	}()
	go func() {
		f.Logger.Fatal(e.Start(":3002"))
	}()
	select {}
}

func handleSocketIO(c echo.Context, server *socketio.Server) error {
	server.ServeHTTP(c.Response(), c.Request())
	return nil
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
	body := c.Get("body")
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid body type"})
	}
	//fmt.Println("vao YourAPIFunction")
	payload := c.Get("payload")
	payloadMap := make(map[string]interface{})
	if payload != nil {
		payloadMap, ok = payload.(map[string]interface{})
		if !ok {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid payload type"})
		}
	}
	result := controllers.ProcessAPI(bodyMap, payloadMap)
	parsedResult := make(map[string]interface{})
	err := json.Unmarshal([]byte(result), &parsedResult)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse result: " + err.Error()})
	}
	return c.JSON(http.StatusOK, parsedResult)
}
func UploadFileFunction(c echo.Context) error {
	// Handle file upload
	return c.String(http.StatusOK, "Upload File Endpoint")
}
