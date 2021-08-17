package repository

import (
	"context"
	"sagara/project/upload-image/model"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BookQueries struct {
	db         *sqlx.DB
	cloudinary *cloudinary.Cloudinary
}

func NewBookQueries(db *sqlx.DB, cloudinary *cloudinary.Cloudinary) *BookQueries {
	return &BookQueries{
		db:         db,
		cloudinary: cloudinary,
	}
}

type BooksRepoInterface interface {
	GetBooks() ([]model.Book, error)
	GetBookByID(id uuid.UUID) (model.Book, error)
	CreateBook(b *model.Book) error
	UpdateBook(id uuid.UUID, b *model.Book) error
	DeleteBook(id uuid.UUID) error
	UploadToCloudinary(file string) (string, error)
	UpdateUrlBookImage(url string, id uuid.UUID) error
}

func (repo *BookQueries) GetBooks() ([]model.Book, error) {
	books := []model.Book{}

	query := `SELECT * FROM books`

	err := repo.db.Select(&books, query)
	if err != nil {
		return books, err
	}

	return books, nil
}

func (repo *BookQueries) GetBookByID(id uuid.UUID) (model.Book, error) {
	book := model.Book{}

	query := `SELECT * FROM books WHERE id = $1`

	err := repo.db.Get(&book, query, id)
	if err != nil {
		return book, err
	}

	return book, nil
}

func (repo *BookQueries) CreateBook(b *model.Book) error {
	query := `INSERT INTO books (id, created_at, updated_at, title, author, book_status, book_attrs)
				VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := repo.db.Exec(query, b.ID, b.CreatedAt, b.UpdatedAt, b.Title, b.Author, b.BookStatus, b.BookAttrs)
	if err != nil {
		return err
	}

	return nil
}

func (repo *BookQueries) UpdateBook(id uuid.UUID, b *model.Book) error {
	query := `UPDATE books SET updated_at = $2, title = $3, author = $4, book_status = $5, book_attrs = $6 WHERE id = $1`

	_, err := repo.db.Exec(query, id, b.UpdatedAt, b.Title, b.Author, b.BookStatus, b.BookAttrs)
	if err != nil {
		return err
	}

	return nil
}
func (repo *BookQueries) UpdateUrlBookImage(url string, id uuid.UUID) error {
	query := `UPDATE books SET 
				book_attrs = jsonb_set("book_attrs", '{"picture"}', to_jsonb($1::text), true),
				updated_at = $2
			WHERE id = $3`

	_, err := repo.db.Exec(query, url, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *BookQueries) DeleteBook(id uuid.UUID) error {
	query := `DELETE FROM books WHERE id = $1`

	_, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *BookQueries) UploadToCloudinary(file string) (string, error) {
	var ctx = context.Background()
	uploadResult, err := repo.cloudinary.Upload.Upload(
		ctx,
		file,
		uploader.UploadParams{PublicID: "my_image"})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
