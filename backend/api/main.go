package main

import (
	"log"

	"Server/database"
	_ "Server/docs"
	"Server/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @title Fiber Golang mongo Grpc websocket etc...
// @version 1.0
// @description This is Swagger docs for resapi Golang Fiber
// @host localhost:5000
// @BasePath /
// @schemes http
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by space and the token

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error Load .Env file")
	}

	database.Connect()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	//setup routes
	routes.SetupAuthRoutes(app)
	routes.SetupUserRoutes(app)

	//Serve swagger doctionation
	app.Get("/swagger/*", swagger.HandlerDefault)
	log.Fatal(app.Listen(":5000"))
}
