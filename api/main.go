package main

import (
	"fmt"
	"log"
	"shorten-url-go-redis/routes"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL) //resolve the URL
	app.Post("/api/v1/shorten", routes.ShortenURL) //shorten the URL

}

func main() {

	err := godotenv.Load() //load the .env file
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	app := fiber.New() //create a new fiber app
	app.Use(logger.New()) //middleware to log the requests
	setupRoutes(app) //setup the routes

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
