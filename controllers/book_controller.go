package controllers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sagara/project/upload-image/model"
	"sagara/project/upload-image/repository"
	"sagara/project/upload-image/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BookController struct {
	BookRepository repository.BooksRepoInterface
}

func NewBookController(bookRepository repository.BooksRepoInterface) *BookController {
	return &BookController{
		BookRepository: bookRepository,
	}
}

type BookControllerInterface interface {
	HealthCheck(c *fiber.Ctx) error
	GetBooks(c *fiber.Ctx) error
	GetBookByID(c *fiber.Ctx) error
	CreateBook(c *fiber.Ctx) error
	UploadBookImage(c *fiber.Ctx) error
	UpdateBook(c *fiber.Ctx) error
	DeleteBook(c *fiber.Ctx) error
	GetNewAccessToken(c *fiber.Ctx) error
}

func (h *BookController) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"server":  true,
		"message": "Server UP Capt 🚀",
	})
}

func (h *BookController) GetBooks(c *fiber.Ctx) error {
	books, err := h.BookRepository.GetBooks()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "books were not found",
			"count": 0,
			"books": nil,
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "success get all book",
		"count": len(books),
		"books": books,
	})
}

func (h *BookController) GetBookByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	book, err := h.BookRepository.GetBookByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "book with the given ID is not found",
			"book":  nil,
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "success get book by id",
		"book":  book,
	})
}

func (h *BookController) CreateBook(c *fiber.Ctx) error {
	book := &model.Book{}

	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()

	book.ID = uuid.New()
	book.CreatedAt = time.Now()
	book.BookStatus = 1 // 0 == draft, 1 == active

	if err := validate.Struct(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	if err := h.BookRepository.CreateBook(book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"msg":   "success create data",
		"book":  book,
	})
}

func (h *BookController) UpdateBook(c *fiber.Ctx) error {
	book := &model.Book{}
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	foundedBook, err := h.BookRepository.GetBookByID(book.ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "book with this ID not found",
		})
	}

	book = BuildUpdateRequestData(book, foundedBook)

	if err := h.BookRepository.UpdateBook(foundedBook.ID, book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"msg":   "success update data",
	})
}

func BuildUpdateRequestData(book *model.Book, foundedBook model.Book) *model.Book {
	if book.Author == "" {
		book.Author = foundedBook.Author
	}

	if book.Title == "" {
		book.Title = foundedBook.Title
	}

	if book.BookAttrs.Description == "" {
		book.BookAttrs.Description = foundedBook.BookAttrs.Description
	}

	if book.BookAttrs.Picture == "" {
		book.BookAttrs.Picture = foundedBook.BookAttrs.Picture
	}

	if book.BookAttrs.Rating == 0 {
		book.BookAttrs.Rating = foundedBook.BookAttrs.Rating
	}

	if book.BookStatus == 0 {
		book.BookStatus = foundedBook.BookStatus
	}
	book.UpdatedAt = time.Now()

	return book
}

func (h *BookController) DeleteBook(c *fiber.Ctx) error {
	book := &model.Book{}
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()

	if err := validate.StructPartial(book, "id"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	foundedBook, err := h.BookRepository.GetBookByID(book.ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "book with this ID not found",
		})
	}

	if err := h.BookRepository.DeleteBook(foundedBook.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   fmt.Sprintf("success delete data with id : %s", book.ID),
	})
}

func (h *BookController) UploadBookImage(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	handler, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Server error",
			"error":   err.Error(),
		})
	}
	uploadedFile, _ := handler.Open()

	fileLocation := filepath.Join(os.TempDir(), handler.Filename)
	tempFile, err := os.Create(fileLocation)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Server error",
			"error":   err.Error(),
		})
	}
	defer tempFile.Close()

	io.Copy(tempFile, uploadedFile)
	urlCloudinary, err := h.BookRepository.UploadToCloudinary(fileLocation)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Server error",
			"error":   err.Error(),
		})
	}

	err = h.BookRepository.UpdateUrlBookImage(urlCloudinary, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Server error",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"msg":   "success upload picture",
	})
}
