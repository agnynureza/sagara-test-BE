package routes

import (
	"sagara/project/upload-image/controllers"
	"sagara/project/upload-image/middleware"

	"github.com/gofiber/fiber/v2"
)

type BookRoute struct {
	BookController controllers.BookControllerInterface
}

func NewBookRoute(BookController controllers.BookControllerInterface) *BookRoute {
	return &BookRoute{
		BookController: BookController,
	}
}

func (b *BookRoute) InitializeBookRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	//check health server
	route.Get("/health", b.BookController.HealthCheck)

	//book route
	route.Post("/book", middleware.JWTProtected(), b.BookController.CreateBook)
	route.Post("/book/:id/upload-image", middleware.JWTProtected(), b.BookController.UploadBookImage)
	route.Get("/books", middleware.JWTProtected(), b.BookController.GetBooks)
	route.Get("/book/:id", middleware.JWTProtected(), b.BookController.GetBookByID)
	route.Put("/book", middleware.JWTProtected(), b.BookController.UpdateBook)
	route.Delete("/book", middleware.JWTProtected(), b.BookController.DeleteBook)

	// login
	route.Get("/login", b.BookController.GetNewAccessToken)
}
