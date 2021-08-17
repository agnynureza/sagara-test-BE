package main

import (
	"sagara/project/upload-image/configs"
	"sagara/project/upload-image/configs/client"
	"sagara/project/upload-image/configs/database"
	"sagara/project/upload-image/controllers"
	"sagara/project/upload-image/middleware"
	"sagara/project/upload-image/repository"
	"sagara/project/upload-image/routes"
	"sagara/project/upload-image/utils"

	"github.com/gofiber/fiber/v2"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config := configs.FiberConfig()

	app := fiber.New(config)

	//middleware
	middleware.FiberMiddleware(app)

	//set up connection
	connDB, err := database.PostgreSQLConnection()
	if err != nil {
		panic(err)
	}

	connCloudinary, err := client.InitCloudinary()
	if err != nil {
		panic(err)
	}

	bookRepository := repository.NewBookQueries(connDB, connCloudinary)
	bookController := controllers.NewBookController(bookRepository)
	BookRoutes := routes.NewBookRoute(bookController)

	// Routes
	BookRoutes.InitializeBookRoutes(app)

	//not found routes
	routes.NotFoundRoute(app)

	// Start server (with graceful shutdown)
	utils.StartServerWithGracefulShutdown(app)
}
