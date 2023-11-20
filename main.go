package main

import (
	"fmt"
	"healthcare/configs"
	"healthcare/middlewares"
	"healthcare/routes"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println("Can't access .env files!")
	// }

	configs.Init()
	e := echo.New()

	// load middlewares
	middlewares.RemoveTrailingSlash(e)
	middlewares.Logger(e)
	middlewares.RateLimiter(e)
	middlewares.Recover(e)
	middlewares.CORS(e)

	// load router
	routes.SetupRoutes(e)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	setPort := fmt.Sprintf(":%s", port)
	e.Logger.Fatal(e.Start(setPort))
}
